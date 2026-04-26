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
// See: https://github.com/medincident/medincident-backend/wiki/Service-Requests#assignexecutors
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

		actorDisplayName, err := s.resolveActorDisplayName(tx, cmd.Caller.ZitadelUserID)
		if err != nil {
			return err
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
				empName := s.resolveEmployeeName(tx, e.EmployeeID)
				if err := projector.ServiceRequestExecutorRemoved(
					tx, sr.ID, e.EmployeeID, empName,
					cmd.Caller.ZitadelUserID, actorDisplayName, now,
				); err != nil {
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
			empName := s.resolveEmployeeName(tx, empID)
			if err := projector.ServiceRequestExecutorAssigned(
				tx, sr.ID, empID, empName,
				cmd.Caller.ZitadelUserID, actorDisplayName, now,
			); err != nil {
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
