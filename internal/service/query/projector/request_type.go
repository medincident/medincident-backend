package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	requesttypev1 "github.com/medincident/medincident-backend/pkg/event/request_type/v1"
)

// RequestTypeCreated inserts a row into projections.request_types.
//
// See: docs/services/request/Classifier.md
func RequestTypeCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *requesttypev1.RequestTypeCreated,
) error {
	id := uuid.MustParse(aggregateID)
	orgID := uuid.MustParse(ev.GetOrganizationId())
	createdAt := ev.GetCreatedAt().AsTime()

	var desc *string
	if sv := ev.GetDescription(); sv != nil {
		v := sv.GetValue()
		desc = &v
	}

	if err := tx.Exec(`
		INSERT INTO projections.request_types
		    (id, organization_id, name, description, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT DO NOTHING`,
		id, orgID, ev.GetName(), desc, ev.GetIsActive(), createdAt, createdAt,
	).Error; err != nil {
		return wrapRequestType(err, "insert", id)
	}
	return nil
}

// RequestTypeDetailsUpdated updates name, description, and updated_at.
//
// See: docs/services/request/Classifier.md
func RequestTypeDetailsUpdated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *requesttypev1.RequestTypeDetailsUpdated,
) error {
	id := uuid.MustParse(aggregateID)
	updatedAt := ev.GetUpdatedAt().AsTime()

	var desc *string
	if sv := ev.GetDescription(); sv != nil {
		v := sv.GetValue()
		desc = &v
	}

	if err := tx.Exec(`
		UPDATE projections.request_types
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetName(), desc, updatedAt, id,
	).Error; err != nil {
		return wrapRequestType(err, "update_details", id)
	}
	return nil
}

// RequestTypeDeactivated sets is_active=false.
//
// See: docs/services/request/Classifier.md
func RequestTypeDeactivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *requesttypev1.RequestTypeDeactivated,
) error {
	id := uuid.MustParse(aggregateID)
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.request_types
		   SET is_active = FALSE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, id,
	).Error; err != nil {
		return wrapRequestType(err, "deactivate", id)
	}
	return nil
}

// RequestTypeReactivated sets is_active=true.
//
// See: docs/services/request/Classifier.md
func RequestTypeReactivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *requesttypev1.RequestTypeReactivated,
) error {
	id := uuid.MustParse(aggregateID)
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.request_types
		   SET is_active = TRUE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, id,
	).Error; err != nil {
		return wrapRequestType(err, "reactivate", id)
	}
	return nil
}

// RequestTypeDeleted removes a row from projections.request_types.
//
// See: docs/services/request/Classifier.md
func RequestTypeDeleted(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	_ *requesttypev1.RequestTypeDeleted,
) error {
	id := uuid.MustParse(aggregateID)

	if err := tx.Exec(
		`DELETE FROM projections.request_types WHERE id = ?`, id,
	).Error; err != nil {
		return wrapRequestType(err, "delete", id)
	}
	return nil
}

func wrapRequestType(err error, action string, id uuid.UUID) error {
	return oops.In("projector.request_type").
		Code(ErrCodeRequestTypeProjectionFailed).
		With("action", action).
		With("request_type_id", id).
		Wrap(err)
}

// ── Projectors forwarding methods ────────────────────────────────────────────

// RequestTypeCreated projects a RequestTypeCreated event.
//
// See: docs/services/request/Classifier.md
func (p *Projectors) RequestTypeCreated(tx *gorm.DB, id string, t time.Time, ev *requesttypev1.RequestTypeCreated) error {
	return RequestTypeCreated(tx, id, t, ev)
}

// RequestTypeDetailsUpdated projects a RequestTypeDetailsUpdated event.
//
// See: docs/services/request/Classifier.md
func (p *Projectors) RequestTypeDetailsUpdated(tx *gorm.DB, id string, t time.Time, ev *requesttypev1.RequestTypeDetailsUpdated) error {
	return RequestTypeDetailsUpdated(tx, id, t, ev)
}

// RequestTypeDeactivated projects a RequestTypeDeactivated event.
//
// See: docs/services/request/Classifier.md
func (p *Projectors) RequestTypeDeactivated(tx *gorm.DB, id string, t time.Time, ev *requesttypev1.RequestTypeDeactivated) error {
	return RequestTypeDeactivated(tx, id, t, ev)
}

// RequestTypeReactivated projects a RequestTypeReactivated event.
//
// See: docs/services/request/Classifier.md
func (p *Projectors) RequestTypeReactivated(tx *gorm.DB, id string, t time.Time, ev *requesttypev1.RequestTypeReactivated) error {
	return RequestTypeReactivated(tx, id, t, ev)
}

// RequestTypeDeleted projects a RequestTypeDeleted event.
//
// See: docs/services/request/Classifier.md
func (p *Projectors) RequestTypeDeleted(tx *gorm.DB, id string, t time.Time, ev *requesttypev1.RequestTypeDeleted) error {
	return RequestTypeDeleted(tx, id, t, ev)
}
