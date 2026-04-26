package request

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

		actorDisplayName, err := s.resolveActorDisplayName(tx, cmd.Caller.ZitadelUserID)
		if err != nil {
			return err
		}

		old := sr.Status
		sr.Status = newStatus
		sr.UpdatedAt = now
		if err := tx.Save(sr).Error; err != nil {
			return oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).Wrap(err)
		}
		return projector.ServiceRequestStatusChanged(
			tx, sr.ID, old, sr.Status, cmd.Caller.ZitadelUserID, actorDisplayName, now)
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
