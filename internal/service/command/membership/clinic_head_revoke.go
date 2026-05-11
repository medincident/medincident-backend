package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// RevokeClinicHeadPayload carries the identifiers needed to remove an
// employee's clinic head role.
type RevokeClinicHeadPayload struct {
	ClinicID   string `validate:"required,uuid"`
	EmployeeID string `validate:"required,uuid"`
}

// RevokeClinicHeadCommand = caller + payload.
type RevokeClinicHeadCommand struct {
	Caller  authz.Caller
	Payload RevokeClinicHeadPayload
}

// RevokeClinicHead removes the role row and publishes the Revoked
// event. If a deputy was assigned, a DeputyRemoved event is published
// FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
//
// See: docs/services/Membership.md
func (s *EmployeeService) RevokeClinicHead(ctx context.Context, cmd RevokeClinicHeadCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	clinicID := uuid.MustParse(cmd.Payload.ClinicID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var row model.ClinicHead
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("clinic_id = ? AND employee_id = ?", clinicID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeClinicHead).
					Code(ErrCodeClinicHeadNotFound).
					Public("Clinic head not found.").
					With("clinic_id", clinicID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadLoadFailed).Wrap(err)
		}

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishClinicHeadDeputyRemoved(tx, clinicID, employeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.ClinicHead{}, "clinic_id = ? AND employee_id = ?",
			clinicID, employeeID).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadDeleteFailed).Wrap(err)
		}

		return publishClinicHeadRevoked(tx, clinicID, employeeID, now)
	})
}

// publishClinicHeadRevoked is a shared helper for Revoke,
// cascade-on-transfer, and cascade-on-terminate. It does NOT delete
// the domain row; the caller owns that.
func publishClinicHeadRevoked(tx *gorm.DB, clinicID, employeeID uuid.UUID, now time.Time) error {
	env, err := buildClinicHeadRevokedEnvelope(clinicID, employeeID, now)
	if err != nil {
		return err
	}
	return outbox.Append(tx, "medincident.event.clinic.v1.clinic_head_revoked", env)
}

// publishClinicHeadDeputyRemoved is a shared helper; it clears the
// deputy slot on the projection row. The caller owns the domain-row
// update.
func publishClinicHeadDeputyRemoved(tx *gorm.DB, clinicID, employeeID uuid.UUID, now time.Time) error {
	env, err := buildClinicHeadDeputyRemovedEnvelope(clinicID, employeeID, now)
	if err != nil {
		return err
	}
	return outbox.Append(tx, "medincident.event.clinic.v1.clinic_head_deputy_removed", env)
}

func buildClinicHeadRevokedEnvelope(clinicID, employeeID uuid.UUID, now time.Time) (*eventv1.Envelope, error) {
	msg := &clinicv1.ClinicHeadRevoked{EmployeeId: employeeID.String()}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeClinicHead).Code(ErrCodeClinicHeadDeleteFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "clinic",
		AggregateId:   clinicID.String(),
		Payload:       payload,
	}, nil
}

func buildClinicHeadDeputyRemovedEnvelope(clinicID, employeeID uuid.UUID, now time.Time) (*eventv1.Envelope, error) {
	msg := &clinicv1.ClinicHeadDeputyRemoved{EmployeeId: employeeID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "clinic",
		AggregateId:   clinicID.String(),
		Payload:       payload,
	}, nil
}
