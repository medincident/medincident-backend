package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	servicerequestv1 "github.com/medincident/medincident-backend/pkg/event/service_request/v1"
)

// ServiceRequestCreated inserts projections.service_requests + the initial
// projections.service_request_status_history entry.
//
// author_display_name is resolved from projections.users so the command
// side does NOT need to cross-read the query DB.
//
// See: docs/services/request/Requests.md
func ServiceRequestCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *servicerequestv1.ServiceRequestCreated,
) error {
	id := uuid.MustParse(aggregateID)
	orgID := uuid.MustParse(ev.GetOrganizationId())
	clinicID := uuid.MustParse(ev.GetClinicId())
	deptID := uuid.MustParse(ev.GetDepartmentId())
	typeID := uuid.MustParse(ev.GetTypeId())
	createdAt := ev.GetCreatedAt().AsTime()

	authorDisplayName := lookupUserDisplayName(tx, ev.GetAuthorId())

	var incidentID *uuid.UUID
	if sv := ev.GetIncidentId(); sv != nil {
		p := uuid.MustParse(sv.GetValue())
		incidentID = &p
	}

	if err := tx.Exec(`
		INSERT INTO projections.service_requests (
			id, organization_id, clinic_id, department_id, type_id, incident_id,
			description, status, author_id, author_display_name, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT DO NOTHING`,
		id, orgID, clinicID, deptID, typeID, incidentID,
		ev.GetDescription(), ev.GetStatus(),
		ev.GetAuthorId(), authorDisplayName, createdAt, createdAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert", id)
	}

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapServiceRequest(err, "generate initial status history id", id)
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_request_status_history (
			id, request_id, old_status, new_status, actor_id, actor_name, changed_at
		) VALUES (?, ?, NULL, ?, ?, ?, ?)
		ON CONFLICT DO NOTHING`,
		histID, id, ev.GetStatus(), ev.GetAuthorId(), authorDisplayName, createdAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert initial status history", id)
	}
	return nil
}

// ServiceRequestStatusChanged updates status + appends a history row.
//
// See: docs/services/request/Requests.md
func ServiceRequestStatusChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *servicerequestv1.ServiceRequestStatusChanged,
) error {
	id := uuid.MustParse(aggregateID)
	changedAt := ev.GetChangedAt().AsTime()

	actorDisplayName := lookupUserDisplayName(tx, ev.GetActorId())

	if err := tx.Exec(`
		UPDATE projections.service_requests
		   SET status = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetNewStatus(), changedAt, id,
	).Error; err != nil {
		return wrapServiceRequest(err, "update status", id)
	}

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapServiceRequest(err, "generate status history id", id)
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_request_status_history (
			id, request_id, old_status, new_status, actor_id, actor_name, changed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		histID, id, ev.GetOldStatus(), ev.GetNewStatus(),
		ev.GetActorId(), actorDisplayName, changedAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert status history", id)
	}
	return nil
}

// ServiceRequestDescriptionUpdated updates description and updated_at.
//
// See: docs/services/request/Requests.md
func ServiceRequestDescriptionUpdated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *servicerequestv1.ServiceRequestDescriptionUpdated,
) error {
	id := uuid.MustParse(aggregateID)
	if err := tx.Exec(`
		UPDATE projections.service_requests
		   SET description = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetDescription(), ev.GetUpdatedAt().AsTime(), id,
	).Error; err != nil {
		return wrapServiceRequest(err, "update description", id)
	}
	return nil
}

// ServiceRequestExecutorAssigned appends an 'assigned' row to the executor history.
//
// See: docs/services/request/Requests.md
func ServiceRequestExecutorAssigned(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *servicerequestv1.ServiceRequestExecutorAssigned,
) error {
	id := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	changedAt := ev.GetChangedAt().AsTime()

	actorDisplayName := lookupUserDisplayName(tx, ev.GetActorId())
	empDisplayName := lookupEmployeeDisplayName(tx, empID)

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapServiceRequest(err, "generate executor history id", id)
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_request_executor_history (
			id, request_id, action, employee_id, employee_name,
			actor_id, actor_name, changed_at
		) VALUES (?, ?, 'assigned', ?, ?, ?, ?, ?)`,
		histID, id, empID, empDisplayName, ev.GetActorId(), actorDisplayName, changedAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert executor assigned history", id)
	}
	return nil
}

