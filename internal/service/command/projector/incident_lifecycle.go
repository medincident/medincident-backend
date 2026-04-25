package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
)

// IncidentRegistrarSnapshot carries the registrar facts to freeze
// next to the incident projection. The service layer reads these
// from projections.users (display_name) and domain.employees
// (organization_id, clinic_id, department_id, position) before
// invoking the projector.
type IncidentRegistrarSnapshot struct {
	EmployeeID     uuid.UUID
	DisplayName    string
	Position       null.String
	OrganizationID uuid.UUID
	ClinicID       uuid.UUID
	DepartmentID   uuid.UUID
}

// IncidentCreated inserts the projection row plus the initial
// status history entry (old_status NULL).
func IncidentCreated(
	tx *gorm.DB,
	inc *model.Incident,
	reg *IncidentRegistrarSnapshot,
) error {
	if err := tx.Exec(`
		INSERT INTO projections.incidents (
			id, organization_id, clinic_id, department_id, category_id, type_id,
			status, priority, description, patient_original_description,
			occurred_at, registrar_employee_id, registrar_display_name,
			registrar_position, registrar_organization_id, registrar_clinic_id,
			registrar_department_id, source_patient_zitadel_user_id,
			source_buffer_id, reopened_from_incident_id, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		inc.ID, inc.OrganizationID, inc.ClinicID, inc.DepartmentID,
		inc.CategoryID, inc.TypeID, inc.Status, inc.Priority,
		inc.Description, inc.PatientOriginalDescription, inc.OccurredAt,
		reg.EmployeeID, reg.DisplayName, reg.Position,
		reg.OrganizationID, reg.ClinicID, reg.DepartmentID,
		inc.SourcePatientZitadelUserID, inc.SourceBufferID,
		inc.ReopenedFromIncidentID, inc.CreatedAt, inc.UpdatedAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert incidents projection", inc.ID)
	}

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapIncidentLifecycle(err, "generate initial status history id", inc.ID)
	}
	if err := tx.Exec(`
		INSERT INTO projections.incident_status_history (
			id, incident_id, old_status, new_status, actor_employee_id,
			actor_display_name, changed_at
		) VALUES (?, ?, NULL, ?, ?, ?, ?)`,
		histID, inc.ID, inc.Status, reg.EmployeeID, reg.DisplayName, inc.CreatedAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert initial status history", inc.ID)
	}
	return nil
}

// IncidentStatusChanged updates projections.incidents.status and
// appends a row to projections.incident_status_history. UpdatedAt is
// the change timestamp.
func IncidentStatusChanged(
	tx *gorm.DB,
	incidentID uuid.UUID,
	oldStatus, newStatus model.IncidentStatus,
	actorEmployeeID uuid.UUID,
	actorDisplayName string,
	changedAt time.Time,
) error {
	if err := tx.Exec(`
		UPDATE projections.incidents
		   SET status = ?, updated_at = ?
		 WHERE id = ?`,
		newStatus, changedAt, incidentID,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "update status projection", incidentID)
	}
	histID, err := uuid.NewV7()
	if err != nil {
		return wrapIncidentLifecycle(err, "generate status history id", incidentID)
	}
	if err := tx.Exec(`
		INSERT INTO projections.incident_status_history (
			id, incident_id, old_status, new_status,
			actor_employee_id, actor_display_name, changed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		histID, incidentID, oldStatus, newStatus,
		actorEmployeeID, actorDisplayName, changedAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert status history", incidentID)
	}
	return nil
}

// IncidentPriorityChanged updates projections.incidents.priority and
// appends a row to projections.incident_priority_history.
func IncidentPriorityChanged(
	tx *gorm.DB,
	incidentID uuid.UUID,
	oldPriority, newPriority model.IncidentPriority,
	actorEmployeeID uuid.UUID,
	actorDisplayName string,
	changedAt time.Time,
) error {
	if err := tx.Exec(`
		UPDATE projections.incidents
		   SET priority = ?, updated_at = ?
		 WHERE id = ?`,
		newPriority, changedAt, incidentID,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "update priority projection", incidentID)
	}
	histID, err := uuid.NewV7()
	if err != nil {
		return wrapIncidentLifecycle(err, "generate priority history id", incidentID)
	}
	if err := tx.Exec(`
		INSERT INTO projections.incident_priority_history (
			id, incident_id, old_priority, new_priority,
			actor_employee_id, actor_display_name, changed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		histID, incidentID, oldPriority, newPriority,
		actorEmployeeID, actorDisplayName, changedAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert priority history", incidentID)
	}
	return nil
}

// IncidentDescriptionUpdated updates description + updated_at on the
// projection. No history row — description edits are not auditable.
func IncidentDescriptionUpdated(
	tx *gorm.DB,
	incidentID uuid.UUID,
	description null.String,
	updatedAt time.Time,
) error {
	if err := tx.Exec(`
		UPDATE projections.incidents
		   SET description = ?, updated_at = ?
		 WHERE id = ?`,
		description, updatedAt, incidentID,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "update description projection", incidentID)
	}
	return nil
}

// BufferCreated inserts a buffer projection row.
func BufferCreated(tx *gorm.DB, b *model.PatientIncidentBuffer) error {
	if err := tx.Exec(`
		INSERT INTO projections.patient_incident_buffer (
			id, organization_id, patient_zitadel_user_id,
			category_id, type_id, description, occurred_at,
			status, published_incident_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		b.ID, b.OrganizationID, b.PatientZitadelUserID,
		b.CategoryID, b.TypeID, b.Description, b.OccurredAt,
		b.Status, b.PublishedIncidentID, b.CreatedAt, b.UpdatedAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert buffer projection", b.ID)
	}
	return nil
}

// BufferUpdated mirrors a buffer mutation (edit / cancel / reject / publish).
func BufferUpdated(tx *gorm.DB, b *model.PatientIncidentBuffer) error {
	if err := tx.Exec(`
		UPDATE projections.patient_incident_buffer
		   SET category_id = ?, type_id = ?, description = ?, occurred_at = ?,
		       status = ?, published_incident_id = ?, updated_at = ?
		 WHERE id = ?`,
		b.CategoryID, b.TypeID, b.Description, b.OccurredAt,
		b.Status, b.PublishedIncidentID, b.UpdatedAt, b.ID,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "update buffer projection", b.ID)
	}
	return nil
}

func wrapIncidentLifecycle(err error, action string, id uuid.UUID) error {
	return oops.In("projector.incident.lifecycle").
		Code(ErrCodeIncidentProjectionFailed).
		With("action", action).
		With("incident_id", id).
		Wrap(err)
}
