package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
)

// Incident classifier projectors mirror domain.incident_categories and
// domain.incident_types into projections.*. Categories/types track
// tree moves, active/inactive toggles, and cascading deletes.

// CategoryCreated writes projections.incident_categories for a new row.
func CategoryCreated(tx *gorm.DB, c *model.IncidentCategory) error {
	var parentID *uuid.UUID
	if c.ParentCategoryID.Valid {
		p := c.ParentCategoryID.UUID
		parentID = &p
	}
	if err := tx.Exec(`
		INSERT INTO projections.incident_categories
		    (id, organization_id, parent_category_id, name, description,
		     is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.OrganizationID, parentID, c.Name, c.Description,
		c.IsActive, c.CreatedAt, c.UpdatedAt,
	).Error; err != nil {
		return wrapIncident(err, "category", c.ID)
	}
	return nil
}

// CategoryUpdateDetails updates name + description + updated_at.
func CategoryUpdateDetails(tx *gorm.DB, c *model.IncidentCategory) error {
	if err := tx.Exec(`
		UPDATE projections.incident_categories
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		c.Name, c.Description, c.UpdatedAt, c.ID,
	).Error; err != nil {
		return wrapIncident(err, "category", c.ID)
	}
	return nil
}

// CategoryMove updates parent_category_id + updated_at.
func CategoryMove(tx *gorm.DB, categoryID uuid.UUID, newParentCategoryID uuid.NullUUID, updatedAt time.Time) error {
	var parentID *uuid.UUID
	if newParentCategoryID.Valid {
		p := newParentCategoryID.UUID
		parentID = &p
	}
	if err := tx.Exec(`
		UPDATE projections.incident_categories
		   SET parent_category_id = ?, updated_at = ?
		 WHERE id = ?`,
		parentID, updatedAt, categoryID,
	).Error; err != nil {
		return wrapIncident(err, "category", categoryID)
	}
	return nil
}

// CategoryDeactivate flips is_active FALSE on a single category row.
// Cascading deactivation is driven by the service layer producing one
// call per affected row.
func CategoryDeactivate(tx *gorm.DB, categoryID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.incident_categories
		   SET is_active = FALSE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, categoryID,
	).Error; err != nil {
		return wrapIncident(err, "category", categoryID)
	}
	return nil
}

// CategoryReactivate flips is_active TRUE on a single category row.
func CategoryReactivate(tx *gorm.DB, categoryID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.incident_categories
		   SET is_active = TRUE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, categoryID,
	).Error; err != nil {
		return wrapIncident(err, "category", categoryID)
	}
	return nil
}

// CategoryDeleted removes the row from the projection.
func CategoryDeleted(tx *gorm.DB, categoryID uuid.UUID) error {
	if err := tx.Exec(
		`DELETE FROM projections.incident_categories WHERE id = ?`, categoryID,
	).Error; err != nil {
		return wrapIncident(err, "category", categoryID)
	}
	return nil
}

// TypeCreated writes projections.incident_types for a new row.
func TypeCreated(tx *gorm.DB, t *model.IncidentType) error {
	if err := tx.Exec(`
		INSERT INTO projections.incident_types
		    (id, organization_id, category_id, name, description,
		     is_active, is_allowed_for_patients, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.OrganizationID, t.CategoryID, t.Name, t.Description,
		t.IsActive, t.IsAllowedForPatients, t.CreatedAt, t.UpdatedAt,
	).Error; err != nil {
		return wrapIncident(err, "type", t.ID)
	}
	return nil
}

// TypeUpdateDetails updates name + description + updated_at.
func TypeUpdateDetails(tx *gorm.DB, t *model.IncidentType) error {
	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		t.Name, t.Description, t.UpdatedAt, t.ID,
	).Error; err != nil {
		return wrapIncident(err, "type", t.ID)
	}
	return nil
}

// TypeMove updates category_id + updated_at.
func TypeMove(tx *gorm.DB, typeID, newCategoryID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET category_id = ?, updated_at = ?
		 WHERE id = ?`,
		newCategoryID, updatedAt, typeID,
	).Error; err != nil {
		return wrapIncident(err, "type", typeID)
	}
	return nil
}

// TypeDeactivate flips is_active FALSE on a type row.
func TypeDeactivate(tx *gorm.DB, typeID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET is_active = FALSE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, typeID,
	).Error; err != nil {
		return wrapIncident(err, "type", typeID)
	}
	return nil
}

// TypeReactivate flips is_active TRUE on a type row.
func TypeReactivate(tx *gorm.DB, typeID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET is_active = TRUE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, typeID,
	).Error; err != nil {
		return wrapIncident(err, "type", typeID)
	}
	return nil
}

// TypeAllowForPatients flips is_allowed_for_patients TRUE on a type row.
func TypeAllowForPatients(tx *gorm.DB, typeID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET is_allowed_for_patients = TRUE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, typeID,
	).Error; err != nil {
		return wrapIncident(err, "type", typeID)
	}
	return nil
}

// TypeDisallowForPatients flips is_allowed_for_patients FALSE on a type row.
func TypeDisallowForPatients(tx *gorm.DB, typeID uuid.UUID, updatedAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET is_allowed_for_patients = FALSE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, typeID,
	).Error; err != nil {
		return wrapIncident(err, "type", typeID)
	}
	return nil
}

// TypeDeleted removes the row from the projection.
func TypeDeleted(tx *gorm.DB, typeID uuid.UUID) error {
	if err := tx.Exec(
		`DELETE FROM projections.incident_types WHERE id = ?`, typeID,
	).Error; err != nil {
		return wrapIncident(err, "type", typeID)
	}
	return nil
}

// wrapIncident wraps a gorm error with the per-aggregate oops scope.
func wrapIncident(err error, aggregate string, id uuid.UUID) error {
	return oops.In("projector.incident."+aggregate).
		Code(ErrCodeIncidentProjectionFailed).
		With(aggregate+"_id", id).
		Wrap(err)
}