// ServiceRequestExecutorRemoved appends a 'removed' row to the executor history.
//
// See: docs/services/request/Requests.md
func ServiceRequestExecutorRemoved(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *servicerequestv1.ServiceRequestExecutorRemoved,
) error {
	id := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	changedAt := ev.GetChangedAt().AsTime()

	actorDisplayName := lookupUserDisplayName(tx, ev.GetActorId())
	empDisplayName := lookupEmployeeDisplayName(tx, empID)

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapServiceRequest(err, "generate executor removal history id", id)
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_request_executor_history (
			id, request_id, action, employee_id, employee_name,
			actor_id, actor_name, changed_at
		) VALUES (?, ?, 'removed', ?, ?, ?, ?, ?)`,
		histID, id, empID, empDisplayName, ev.GetActorId(), actorDisplayName, changedAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert executor removed history", id)
	}
	return nil
}

func wrapServiceRequest(err error, action string, id uuid.UUID) error {
	return oops.In("projector.service_request").
		Code(ErrCodeServiceRequestProjectionFailed).
		With("action", action).
		With("service_request_id", id).
		Wrap(err)
}

// lookupEmployeeDisplayName looks up the display_name for an employee
// from projections.employee_cards, returning empty string when absent.
func lookupEmployeeDisplayName(tx *gorm.DB, empID uuid.UUID) string {
	var name string
	_ = tx.Raw(
		`SELECT display_name FROM projections.employee_cards WHERE employee_id = ?`, empID,
	).Row().Scan(&name)
	return name
}

// ── Projectors forwarding methods ────────────────────────────────────────────

// ServiceRequestCreated projects a ServiceRequestCreated event.
//
// See: docs/services/request/Requests.md
func (p *Projectors) ServiceRequestCreated(tx *gorm.DB, id string, t time.Time, ev *servicerequestv1.ServiceRequestCreated) error {
	return ServiceRequestCreated(tx, id, t, ev)
}

// ServiceRequestStatusChanged projects a ServiceRequestStatusChanged event.
//
// See: docs/services/request/Requests.md
func (p *Projectors) ServiceRequestStatusChanged(tx *gorm.DB, id string, t time.Time, ev *servicerequestv1.ServiceRequestStatusChanged) error {
	return ServiceRequestStatusChanged(tx, id, t, ev)
}

// ServiceRequestDescriptionUpdated projects a ServiceRequestDescriptionUpdated event.
//
// See: docs/services/request/Requests.md
func (p *Projectors) ServiceRequestDescriptionUpdated(tx *gorm.DB, id string, t time.Time, ev *servicerequestv1.ServiceRequestDescriptionUpdated) error {
	return ServiceRequestDescriptionUpdated(tx, id, t, ev)
}

// ServiceRequestExecutorAssigned projects a ServiceRequestExecutorAssigned event.
//
// See: docs/services/request/Requests.md
func (p *Projectors) ServiceRequestExecutorAssigned(tx *gorm.DB, id string, t time.Time, ev *servicerequestv1.ServiceRequestExecutorAssigned) error {
	return ServiceRequestExecutorAssigned(tx, id, t, ev)
}

// ServiceRequestExecutorRemoved projects a ServiceRequestExecutorRemoved event.
//
// See: docs/services/request/Requests.md
func (p *Projectors) ServiceRequestExecutorRemoved(tx *gorm.DB, id string, t time.Time, ev *servicerequestv1.ServiceRequestExecutorRemoved) error {
	return ServiceRequestExecutorRemoved(tx, id, t, ev)
}
