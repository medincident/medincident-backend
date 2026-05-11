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

// AssignExecutorsPayload carries the new executor set.
type AssignExecutorsPayload struct {
	ServiceRequestID    string   `validate:"required,uuid"`
	ExecutorEmployeeIDs []string `validate:"required,min=1,dive,required,uuid"`
}

// AssignExecutorsCommand = caller + payload.
type AssignExecutorsCommand struct {
	Caller  authz.Caller
	Payload AssignExecutorsPayload
}

// AssignExecutors replaces the executor set of a service request using
// a diff-based approach: remove absent, add new.
// See: docs/services/request/Requests.md
func (s *ServiceRequestService) AssignExecutors(
	ctx context.Context,
	cmd AssignExecutorsCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ServiceRequestID)
	newIDs := make([]uuid.UUID, 0, len(cmd.Payload.ExecutorEmployeeIDs))
	for _, raw := range cmd.Payload.ExecutorEmployeeIDs {
		newIDs = append(newIDs, uuid.MustParse(raw))
	}
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
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			privilegedActorPolicy(sr.OrganizationID, sr.ClinicID, sr.DepartmentID)); err != nil {
			return err
		}

		for _, empID := range newIDs {
			var emp model.Employee
			if err := tx.First(&emp, "id = ?", empID).Error; err != nil {
				return oops.In(scope).Code(ErrCodeServiceRequestEmployeeNotFound).
					Public("Executor employee not found.").With("employee_id", empID).Wrap(err)
			}
			if emp.DepartmentID != sr.DepartmentID {
				return oops.In(scope).Code(ErrCodeServiceRequestEmployeeDeptMismatch).
					Public("Executor must be an employee of the request's department.").
					With("employee_id", empID).With("department_id", sr.DepartmentID).
					Errorf("dept mismatch")
			}
		}

		var existing []model.ServiceRequestExecutor
		if err := tx.Where("request_id = ?", sr.ID).Find(&existing).Error; err != nil {
			return oops.In(scope).Code(ErrCodeServiceRequestLoadFailed).Wrap(err)
		}

		existingMap := make(map[uuid.UUID]model.ServiceRequestExecutor, len(existing))
		for _, e := range existing {
			existingMap[e.EmployeeID] = e
		}
		newSet := make(map[uuid.UUID]bool, len(newIDs))
		for _, eid := range newIDs {
			newSet[eid] = true
		}

		// Remove executors not in the new set.
		for _, e := range existing {
			if !newSet[e.EmployeeID] {
				if err := tx.Delete(&model.ServiceRequestExecutor{}, "id = ?", e.ID).Error; err != nil {
					return oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).Wrap(err)
				}
				env, err := buildServiceRequestExecutorRemovedEnvelope(sr.ID, e.EmployeeID, cmd.Caller.ZitadelUserID, now)
				if err != nil {
					return err
				}
				if err := outbox.Append(tx, "medincident.event.service_request.v1.executor_removed", env); err != nil {
					return err
				}
			}
		}

		// Add new executors not already existing.
		for _, empID := range newIDs {
			if _, exists := existingMap[empID]; exists {
				continue
			}
			execID, err := uuid.NewV7()
			if err != nil {
				return oops.In(scope).Code(ErrCodeServiceRequestExecutorIDGenFailed).Wrap(err)
			}
			exec := model.ServiceRequestExecutor{
				ID:           execID,
				RequestID:    sr.ID,
				EmployeeID:   empID,
				AssignedAt:   now,
				AssignedByID: cmd.Caller.ZitadelUserID,
			}
			if err := tx.Create(&exec).Error; err != nil {
				return oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).Wrap(err)
			}
			env, err := buildServiceRequestExecutorAssignedEnvelope(sr.ID, empID, cmd.Caller.ZitadelUserID, now)
			if err != nil {
				return err
			}
			if err := outbox.Append(tx, "medincident.event.service_request.v1.executor_assigned", env); err != nil {
				return err
			}
		}

		sr.UpdatedAt = now
		if err := tx.Save(sr).Error; err != nil {
			return oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).Wrap(err)
		}
		return nil
	})
}

func buildServiceRequestExecutorRemovedEnvelope(requestID, employeeID uuid.UUID, actorZitadelUserID string, changedAt time.Time) (*eventv1.Envelope, error) {
	msg := &servicerequestv1.ServiceRequestExecutorRemoved{
		RequestId:  requestID.String(),
		EmployeeId: employeeID.String(),
		ActorId:    actorZitadelUserID,
		ChangedAt:  timestamppb.New(changedAt),
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
