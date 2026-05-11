package request

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	servicerequestv1 "github.com/medincident/medincident-backend/pkg/event/service_request/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// UpdateServiceRequestStatusPayload carries the new status.
type UpdateServiceRequestStatusPayload struct {
	ServiceRequestID string `validate:"required,uuid"`
	NewStatus        string `validate:"required,oneof=in_work on_hold pending_review completed cancelled"`
}

// UpdateServiceRequestStatusCommand = caller + payload.
type UpdateServiceRequestStatusCommand struct {
	Caller  authz.Caller
	Payload UpdateServiceRequestStatusPayload
}

// UpdateStatus transitions a service request to a new status.
// See: docs/services/request/Requests.md
func (s *ServiceRequestService) UpdateStatus(
	ctx context.Context,
	cmd UpdateServiceRequestStatusCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ServiceRequestID)
	newStatus := model.ServiceRequestStatus(cmd.Payload.NewStatus)
	now := time.Now()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sr, err := s.loadServiceRequest(tx, id)
		if err != nil {
			return err
		}
		if isTerminalStatus(sr.Status) {
			return oops.In(scope).Code(ErrCodeServiceRequestFrozen).
				Public("Cannot modify a completed or cancelled request.").
				With("service_request_id", id).With("status", sr.Status).
				Errorf("frozen")
		}

		isExecutor, err := s.isCallerExecutor(tx, cmd.Caller.ZitadelUserID, sr.ID)
		if err != nil {
			return err
		}

		if err := s.validateStatusTransition(ctx, sr, newStatus, isExecutor, cmd.Caller.ZitadelUserID); err != nil {
			return err
		}

		old := sr.Status
		sr.Status = newStatus
		sr.UpdatedAt = now
		if err := tx.Save(sr).Error; err != nil {
			return oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).Wrap(err)
		}
		env, err := buildServiceRequestStatusChangedEnvelope(sr.ID, old, sr.Status, cmd.Caller.ZitadelUserID, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.service_request.v1.status_changed", env)
	})
}

func (s *ServiceRequestService) validateStatusTransition(
	ctx context.Context,
	sr *model.ServiceRequest,
	newStatus model.ServiceRequestStatus,
	isExecutor bool,
	callerID string,
) error {
	from := sr.Status
	to := newStatus

	executorTransitions := map[model.ServiceRequestStatus]map[model.ServiceRequestStatus]bool{
		model.ServiceRequestStatusCreated: {model.ServiceRequestStatusInWork: true},
		model.ServiceRequestStatusInWork:  {model.ServiceRequestStatusOnHold: true, model.ServiceRequestStatusPendingReview: true},
		model.ServiceRequestStatusOnHold:  {model.ServiceRequestStatusInWork: true},
	}

	responsibleTransitions := map[model.ServiceRequestStatus]map[model.ServiceRequestStatus]bool{
		model.ServiceRequestStatusPendingReview: {model.ServiceRequestStatusCompleted: true, model.ServiceRequestStatusInWork: true},
	}

	if isExecutor {
		if targets, ok := executorTransitions[from]; ok && targets[to] {
			return nil
		}
	}

	if targets, ok := responsibleTransitions[from]; ok && targets[to] {
		return s.authz.Require(ctx, callerID,
			privilegedActorPolicy(sr.OrganizationID, sr.ClinicID, sr.DepartmentID))
	}

	if to == model.ServiceRequestStatusCancelled && !isTerminalStatus(from) {
		return s.authz.Require(ctx, callerID,
			privilegedActorPolicy(sr.OrganizationID, sr.ClinicID, sr.DepartmentID))
	}

	return oops.In(scope).
		Code(ErrCodeServiceRequestInvalidStatusFlow).
		Public("This status transition is not allowed.").
		With("from", from).With("to", to).
		Errorf("invalid transition")
}

func buildServiceRequestStatusChangedEnvelope(requestID uuid.UUID, oldStatus, newStatus model.ServiceRequestStatus, actorZitadelUserID string, changedAt time.Time) (*eventv1.Envelope, error) {
	msg := &servicerequestv1.ServiceRequestStatusChanged{
		RequestId: requestID.String(),
		OldStatus: string(oldStatus),
		NewStatus: string(newStatus),
		ActorId:   actorZitadelUserID,
		ChangedAt: timestamppb.New(changedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(changedAt),
		AggregateType: "service_request",
		AggregateId:   requestID.String(),
		Payload:       payload,
	}, nil
}
