package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// Error codes emitted by Clinic-aggregate commands that are not
// primitive validation. Validation codes are generic and live in
// internal/validation.
const (
	ErrCodeClinicIDGenerationFailed   = "clinic_id_generation_failed"
	ErrCodeClinicSaveFailed           = "clinic_save_failed"
	ErrCodeClinicLoadFailed           = "clinic_load_failed"
	ErrCodeClinicNotFound             = "clinic_not_found"
	ErrCodeClinicOrganizationNotFound = "clinic_organization_not_found"
)

// CreateClinicPayload is the validated client-facing payload of
// CreateClinic.
type CreateClinicPayload struct {
	OrganizationID  string  `validate:"required,uuid"`
	Name            string  `validate:"required,no_extra_ws,min=4,max=256"`
	Description     *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
	PhysicalAddress AddressInput
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

// Create persists a new Clinic under the given organization.
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

		if err := projector.ClinicCreated(tx, &clinic); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
