// Package request is the write-side service for service requests.
// The struct exposes all command methods on a single service type.
package request

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
)

// Error codes shared across service request command methods.
const (
	ErrCodeServiceRequestNotFound             = "service_request_not_found"
	ErrCodeServiceRequestLoadFailed           = "service_request_load_failed"
	ErrCodeServiceRequestSaveFailed           = "service_request_save_failed"
	ErrCodeServiceRequestIDGenerationFailed   = "service_request_id_generation_failed"
	ErrCodeServiceRequestDeptNotFound         = "service_request_department_not_found"
	ErrCodeServiceRequestClinicNotFound       = "service_request_clinic_not_found"
	ErrCodeServiceRequestTypeNotFound         = "service_request_type_not_found"
	ErrCodeServiceRequestTypeInactive         = "service_request_type_inactive"
	ErrCodeServiceRequestTypeOrgMismatch      = "service_request_type_org_mismatch"
	ErrCodeServiceRequestIncidentNotFound     = "service_request_incident_not_found"
	ErrCodeServiceRequestIncidentOrgMismatch  = "service_request_incident_org_mismatch"
	ErrCodeServiceRequestInvalidStatusFlow    = "service_request_invalid_status_transition"
	ErrCodeServiceRequestFrozen               = "service_request_frozen"
	ErrCodeServiceRequestEmployeeNotFound     = "service_request_employee_not_found"
	ErrCodeServiceRequestEmployeeDeptMismatch = "service_request_employee_dept_mismatch"
	ErrCodeServiceRequestExecutorIDGenFailed  = "service_request_executor_id_generation_failed"
	ErrCodeServiceRequestActorNotFound        = "service_request_actor_not_found"
)

// scope is the oops scope tag for every error emitted by this package.
const scope = "services.command.request"

// ServiceRequestService is the concrete write-side service.
type ServiceRequestService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewServiceRequestService wires the service.
func NewServiceRequestService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *ServiceRequestService {
	return &ServiceRequestService{db: db, authz: az, logger: logger}
}

// privilegedActorPolicy is the standard "any privileged role for the
// request's scope" battery, used by create, edit, assign, and
// responsible-only transitions (approve, reject, cancel).
func privilegedActorPolicy(orgID, clinicID, deptID uuid.UUID) authz.Policy {
	return authz.AnyOf(
		authz.SystemAdmin,
		authz.OrgAdminOf.Organization(orgID),
		authz.OrgHeadOf.Organization(orgID),
		authz.ClinicHeadOf.Clinic(clinicID),
		authz.DeptResponsibleOf.Department(deptID),
	)
}

// loadServiceRequest reads a service request by id; not-found is
// mapped to a public error.
func (s *ServiceRequestService) loadServiceRequest(tx *gorm.DB, id uuid.UUID) (*model.ServiceRequest, error) {
	var sr model.ServiceRequest
	if err := tx.First(&sr, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, oops.In(scope).
				Code(ErrCodeServiceRequestNotFound).
				Public("Service request not found.").
				With("service_request_id", id).
				Wrap(err)
		}
		return nil, oops.In(scope).
			Code(ErrCodeServiceRequestLoadFailed).
			With("service_request_id", id).
			Wrap(err)
	}
	return &sr, nil
}

func isTerminalStatus(status model.ServiceRequestStatus) bool {
	return status == model.ServiceRequestStatusCompleted ||
		status == model.ServiceRequestStatusCancelled
}

func (s *ServiceRequestService) isCallerExecutor(tx *gorm.DB, callerZitadelID string, requestID uuid.UUID) (bool, error) {
	var count int64
	err := tx.Raw(`
		SELECT COUNT(*)
		FROM domain.service_request_executors sre
		JOIN domain.employees e ON e.id = sre.employee_id
		WHERE sre.request_id = ? AND e.zitadel_user_id = ?`,
		requestID, callerZitadelID,
	).Scan(&count).Error
	if err != nil {
		return false, oops.In(scope).Code(ErrCodeServiceRequestLoadFailed).Wrap(err)
	}
	return count > 0, nil
}

func (s *ServiceRequestService) resolveActorDisplayName(tx *gorm.DB, callerID string) (string, error) {
	var displayName string
	if err := tx.Raw(
		`SELECT display_name FROM projections.users WHERE id = ?`, callerID,
	).Scan(&displayName).Error; err != nil {
		return "", oops.In(scope).
			Code(ErrCodeServiceRequestActorNotFound).
			Public("Actor user record is missing.").
			With("zitadel_user_id", callerID).
			Wrap(err)
	}
	if displayName == "" {
		return "", oops.In(scope).
			Code(ErrCodeServiceRequestActorNotFound).
			Public("Actor user record is missing.").
			With("zitadel_user_id", callerID).
			Errorf("display_name is empty")
	}
	return displayName, nil
}

func (s *ServiceRequestService) resolveEmployeeName(tx *gorm.DB, employeeID uuid.UUID) string {
	var displayName string
	err := tx.Raw(`
		SELECT pu.display_name
		FROM domain.employees e
		JOIN projections.users pu ON pu.id = e.zitadel_user_id
		WHERE e.id = ?`, employeeID,
	).Scan(&displayName).Error
	if err != nil || displayName == "" {
		return "unknown"
	}
	return displayName
}
