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

type UpdateIncidentPriorityPayload struct {
	IncidentID string                 `validate:"required,uuid"`
	Priority   model.IncidentPriority `validate:"required,oneof=low normal high critical"`
}

type UpdateIncidentPriorityCommand struct {
	Caller  authz.Caller
	Payload UpdateIncidentPriorityPayload
}

// UpdatePriority sets the incident priority. Allowed only for
// privileged roles, only while the incident is not terminal.
// No-op if priority is unchanged (no history row written).
//
// See: docs/services/incident/Incidents.md
func (s *IncidentService) UpdatePriority(
	ctx context.Context, cmd *UpdateIncidentPriorityCommand,
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
		if inc.Status.IsTerminal() {
			return oops.In(scope).Code(ErrCodeIncidentFrozen).
				Public("Incident is frozen and cannot be modified.").
				With("status", inc.Status).Errorf("frozen")
		}
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			privilegedActorPolicy(inc.OrganizationID, inc.ClinicID, inc.DepartmentID)); err != nil {
			return err
		}
		if inc.Priority == cmd.Payload.Priority {
			return nil
		}
		actorEmpID, displayName, err := s.resolveActor(tx, cmd.Caller.ZitadelUserID)
		if err != nil {
			return err
		}
		old := inc.Priority
		inc.Priority = cmd.Payload.Priority
		inc.UpdatedAt = now
		if err := tx.Save(inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
		}
		return projector.IncidentPriorityChanged(tx, inc.ID, old, inc.Priority,
			actorEmpID, displayName, now)
	})
}
