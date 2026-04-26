// Package incident is the write-side service for incidents created
// directly by employees and incidents materialised from the patient
// buffer. The struct exposes all command methods on a single service
// type — Create, Cancel, UpdateStatus, UpdatePriority,
// UpdateDescription, Reopen.
//
// See: docs/services/incident/Incidents.md
package incident

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
)

// Error codes shared across incident command methods.
const (
	ErrCodeIncidentNotFound              = "incident_not_found"
	ErrCodeIncidentLoadFailed            = "incident_load_failed"
	ErrCodeIncidentSaveFailed            = "incident_save_failed"
	ErrCodeIncidentIDGenerationFailed    = "incident_id_generation_failed"
	ErrCodeIncidentDeptNotFound          = "incident_department_not_found"
	ErrCodeIncidentClinicNotFound        = "incident_clinic_not_found"
	ErrCodeIncidentCategoryNotFound      = "incident_category_not_found"
	ErrCodeIncidentTypeNotFound          = "incident_type_not_found"
	ErrCodeIncidentCategoryInactive      = "incident_category_inactive"
	ErrCodeIncidentTypeInactive          = "incident_type_inactive"
	ErrCodeIncidentTypeOrgMismatch       = "incident_type_organization_mismatch"
	ErrCodeIncidentTypeCategoryMismatch  = "incident_type_category_mismatch"
	ErrCodeIncidentEmployeeNotFound      = "incident_employee_not_found"
	ErrCodeIncidentRegistrarUserNotFound = "incident_registrar_user_not_found"
	ErrCodeIncidentOccurredAtInvalid     = "incident_occurred_at_invalid"
	ErrCodeIncidentOccurredAtFuture      = "incident_occurred_at_in_future"
	ErrCodeIncidentOccurredAtTooOld      = "incident_occurred_at_too_old"
	ErrCodeIncidentInvalidStatusFlow     = "incident_invalid_status_transition"
	ErrCodeIncidentFrozen                = "incident_frozen"
	ErrCodeIncidentNotCancellable        = "incident_not_cancellable"
	ErrCodeIncidentNotReopenable         = "incident_not_reopenable"
)

// incidentMaxOccurredAtAge bounds how far in the past an employee
// may file an incident.
const incidentMaxOccurredAtAge = 48 * time.Hour

// scope is the oops scope tag for every error emitted by this package.
const scope = "services.command.incident"

// IncidentService is the concrete write-side service.
type IncidentService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewIncidentService wires the service.
func NewIncidentService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *IncidentService {
	return &IncidentService{db: db, authz: az, logger: logger}
}

// privilegedActorPolicy is the standard "any privileged role for the
// incident's scope" battery, used by every action that mutates an
// incident other than registrar-only cancellation.
func privilegedActorPolicy(orgID, clinicID, deptID uuid.UUID) authz.Policy {
	return authz.AnyOf(
		authz.SystemAdmin,
		authz.OrgAdminOf.Organization(orgID),
		authz.OrgHeadOf.Organization(orgID),
		authz.OrgDispatcherOf.Organization(orgID),
		authz.ClinicHeadOf.Clinic(clinicID),
		authz.DeptResponsibleOf.Department(deptID),
	)
}

// loadIncident reads an incident by id; not-found is mapped to a
// public error.
func (s *IncidentService) loadIncident(tx *gorm.DB, id uuid.UUID) (*model.Incident, error) {
	var inc model.Incident
	if err := tx.First(&inc, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, oops.In(scope).
				Code(ErrCodeIncidentNotFound).
				Public("Incident not found.").
				With("incident_id", id).
				Wrap(err)
		}
		return nil, oops.In(scope).
			Code(ErrCodeIncidentLoadFailed).
			With("incident_id", id).
			Wrap(err)
	}
	return &inc, nil
}

// callerIsActiveRegistrar returns true iff the caller's zitadel id
// matches the registrar's employee row. Runs against the supplied
// gorm transaction (NOT the bare s.db) so the check participates in
// the surrounding transaction's isolation — a registrar deleted in a
// concurrent tx is correctly invisible, and uncommitted test
// fixtures within the same tx are visible.
func (s *IncidentService) callerIsActiveRegistrar(
	tx *gorm.DB, callerID string, registrarEmployeeID uuid.UUID,
) (bool, error) {
	var count int64
	if err := tx.
		Model(&model.Employee{}).
		Where("id = ? AND zitadel_user_id = ?", registrarEmployeeID, callerID).
		Count(&count).Error; err != nil {
		return false, oops.In(scope).
			Code(ErrCodeIncidentEmployeeNotFound).
			With("registrar_employee_id", registrarEmployeeID).
			Wrap(err)
	}
	return count > 0, nil
}
