package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	incidentv1 "github.com/medincident/medincident-backend/pkg/event/incident/v1"
)

// IncidentCreated inserts projections.incidents + the initial
// projections.incident_status_history entry (old_status NULL).
//
// See: docs/services/incident/Incidents.md
func IncidentCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *incidentv1.IncidentCreated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	orgID, err := parseUUID(ev.GetOrganizationId(), "organization_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	clinicID, err := parseUUID(ev.GetClinicId(), "clinic_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	deptID, err := parseUUID(ev.GetDepartmentId(), "department_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	catID, err := parseUUID(ev.GetCategoryId(), "category_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	typeID, err := parseUUID(ev.GetTypeId(), "type_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	registrarEmpID, err := parseUUID(ev.GetRegistrarEmployeeId(), "registrar_employee_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	createdAt := ev.GetCreatedAt().AsTime()
	occAt := ev.GetOccurredAt().AsTime()

	// Resolve registrar display name from the identity projection.
	registrarDisplayName := lookupUserDisplayName(tx, ev.GetRegistrarZitadelUserId())

	var desc null.String
	if sv := ev.GetDescription(); sv != nil {
		desc = null.StringFrom(sv.GetValue())
	}
	var patientOrigDesc null.String
	if sv := ev.GetPatientOriginalDescription(); sv != nil {
		patientOrigDesc = null.StringFrom(sv.GetValue())
	}
	var srcPatientUserID null.String
	if sv := ev.GetSourcePatientZitadelUserId(); sv != nil {
		srcPatientUserID = null.StringFrom(sv.GetValue())
	}
	var srcBufID *uuid.UUID
	if sv := ev.GetSourceBufferId(); sv != nil {
		p, parseErr := parseUUID(sv.GetValue(), "source_buffer_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
		if parseErr != nil {
			return parseErr
		}
		srcBufID = &p
	}
	var reopenedFromID *uuid.UUID
	if sv := ev.GetReopenedFromIncidentId(); sv != nil {
		p, parseErr := parseUUID(sv.GetValue(), "reopened_from_incident_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
		if parseErr != nil {
			return parseErr
		}
		reopenedFromID = &p
	}

	if err := tx.Exec(`
		INSERT INTO projections.incidents (
			id, organization_id, clinic_id, department_id, category_id, type_id,
			status, priority, description, patient_original_description,
			occurred_at, registrar_employee_id, registrar_display_name,
			source_patient_zitadel_user_id, source_buffer_id,
			reopened_from_incident_id, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		) ON CONFLICT DO NOTHING`,
		id, orgID, clinicID, deptID, catID, typeID,
		ev.GetStatus(), ev.GetPriority(), desc, patientOrigDesc,
		occAt, registrarEmpID, registrarDisplayName,
		srcPatientUserID, srcBufID,
		reopenedFromID, createdAt, createdAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert incidents", id)
	}

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapIncidentLifecycle(err, "generate initial status history id", id)
	}
	if err := tx.Exec(`
		INSERT INTO projections.incident_status_history (
			id, incident_id, old_status, new_status, actor_employee_id,
			actor_display_name, changed_at
		) VALUES (?, ?, NULL, ?, ?, ?, ?)
		ON CONFLICT DO NOTHING`,
		histID, id, ev.GetStatus(), registrarEmpID, registrarDisplayName, createdAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert initial status history", id)
	}
	return nil
}

// IncidentStatusChanged updates projections.incidents.status and appends
// a row to projections.incident_status_history.
//
// See: docs/services/incident/Incidents.md
func IncidentStatusChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *incidentv1.IncidentStatusChanged,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	changedAt := ev.GetChangedAt().AsTime()

	actorDisplayName := lookupUserDisplayName(tx, ev.GetActorZitadelUserId())

	var actorEmpID *uuid.UUID
	if sv := ev.GetActorEmployeeId(); sv != nil {
		p, parseErr := parseUUID(sv.GetValue(), "actor_employee_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
		if parseErr != nil {
			return parseErr
		}
		actorEmpID = &p
	}

	if err := tx.Exec(`
		UPDATE projections.incidents
		   SET status = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetNewStatus(), changedAt, id,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "update status", id)
	}

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapIncidentLifecycle(err, "generate status history id", id)
	}
	if err := tx.Exec(`
		INSERT INTO projections.incident_status_history (
			id, incident_id, old_status, new_status,
			actor_employee_id, actor_display_name, changed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		histID, id, ev.GetOldStatus(), ev.GetNewStatus(),
		actorEmpID, actorDisplayName, changedAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert status history", id)
	}
	return nil
}

// IncidentPriorityChanged updates projections.incidents.priority and appends
// a row to projections.incident_priority_history.
//
// See: docs/services/incident/Incidents.md
func IncidentPriorityChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *incidentv1.IncidentPriorityChanged,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	changedAt := ev.GetChangedAt().AsTime()

	actorDisplayName := lookupUserDisplayName(tx, ev.GetActorZitadelUserId())

	var actorEmpID *uuid.UUID
	if sv := ev.GetActorEmployeeId(); sv != nil {
		p, parseErr := parseUUID(sv.GetValue(), "actor_employee_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
		if parseErr != nil {
			return parseErr
		}
		actorEmpID = &p
	}

	if err := tx.Exec(`
		UPDATE projections.incidents
		   SET priority = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetNewPriority(), changedAt, id,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "update priority", id)
	}

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapIncidentLifecycle(err, "generate priority history id", id)
	}
	if err := tx.Exec(`
		INSERT INTO projections.incident_priority_history (
			id, incident_id, old_priority, new_priority,
			actor_employee_id, actor_display_name, changed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		histID, id, ev.GetOldPriority(), ev.GetNewPriority(),
		actorEmpID, actorDisplayName, changedAt,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "insert priority history", id)
	}
	return nil
}

// IncidentDescriptionUpdated updates description + updated_at on the projection.
//
// See: docs/services/incident/Incidents.md
func IncidentDescriptionUpdated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *incidentv1.IncidentDescriptionUpdated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_lifecycle", ErrCodeIncidentLifecycleProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	var desc null.String
	if sv := ev.GetDescription(); sv != nil {
		desc = null.StringFrom(sv.GetValue())
	}

	if err := tx.Exec(`
		UPDATE projections.incidents
		   SET description = ?, updated_at = ?
		 WHERE id = ?`,
		desc, updatedAt, id,
	).Error; err != nil {
		return wrapIncidentLifecycle(err, "update description", id)
	}
	return nil
}

func wrapIncidentLifecycle(err error, action string, id uuid.UUID) error {
	return oops.In("projector.incident.lifecycle").
		Code(ErrCodeIncidentLifecycleProjectionFailed).
		With("action", action).
		With("incident_id", id).
		Wrap(err)
}

// ── Projectors forwarding methods ────────────────────────────────────────────

// IncidentCreated forwards to the package-level projector function.
func (p *Projectors) IncidentCreated(tx *gorm.DB, id string, t time.Time, ev *incidentv1.IncidentCreated) error {
	return IncidentCreated(tx, id, t, ev)
}

// IncidentStatusChanged forwards to the package-level projector function.
func (p *Projectors) IncidentStatusChanged(tx *gorm.DB, id string, t time.Time, ev *incidentv1.IncidentStatusChanged) error {
	return IncidentStatusChanged(tx, id, t, ev)
}

// IncidentPriorityChanged forwards to the package-level projector function.
func (p *Projectors) IncidentPriorityChanged(tx *gorm.DB, id string, t time.Time, ev *incidentv1.IncidentPriorityChanged) error {
	return IncidentPriorityChanged(tx, id, t, ev)
}

// IncidentDescriptionUpdated forwards to the package-level projector function.
func (p *Projectors) IncidentDescriptionUpdated(tx *gorm.DB, id string, t time.Time, ev *incidentv1.IncidentDescriptionUpdated) error {
	return IncidentDescriptionUpdated(tx, id, t, ev)
}
