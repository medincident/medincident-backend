package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
)

type ServiceRequestAuthorSnapshot struct {
	DisplayName string
}

func ServiceRequestCreated(
	tx *gorm.DB,
	sr *model.ServiceRequest,
	author *ServiceRequestAuthorSnapshot,
) error {
	var incidentID *uuid.UUID
	if sr.IncidentID.Valid {
		p := sr.IncidentID.UUID
		incidentID = &p
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_requests (
			id, organization_id, clinic_id, department_id, type_id, incident_id,
			description, status, author_id, author_display_name, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sr.ID, sr.OrganizationID, sr.ClinicID, sr.DepartmentID,
		sr.TypeID, incidentID, sr.Description, sr.Status,
		sr.AuthorID, author.DisplayName, sr.CreatedAt, sr.UpdatedAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert projection", sr.ID)
	}

	histID, err := uuid.NewV7()
	if err != nil {
		return wrapServiceRequest(err, "generate initial status history id", sr.ID)
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_request_status_history (
			id, request_id, old_status, new_status, actor_id, actor_name, changed_at
		) VALUES (?, ?, NULL, ?, ?, ?, ?)`,
		histID, sr.ID, sr.Status, sr.AuthorID, author.DisplayName, sr.CreatedAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert initial status history", sr.ID)
	}
	return nil
}

func ServiceRequestStatusChanged(
	tx *gorm.DB,
	requestID uuid.UUID,
	oldStatus, newStatus model.ServiceRequestStatus,
	actorID string,
	actorDisplayName string,
	changedAt time.Time,
) error {
	if err := tx.Exec(`
		UPDATE projections.service_requests
		   SET status = ?, updated_at = ?
		 WHERE id = ?`,
		newStatus, changedAt, requestID,
	).Error; err != nil {
		return wrapServiceRequest(err, "update status projection", requestID)
	}
	histID, err := uuid.NewV7()
	if err != nil {
		return wrapServiceRequest(err, "generate status history id", requestID)
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_request_status_history (
			id, request_id, old_status, new_status, actor_id, actor_name, changed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		histID, requestID, oldStatus, newStatus, actorID, actorDisplayName, changedAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert status history", requestID)
	}
	return nil
}

func ServiceRequestDescriptionUpdated(
	tx *gorm.DB,
	requestID uuid.UUID,
	description string,
	updatedAt time.Time,
) error {
	if err := tx.Exec(`
		UPDATE projections.service_requests
		   SET description = ?, updated_at = ?
		 WHERE id = ?`,
		description, updatedAt, requestID,
	).Error; err != nil {
		return wrapServiceRequest(err, "update description projection", requestID)
	}
	return nil
}

func ServiceRequestExecutorAssigned(
	tx *gorm.DB,
	requestID, employeeID uuid.UUID,
	employeeName, actorID, actorName string,
	changedAt time.Time,
) error {
	histID, err := uuid.NewV7()
	if err != nil {
		return wrapServiceRequest(err, "generate executor history id", requestID)
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_request_executor_history (
			id, request_id, action, employee_id, employee_name,
			actor_id, actor_name, changed_at
		) VALUES (?, ?, 'assigned', ?, ?, ?, ?, ?)`,
		histID, requestID, employeeID, employeeName, actorID, actorName, changedAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert executor assigned history", requestID)
	}
	return nil
}

func ServiceRequestExecutorRemoved(
	tx *gorm.DB,
	requestID, employeeID uuid.UUID,
	employeeName, actorID, actorName string,
	changedAt time.Time,
) error {
	histID, err := uuid.NewV7()
	if err != nil {
		return wrapServiceRequest(err, "generate executor removal history id", requestID)
	}
	if err := tx.Exec(`
		INSERT INTO projections.service_request_executor_history (
			id, request_id, action, employee_id, employee_name,
			actor_id, actor_name, changed_at
		) VALUES (?, ?, 'removed', ?, ?, ?, ?, ?)`,
		histID, requestID, employeeID, employeeName, actorID, actorName, changedAt,
	).Error; err != nil {
		return wrapServiceRequest(err, "insert executor removed history", requestID)
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
