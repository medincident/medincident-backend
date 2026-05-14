package orgstructure

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	"github.com/medincident/medincident-backend/internal/service/validation"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// DeleteClinicPayload identifies the clinic to delete.
type DeleteClinicPayload struct {
	ID string `validate:"required,uuid"`
}

// DeleteClinicCommand = caller + payload.
type DeleteClinicCommand struct {
	Caller  authz.Caller
	Payload DeleteClinicPayload
}

// Delete hard-deletes a clinic and all its children in the correct FK
// order: roles → employees → departments → clinic.
//
// See: docs/services/OrgStructure.md
func (s *ClinicService) Delete(
	ctx context.Context,
	cmd DeleteClinicCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	clinicID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		// Collect and delete all employees in this clinic.
		var empIDs []uuid.UUID
		if err := tx.Raw(
			`SELECT id FROM domain.employees WHERE department_id IN (SELECT id FROM domain.departments WHERE clinic_id = ?) FOR UPDATE`,
			clinicID,
		).Scan(&empIDs).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}

		for _, empID := range empIDs {
			if err := membership.RevokeAllEmployeeRoles(tx, empID, now); err != nil {
				return err
			}
		}
		if len(empIDs) > 0 {
			if err := tx.Exec(
				`DELETE FROM domain.employees WHERE id = ANY(?)`, empIDs,
			).Error; err != nil {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicDeleteFailed).
					With("clinic_id", clinicID).
					Wrap(err)
			}
		}

		// Delete departments.
		if err := tx.Exec(
			`DELETE FROM domain.departments WHERE clinic_id = ?`, clinicID,
		).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicDeleteFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}

		// Delete the clinic itself.
		if err := tx.Delete(&model.Clinic{}, "id = ?", clinicID).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicDeleteHasDependents).
					Public("Clinic has dependents and cannot be deleted.").
					With("clinic_id", clinicID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicDeleteFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}

		return appendClinicDeletedEvent(tx, clinicID, now)
	})
}

func appendClinicDeletedEvent(tx *gorm.DB, clinicID uuid.UUID, now time.Time) error {
	msg := &clinicv1.ClinicDeleted{ClinicId: clinicID.String(), DeletedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.clinic").Code(ErrCodeClinicDeleteFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "clinic",
		AggregateId:   clinicID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.clinic.v1.deleted", env)
}
