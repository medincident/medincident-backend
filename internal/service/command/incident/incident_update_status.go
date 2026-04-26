package incident

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

type UpdateIncidentStatusPayload struct {
	IncidentID string `validate:"required,uuid"`
	NewStatus  string `validate:"required,oneof=in_progress done rejected"`
}

type UpdateIncidentStatusCommand struct {
	Caller  authz.Caller
	Payload UpdateIncidentStatusPayload
}

// UpdateStatus performs forward-only transitions. The set of valid
// transitions is fixed: pending->in_progress, in_progress->done,
// in_progress->rejected. Cancellation is a separate RPC.
//
// See: docs/services/incident/Incidents.md
func (s *IncidentService) UpdateStatus(
	ctx context.Context, cmd UpdateIncidentStatusCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.IncidentID)
	newStatus := model.IncidentStatus(cmd.Payload.NewStatus)
	now := time.Now()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		inc, err := s.loadIncident(tx, id)
		if err != nil {
			return err
		}
		if !validStatusTransition(inc.Status, newStatus) {
			return oops.In(scope).
				Code(ErrCodeIncidentInvalidStatusFlow).
				Public("This status transition is not allowed.").
				With("from", inc.Status).
				With("to", newStatus).
				Errorf("invalid transition")
		}
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			privilegedActorPolicy(inc.OrganizationID, inc.ClinicID, inc.DepartmentID)); err != nil {
			return err
		}
		actorEmpID, displayName, err := s.resolveActor(tx, cmd.Caller.ZitadelUserID)
		if err != nil {
			return err
		}

		old := inc.Status
		inc.Status = newStatus
		inc.UpdatedAt = now
		if err := tx.Save(inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
		}
		return projector.IncidentStatusChanged(tx, inc.ID, old, inc.Status,
			actorEmpID, displayName, now)
	})
}

// validStatusTransition encodes the allowed forward transitions.
func validStatusTransition(from, to model.IncidentStatus) bool {
	switch from {
	case model.IncidentStatusPending:
		return to == model.IncidentStatusInProgress
	case model.IncidentStatusInProgress:
		return to == model.IncidentStatusDone || to == model.IncidentStatusRejected
	default:
		return false
	}
}

// resolveActor returns the caller's employee_id (in any org) and the
// display_name from projections.users. Used for history rows. The
// returned NullUUID is invalid when the caller has no employee row
// (e.g. SystemAdmin acting outside any org) so history tables get a
// proper SQL NULL rather than the zero UUID.
func (s *IncidentService) resolveActor(tx *gorm.DB, callerID string) (uuid.NullUUID, string, error) {
	// Caller may be a SystemAdmin without a domain.employees row; that
	// is not an error here — we treat the actor's employee link as
	// optional. Any error other than "not found" is propagated below.
	var emp model.Employee
	empID := uuid.NullUUID{}
	if err := tx.Where("zitadel_user_id = ?", callerID).Limit(1).First(&emp).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.NullUUID{}, "", oops.In(scope).
				Code(ErrCodeIncidentEmployeeNotFound).
				With("zitadel_user_id", callerID).Wrap(err)
		}
	} else {
		empID = uuid.NullUUID{UUID: emp.ID, Valid: true}
	}

	var displayName string
	if err := tx.Raw(
		`SELECT display_name FROM projections.users WHERE id = ?`,
		callerID,
	).Scan(&displayName).Error; err != nil {
		return uuid.NullUUID{}, "", oops.In(scope).
			Code(ErrCodeIncidentRegistrarUserNotFound).
			With("zitadel_user_id", callerID).Wrap(err)
	}
	return empID, displayName, nil
}
