package incident

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

type CancelIncidentPayload struct {
	IncidentID string `validate:"required,uuid"`
}

type CancelIncidentCommand struct {
	Caller  authz.Caller
	Payload CancelIncidentPayload
}

// Cancel marks the incident cancelled. Allowed only if the caller is
// the active registrar AND the incident is still pending.
//
// See: docs/services/incident/Incidents.md
func (s *IncidentService) Cancel(
	ctx context.Context, cmd CancelIncidentCommand,
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
		isReg, err := s.callerIsActiveRegistrar(tx, cmd.Caller.ZitadelUserID,
			inc.RegistrarEmployeeID)
		if err != nil {
			return err
		}
		if !isReg {
			return oops.In(scope).
				Code(authz.ErrCodePermissionDenied).
				Public("Only the registrar may cancel this incident.").
				Errorf("not registrar")
		}
		if inc.Status != model.IncidentStatusPending {
			return oops.In(scope).
				Code(ErrCodeIncidentNotCancellable).
				Public("Incident can only be cancelled while pending.").
				With("status", inc.Status).
				Errorf("not pending")
		}

		actorEmpID, err := s.resolveActorEmployeeID(tx, cmd.Caller.ZitadelUserID)
		if err != nil {
			return err
		}

		old := inc.Status
		inc.Status = model.IncidentStatusCancelled
		inc.UpdatedAt = now
		if err := tx.Save(inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
		}
		env, err := buildIncidentStatusChangedEnvelope(inc.ID, old, inc.Status, actorEmpID, cmd.Caller.ZitadelUserID, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.incident.v1.status_changed", env)
	})
}
