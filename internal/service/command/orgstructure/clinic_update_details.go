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

// UpdateClinicDetailsPayload carries the new name and (optional)
// description for an existing clinic.
type UpdateClinicDetailsPayload struct {
	ID          string  `validate:"required,uuid"`
	Name        string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// UpdateClinicDetailsCommand = caller + payload.
type UpdateClinicDetailsCommand struct {
	Caller  authz.Caller
	Payload UpdateClinicDetailsPayload
}

func buildClinicDetailsChangedEnvelope(c *model.Clinic) (*eventv1.Envelope, error) {
	var desc string
	if c.Description.Valid {
		desc = c.Description.String
	}
	msg := &clinicv1.ClinicDetailsChanged{
		Name:        c.Name,
		Description: desc,
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
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

// UpdateDetails changes a clinic's name and description.
//
// See: docs/services/OrgStructure.md
func (s *ClinicService) UpdateDetails(
	ctx context.Context,
	cmd UpdateClinicDetailsCommand,
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

		newName := strings.TrimSpace(cmd.Payload.Name)
		var newDesc null.String
		if cmd.Payload.Description != nil {
			newDesc = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if clinic.Name == newName && clinic.Description == newDesc {
			return nil
		}
		clinic.Name = newName
		clinic.Description = newDesc
		if err := tx.Save(&clinic).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", id).
				Wrap(err)
		}

		env, err := buildClinicDetailsChangedEnvelope(&clinic)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.clinic.v1.details_changed", env)
	})
}
