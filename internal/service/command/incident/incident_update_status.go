package incident

import (
	"context"
	"errors"
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
		actorEmpID, err := s.resolveActorEmployeeID(tx, cmd.Caller.ZitadelUserID)
		if err != nil {
			return err
		}

		old := inc.Status
		inc.Status = newStatus
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

// resolveActorEmployeeID returns the caller's employee_id (in any org).
// The returned NullUUID is invalid when the caller has no employee row
// (e.g. SystemAdmin acting outside any org) so history tables get a
// proper SQL NULL rather than the zero UUID. Display name is resolved
// query-side from the event's actor_zitadel_user_id field.
func (s *IncidentService) resolveActorEmployeeID(tx *gorm.DB, callerID string) (uuid.NullUUID, error) {
	var emp model.Employee
	empID := uuid.NullUUID{}
	if err := tx.Where("zitadel_user_id = ?", callerID).Limit(1).First(&emp).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.NullUUID{}, oops.In(scope).
				Code(ErrCodeIncidentEmployeeNotFound).
				With("zitadel_user_id", callerID).Wrap(err)
		}
	} else {
		empID = uuid.NullUUID{UUID: emp.ID, Valid: true}
	}
	return empID, nil
}

func buildIncidentStatusChangedEnvelope(
	incidentID uuid.UUID,
	oldStatus, newStatus model.IncidentStatus,
	actorEmployeeID uuid.NullUUID,
	actorZitadelUserID string,
	changedAt time.Time,
) (*eventv1.Envelope, error) {
	msg := &incidentv1.IncidentStatusChanged{
		IncidentId:         incidentID.String(),
		OldStatus:          string(oldStatus),
		NewStatus:          string(newStatus),
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
