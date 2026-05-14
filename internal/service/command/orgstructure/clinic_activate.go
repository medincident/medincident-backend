package orgstructure

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
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

// ActivateClinicPayload identifies the clinic to activate.
type ActivateClinicPayload struct {
	ID string `validate:"required,uuid"`
}

// ActivateClinicCommand = caller + payload.
type ActivateClinicCommand struct {
	Caller  authz.Caller
	Payload ActivateClinicPayload
}

// Activate marks a deactivated clinic as active. The parent organization
// must be active; child entities keep their individual is_active state.
//
// See: docs/services/OrgStructure.md
func (s *ClinicService) Activate(
	ctx context.Context,
	cmd ActivateClinicCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	clinicID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var clinic model.Clinic
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&clinic, "id = ?", clinicID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", clinicID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}

		if clinic.IsActive {
			return nil
		}

		// Check that the parent organization is active.
		var orgActive bool
		if err := tx.Raw(
			`SELECT is_active FROM domain.organizations WHERE id = ?`,
			clinic.OrganizationID,
		).Row().Scan(&orgActive); err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}
		if !orgActive {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicActivateParentInactive).
				Public("Cannot activate clinic: parent organization is inactive.").
				With("clinic_id", clinicID).
				With("organization_id", clinic.OrganizationID).
				Errorf("inactive parent organization")
		}

		var updatedAt time.Time
		if err := tx.Raw(
			`UPDATE domain.clinics SET is_active = TRUE, updated_at = now() WHERE id = ? RETURNING updated_at`,
			clinicID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}
		return appendClinicActivatedEvent(tx, clinicID, updatedAt)
	})
}

func appendClinicActivatedEvent(tx *gorm.DB, clinicID uuid.UUID, now time.Time) error {
	msg := &clinicv1.ClinicActivated{ClinicId: clinicID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.clinic").Code(ErrCodeClinicSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "clinic",
		AggregateId:   clinicID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.clinic.v1.activated", env)
}
