package incident

import (
	"context"
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
	IncidentID string               `validate:"required,uuid"`
	NewStatus  model.IncidentStatus `validate:"required,oneof=in_progress done rejected"`
}

type UpdateIncidentStatusCommand struct {
	Caller  authz.Caller
	Payload UpdateIncidentStatusPayload
}

// UpdateStatus performs forward-only transitions. The set of valid
// transitions is fixed: pending->in_progress, in_progress->done,
// in_progress->rejected. Cancellation is a separate RPC.
func (s *IncidentService) UpdateStatus(
	ctx context.Context, cmd *UpdateIncidentStatusCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.IncidentID)
	now := time.Now()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		inc, err := s.loadIncident(tx, id)
		if err != nil {
			return err
		}
		if !validStatusTransition(inc.Status, cmd.Payload.NewStatus) {
			return oops.In(scope).
				Code(ErrCodeIncidentInvalidStatusFlow).
				Public("This status transition is not allowed.").
				With("from", inc.Status).
				With("to", cmd.Payload.NewStatus).
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
		inc.Status = cmd.Payload.NewStatus
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
// display_name from projections.users. Used for history rows.
func (s *IncidentService) resolveActor(tx *gorm.DB, callerID string) (uuid.UUID, string, error) {
	// Caller might be SystemAdmin without a domain.employees row;
	// in that case actor_employee_id is the zero UUID and display_name
	// comes from projections.users only.
	var emp model.Employee
	_ = tx.Where("zitadel_user_id = ?", callerID).Limit(1).First(&emp).Error

	var displayName string
	if err := tx.Raw(
		`SELECT display_name FROM projections.users WHERE id = ?`,
		callerID,
	).Scan(&displayName).Error; err != nil {
		return uuid.Nil, "", oops.In(scope).
			Code(ErrCodeIncidentRegistrarUserNotFound).
			With("zitadel_user_id", callerID).Wrap(err)
	}
	return emp.ID, displayName, nil
}
