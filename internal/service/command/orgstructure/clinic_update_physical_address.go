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
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// UpdateClinicPhysicalAddressPayload carries the new physical address
// for an existing clinic.
type UpdateClinicPhysicalAddressPayload struct {
	ID      string       `validate:"required,uuid"`
	Address AddressInput `validate:"required"`
}

// UpdateClinicPhysicalAddressCommand = caller + payload.
type UpdateClinicPhysicalAddressCommand struct {
	Caller  authz.Caller
	Payload UpdateClinicPhysicalAddressPayload
}

func buildClinicPhysicalAddressChangedEnvelope(c *model.Clinic) (*eventv1.Envelope, error) {
	msg := &clinicv1.ClinicPhysicalAddressChanged{
		PhysicalAddress: buildClinicAddressProto(c.PhysicalAddress),
		UpdatedAt:       timestamppb.New(c.UpdatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(c.UpdatedAt),
		AggregateType: "clinic",
		AggregateId:   c.ID.String(),
		Payload:       payload,
	}, nil
}

// UpdatePhysicalAddress replaces the clinic's physical address.
//
// See: docs/services/OrgStructure.md
func (s *ClinicService) UpdatePhysicalAddress(
	ctx context.Context,
	cmd UpdateClinicPhysicalAddressCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(id)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var clinic model.Clinic
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&clinic, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", id).
					Errorf("clinic not found")
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", id).
				Wrap(err)
		}

		newAddress := model.Address{Text: strings.TrimSpace(cmd.Payload.Address.Text)}
		if cmd.Payload.Address.Point != nil {
			newAddress.Point = null.ValueFrom(model.Point{
				Longitude: cmd.Payload.Address.Point.Longitude,
				Latitude:  cmd.Payload.Address.Point.Latitude,
			})
		}
		if clinic.PhysicalAddress.Equal(newAddress) {
			return nil
		}
		clinic.PhysicalAddress = newAddress
		if err := tx.Save(&clinic).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", id).
				Wrap(err)
		}

		env, err := buildClinicPhysicalAddressChangedEnvelope(&clinic)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.clinic.v1.physical_address_changed", env)
	})
}
