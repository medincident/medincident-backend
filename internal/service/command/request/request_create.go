package request

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// CreateServiceRequestPayload is the validated client-facing payload.
type CreateServiceRequestPayload struct {
	DepartmentID        string   `validate:"required,uuid"`
	TypeID              string   `validate:"required,uuid"`
	IncidentID          *string  `validate:"omitnil,uuid"`
	Description         string   `validate:"required,no_extra_ws,min=1,max=10000"`
	ExecutorEmployeeIDs []string `validate:"required,min=1,dive,required,uuid"`
}

// CreateServiceRequestCommand = caller + payload.
type CreateServiceRequestCommand struct {
	Caller  authz.Caller
	Payload CreateServiceRequestPayload
}

// CreateServiceRequestResult is the output of ServiceRequestService.Create.
type CreateServiceRequestResult struct {
	ID uuid.UUID
}

// Create creates a new service request.
// See: https://github.com/medincident/medincident-backend/wiki/Service-Requests#createservicerequest
func (s *ServiceRequestService) Create(
	ctx context.Context,
	cmd *CreateServiceRequestCommand,
) (CreateServiceRequestResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateServiceRequestResult{}, err
	}
	deptID := uuid.MustParse(cmd.Payload.DepartmentID)
	typeID := uuid.MustParse(cmd.Payload.TypeID)

	executorIDs := make([]uuid.UUID, 0, len(cmd.Payload.ExecutorEmployeeIDs))
	for _, raw := range cmd.Payload.ExecutorEmployeeIDs {
		executorIDs = append(executorIDs, uuid.MustParse(raw))
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateServiceRequestResult{}, oops.In(scope).
			Code(ErrCodeServiceRequestIDGenerationFailed).
			Public("Failed to create service request.").Wrap(err)
	}

	var result CreateServiceRequestResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dept model.Department
		if err := tx.First(&dept, "id = ?", deptID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).Code(ErrCodeServiceRequestDeptNotFound).
					Public("Department not found.").With("department_id", deptID).Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeServiceRequestLoadFailed).Wrap(err)
		}
		var clinic model.Clinic
		if err := tx.First(&clinic, "id = ?", dept.ClinicID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeServiceRequestClinicNotFound).
				With("clinic_id", dept.ClinicID).Wrap(err)
		}
		orgID := clinic.OrganizationID

		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			privilegedActorPolicy(orgID, clinic.ID, deptID)); err != nil {
			return err
		}

		var reqType model.RequestType
		if err := tx.First(&reqType, "id = ?", typeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).Code(ErrCodeServiceRequestTypeNotFound).
					Public("Request type not found.").With("type_id", typeID).Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeServiceRequestLoadFailed).Wrap(err)
		}
		if reqType.OrganizationID != orgID {
			return oops.In(scope).Code(ErrCodeServiceRequestTypeOrgMismatch).
				Public("Request type does not belong to the department's organization.").
				With("type_id", typeID).Errorf("org mismatch")
		}
		if !reqType.IsActive {
			return oops.In(scope).Code(ErrCodeServiceRequestTypeInactive).
				Public("Request type is inactive.").With("type_id", typeID).Errorf("inactive")
		}

		incidentID := uuid.NullUUID{}
		if cmd.Payload.IncidentID != nil {
			incID := uuid.MustParse(*cmd.Payload.IncidentID)
			var inc model.Incident
			if err := tx.First(&inc, "id = ?", incID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return oops.In(scope).Code(ErrCodeServiceRequestIncidentNotFound).
						Public("Incident not found.").With("incident_id", incID).Wrap(err)
				}
				return oops.In(scope).Code(ErrCodeServiceRequestLoadFailed).Wrap(err)
			}
			if inc.OrganizationID != orgID {
				return oops.In(scope).Code(ErrCodeServiceRequestIncidentOrgMismatch).
					Public("Incident does not belong to the same organization.").
					With("incident_id", incID).Errorf("org mismatch")
			}
			incidentID = uuid.NullUUID{UUID: incID, Valid: true}
		}

		for _, empID := range executorIDs {
			var emp model.Employee
			if err := tx.First(&emp, "id = ?", empID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return oops.In(scope).Code(ErrCodeServiceRequestEmployeeNotFound).
						Public("Executor employee not found.").With("employee_id", empID).Wrap(err)
				}
				return oops.In(scope).Code(ErrCodeServiceRequestLoadFailed).Wrap(err)
			}
			if emp.DepartmentID != deptID {
				return oops.In(scope).Code(ErrCodeServiceRequestEmployeeDeptMismatch).
					Public("Executor must be an employee of the request's department.").
					With("employee_id", empID).With("department_id", deptID).Errorf("dept mismatch")
			}
		}

		authorDisplayName, err := s.resolveActorDisplayName(tx, cmd.Caller.ZitadelUserID)
		if err != nil {
			return err
		}

		now := time.Now()
		sr := model.ServiceRequest{
			ID:             id,
			OrganizationID: orgID,
			ClinicID:       clinic.ID,
			DepartmentID:   deptID,
			TypeID:         typeID,
			IncidentID:     incidentID,
			Description:    strings.TrimSpace(cmd.Payload.Description),
			Status:         model.ServiceRequestStatusCreated,
			AuthorID:       cmd.Caller.ZitadelUserID,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := tx.Create(&sr).Error; err != nil {
			return oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).
				With("service_request_id", id).Wrap(err)
		}

		for _, empID := range executorIDs {
			execID, err := uuid.NewV7()
			if err != nil {
				return oops.In(scope).Code(ErrCodeServiceRequestExecutorIDGenFailed).Wrap(err)
			}
			exec := model.ServiceRequestExecutor{
				ID:           execID,
				RequestID:    id,
				EmployeeID:   empID,
				AssignedAt:   now,
				AssignedByID: cmd.Caller.ZitadelUserID,
			}
			if err := tx.Create(&exec).Error; err != nil {
				return oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).
					With("service_request_id", id).With("employee_id", empID).Wrap(err)
			}

			empName := s.resolveEmployeeName(tx, empID)
			if err := projector.ServiceRequestExecutorAssigned(
				tx, id, empID, empName, cmd.Caller.ZitadelUserID, authorDisplayName, now,
			); err != nil {
				return err
			}
		}

		if err := projector.ServiceRequestCreated(tx, &sr, &projector.ServiceRequestAuthorSnapshot{
			DisplayName: authorDisplayName,
		}); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
