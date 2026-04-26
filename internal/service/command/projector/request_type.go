package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
)

func RequestTypeCreated(tx *gorm.DB, rt *model.RequestType) error {
	if err := tx.Exec(`
		INSERT INTO projections.request_types
		    (id, organization_id, name, description, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		rt.ID, rt.OrganizationID, rt.Name, rt.Description,
		rt.IsActive, rt.CreatedAt, rt.UpdatedAt,
	).Error; err != nil {
		return wrapRequestType(err, "create", rt.ID)
	}
	return nil
}

func RequestTypeDetailsUpdated(tx *gorm.DB, rt *model.RequestType) error {
	if err := tx.Exec(`
		UPDATE projections.request_types
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		rt.Name, rt.Description, rt.UpdatedAt, rt.ID,
	).Error; err != nil {
		return wrapRequestType(err, "update details", rt.ID)
	}
	return nil
}

func RequestTypeDeactivated(tx *gorm.DB, typeID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.request_types
		   SET is_active = FALSE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, typeID,
	).Error; err != nil {
		return wrapRequestType(err, "deactivate", typeID)
	}
	return nil
}

func RequestTypeReactivated(tx *gorm.DB, typeID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.request_types
		   SET is_active = TRUE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, typeID,
	).Error; err != nil {
		return wrapRequestType(err, "reactivate", typeID)
	}
	return nil
}

func RequestTypeDeleted(tx *gorm.DB, typeID uuid.UUID) error {
	if err := tx.Exec(
		`DELETE FROM projections.request_types WHERE id = ?`, typeID,
	).Error; err != nil {
		return wrapRequestType(err, "delete", typeID)
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
