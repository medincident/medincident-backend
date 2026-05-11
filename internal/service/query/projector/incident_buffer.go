package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	bufferv1 "github.com/medincident/medincident-backend/pkg/event/incident/buffer/v1"
)

// PatientIncidentBufferCreated inserts a buffer projection row.
//
// See: docs/services/incident/Buffer.md
func PatientIncidentBufferCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *bufferv1.PatientIncidentBufferCreated,
) error {
	id := uuid.MustParse(aggregateID)
	orgID := uuid.MustParse(ev.GetOrganizationId())
	createdAt := ev.GetCreatedAt().AsTime()
	occAt := ev.GetOccurredAt().AsTime()

	var catID *uuid.UUID
	if sv := ev.GetCategoryId(); sv != nil {
		p := uuid.MustParse(sv.GetValue())
		catID = &p
	}
	var typeID *uuid.UUID
	if sv := ev.GetTypeId(); sv != nil {
		p := uuid.MustParse(sv.GetValue())
		typeID = &p
	}

	if err := tx.Exec(`
		INSERT INTO projections.patient_incident_buffer (
			id, organization_id, patient_zitadel_user_id,
			category_id, type_id, description, occurred_at,
			status, published_incident_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, ?)
		ON CONFLICT DO NOTHING`,
		id, orgID, ev.GetPatientZitadelUserId(),
		catID, typeID, ev.GetDescription(), occAt,
		ev.GetStatus(), createdAt, createdAt,
	).Error; err != nil {
		return wrapIncidentBuffer(err, "insert buffer", id)
	}
	return nil
}

// PatientIncidentBufferUpdated mirrors a buffer mutation (edit / cancel / reject / publish).
//
// See: docs/services/incident/Buffer.md
func PatientIncidentBufferUpdated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *bufferv1.PatientIncidentBufferUpdated,
) error {
	id := uuid.MustParse(aggregateID)
	updatedAt := ev.GetUpdatedAt().AsTime()
	occAt := ev.GetOccurredAt().AsTime()

	var catID *uuid.UUID
	if sv := ev.GetCategoryId(); sv != nil {
		p := uuid.MustParse(sv.GetValue())
		catID = &p
	}
	var typeID *uuid.UUID
	if sv := ev.GetTypeId(); sv != nil {
		p := uuid.MustParse(sv.GetValue())
		typeID = &p
	}
	var pubIncID null.String
	if sv := ev.GetPublishedIncidentId(); sv != nil {
		pubIncID = null.StringFrom(sv.GetValue())
	}

	if err := tx.Exec(`
		UPDATE projections.patient_incident_buffer
		   SET category_id = ?, type_id = ?, description = ?, occurred_at = ?,
		       status = ?, published_incident_id = ?, updated_at = ?
		 WHERE id = ?`,
		catID, typeID, ev.GetDescription(), occAt,
		ev.GetStatus(), pubIncID, updatedAt, id,
	).Error; err != nil {
		return wrapIncidentBuffer(err, "update buffer", id)
	}
	return nil
}

func wrapIncidentBuffer(err error, action string, id uuid.UUID) error {
	return oops.In("projector.incident.buffer").
		Code(ErrCodeIncidentBufferProjectionFailed).
		With("action", action).
		With("buffer_id", id).
		Wrap(err)
}

// ── Projectors forwarding methods ────────────────────────────────────────────

// PatientIncidentBufferCreated forwards to the package-level projector function.
func (p *Projectors) PatientIncidentBufferCreated(tx *gorm.DB, id string, t time.Time, ev *bufferv1.PatientIncidentBufferCreated) error {
	return PatientIncidentBufferCreated(tx, id, t, ev)
}

// PatientIncidentBufferUpdated forwards to the package-level projector function.
func (p *Projectors) PatientIncidentBufferUpdated(tx *gorm.DB, id string, t time.Time, ev *bufferv1.PatientIncidentBufferUpdated) error {
	return PatientIncidentBufferUpdated(tx, id, t, ev)
}
