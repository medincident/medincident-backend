package incident

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	incidentv1 "github.com/medincident/medincident-backend/pkg/event/incident/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

type UpdateIncidentPriorityPayload struct {
	IncidentID string `validate:"required,uuid"`
	Priority   string `validate:"required,oneof=low normal high critical"`
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
	ctx context.Context, cmd UpdateIncidentPriorityCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.IncidentID)
	priority := model.IncidentPriority(cmd.Payload.Priority)
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
		if inc.Priority == priority {
			return nil
		}
		actorEmpID, err := s.resolveActorEmployeeID(tx, cmd.Caller.ZitadelUserID)
		if err != nil {
			return err
		}
		old := inc.Priority
		inc.Priority = priority
		inc.UpdatedAt = now
		if err := tx.Save(inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
		}
		env, err := buildIncidentPriorityChangedEnvelope(inc.ID, old, inc.Priority, actorEmpID, cmd.Caller.ZitadelUserID, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.incident.v1.priority_changed", env)
	})
}

func buildIncidentPriorityChangedEnvelope(
	incidentID uuid.UUID,
	oldPriority, newPriority model.IncidentPriority,
	actorEmployeeID uuid.NullUUID,
	actorZitadelUserID string,
	changedAt time.Time,
) (*eventv1.Envelope, error) {
	msg := &incidentv1.IncidentPriorityChanged{
		IncidentId:         incidentID.String(),
		OldPriority:        string(oldPriority),
		NewPriority:        string(newPriority),
		ActorZitadelUserId: actorZitadelUserID,
		ChangedAt:          timestamppb.New(changedAt),
	}
	if actorEmployeeID.Valid {
		msg.ActorEmployeeId = wrapperspb.String(actorEmployeeID.UUID.String())
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(changedAt),
		AggregateType: "incident",
		AggregateId:   incidentID.String(),
		Payload:       payload,
	}, nil
}
