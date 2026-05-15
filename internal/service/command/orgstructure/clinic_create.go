package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// Error codes emitted by Clinic-aggregate commands that are not
// primitive validation. Validation codes are generic and live in
// internal/validation.
const (
	ErrCodeClinicIDGenerationFailed     = "clinic_id_generation_failed"
	ErrCodeClinicSaveFailed             = "clinic_save_failed"
	ErrCodeClinicLoadFailed             = "clinic_load_failed"
	ErrCodeClinicNotFound               = "clinic_not_found"
	ErrCodeClinicOrganizationNotFound   = "clinic_organization_not_found"
	ErrCodeClinicDeleteFailed           = "clinic_delete_failed"
	ErrCodeClinicDeleteHasDependents    = "clinic_delete_has_dependents"
	ErrCodeClinicActivateParentInactive = "clinic_activate_parent_inactive"
)

// CreateClinicPayload is the validated client-facing payload of
// CreateClinic.
type CreateClinicPayload struct {
	OrganizationID  string       `validate:"required,uuid"`
	Name            string       `validate:"required,no_extra_ws,min=2,max=256"`
	Description     *string      `validate:"omitnil,no_extra_ws,min=8,max=2048"`
	PhysicalAddress AddressInput `validate:"required"`
}

// CreateClinicCommand = caller + payload.
type CreateClinicCommand struct {
	Caller  authz.Caller
	Payload CreateClinicPayload
}

// CreateClinicResult is the output of ClinicService.Create.
type CreateClinicResult struct {
	ID uuid.UUID
}

// buildClinicAddressProto converts a model.Address to the clinic proto Address.
func buildClinicAddressProto(a model.Address) *clinicv1.Address {
	addr := &clinicv1.Address{Text: a.Text}
	if a.Point.Valid {
		addr.Point = &clinicv1.Point{
			Longitude: a.Point.V.Longitude,
			Latitude:  a.Point.V.Latitude,
		}
	}
	return addr
}

func buildClinicCreatedEnvelope(c *model.Clinic) (*eventv1.Envelope, error) {
	addr := &clinicv1.Address{Text: c.PhysicalAddress.Text}
	if c.PhysicalAddress.Point.Valid {
		addr.Point = &clinicv1.Point{
			Longitude: c.PhysicalAddress.Point.V.Longitude,
			Latitude:  c.PhysicalAddress.Point.V.Latitude,
		}
	}
	var desc string
	if c.Description.Valid {
		desc = c.Description.String
	}
	msg := &clinicv1.ClinicCreated{
		OrganizationId:  c.OrganizationID.String(),
		Name:            c.Name,
		Description:     desc,
		PhysicalAddress: addr,
		CreatedAt:       timestamppb.New(c.CreatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(c.CreatedAt),
		AggregateType: "clinic",
		AggregateId:   c.ID.String(),
		Payload:       payload,
	}, nil
}

// Create persists a new Clinic under the given organization.
//
// See: docs/services/OrgStructure.md
//
//nolint:gocritic // hugeParam: Command is passed by value across the whole service layer for consistency; CreateClinicCommand is borderline at 80 bytes but not worth breaking the convention for.
func (s *ClinicService) Create(
	ctx context.Context,
	cmd CreateClinicCommand,
) (CreateClinicResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateClinicResult{}, err
	}
	orgID := uuid.MustParse(cmd.Payload.OrganizationID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return CreateClinicResult{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateClinicResult{}, oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicIDGenerationFailed).
			Public("Failed to create clinic.").
			Wrap(err)
	}

	clinic := model.Clinic{
		ID:             id,
		OrganizationID: orgID,
		Name:           strings.TrimSpace(cmd.Payload.Name),
		IsActive:       true,
		PhysicalAddress: model.Address{
			Text: strings.TrimSpace(cmd.Payload.PhysicalAddress.Text),
		},
	}
	if cmd.Payload.Description != nil {
		clinic.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
	}
	if cmd.Payload.PhysicalAddress.Point != nil {
		clinic.PhysicalAddress.Point = null.ValueFrom(model.Point{
			Longitude: cmd.Payload.PhysicalAddress.Point.Longitude,
			Latitude:  cmd.Payload.PhysicalAddress.Point.Latitude,
		})
	}

	var result CreateClinicResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&clinic).Error; err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", orgID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", id).
				Wrap(err)
		}

		env, err := buildClinicCreatedEnvelope(&clinic)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.clinic.v1.created", env); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
