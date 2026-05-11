package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	classifierv1 "github.com/medincident/medincident-backend/pkg/event/incident/classifier/v1"
)

// ── Category projectors ───────────────────────────────────────────────────────

// CategoryCreated inserts a row into projections.incident_categories.
//
// See: docs/services/incident/Classifier.md
func CategoryCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentCategoryCreated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	orgID, err := parseUUID(ev.GetOrganizationId(), "organization_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	createdAt := ev.GetCreatedAt().AsTime()

	var parentID *uuid.UUID
	if sv := ev.GetParentCategoryId(); sv != nil {
		p, parseErr := parseUUID(sv.GetValue(), "parent_category_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
		if parseErr != nil {
			return parseErr
		}
		parentID = &p
	}
	var desc *string
	if sv := ev.GetDescription(); sv != nil {
		v := sv.GetValue()
		desc = &v
	}

	if err := tx.Exec(`
		INSERT INTO projections.incident_categories
		    (id, organization_id, parent_category_id, name, description,
		     is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT DO NOTHING`,
		id, orgID, parentID, ev.GetName(), desc,
		ev.GetIsActive(), createdAt, createdAt,
	).Error; err != nil {
		return wrapClassifier(err, "category", "insert", id)
	}
	return nil
}

// CategoryDetailsUpdated updates name and description in projections.incident_categories.
//
// See: docs/services/incident/Classifier.md
func CategoryDetailsUpdated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentCategoryDetailsUpdated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	var desc *string
	if sv := ev.GetDescription(); sv != nil {
		v := sv.GetValue()
		desc = &v
	}

	if err := tx.Exec(`
		UPDATE projections.incident_categories
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetName(), desc, updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "category", "update_details", id)
	}
	return nil
}

// CategoryMoved updates parent_category_id in projections.incident_categories.
//
// See: docs/services/incident/Classifier.md
func CategoryMoved(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentCategoryMoved,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	var parentID *uuid.UUID
	if sv := ev.GetNewParentCategoryId(); sv != nil {
		p, parseErr := parseUUID(sv.GetValue(), "new_parent_category_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
		if parseErr != nil {
			return parseErr
		}
		parentID = &p
	}

	if err := tx.Exec(`
		UPDATE projections.incident_categories
		   SET parent_category_id = ?, updated_at = ?
		 WHERE id = ?`,
		parentID, updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "category", "move", id)
	}
	return nil
}

// CategoryDeactivated sets is_active=false in projections.incident_categories.
//
// See: docs/services/incident/Classifier.md
func CategoryDeactivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentCategoryDeactivated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.incident_categories
		   SET is_active = FALSE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "category", "deactivate", id)
	}
	return nil
}

// CategoryReactivated sets is_active=true in projections.incident_categories.
//
// See: docs/services/incident/Classifier.md
func CategoryReactivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentCategoryReactivated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.incident_categories
		   SET is_active = TRUE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "category", "reactivate", id)
	}
	return nil
}

// CategoryDeleted removes a row from projections.incident_categories.
//
// See: docs/services/incident/Classifier.md
func CategoryDeleted(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	_ *classifierv1.IncidentCategoryDeleted,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}

	if err := tx.Exec(
		`DELETE FROM projections.incident_categories WHERE id = ?`, id,
	).Error; err != nil {
		return wrapClassifier(err, "category", "delete", id)
	}
	return nil
}

// ── Type projectors ───────────────────────────────────────────────────────────

// TypeCreated inserts a row into projections.incident_types.
//
// See: docs/services/incident/Classifier.md
func TypeCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentTypeCreated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	orgID, err := parseUUID(ev.GetOrganizationId(), "organization_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	catID, err := parseUUID(ev.GetCategoryId(), "category_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	createdAt := ev.GetCreatedAt().AsTime()

	var desc *string
	if sv := ev.GetDescription(); sv != nil {
		v := sv.GetValue()
		desc = &v
	}

	if err := tx.Exec(`
		INSERT INTO projections.incident_types
		    (id, organization_id, category_id, name, description,
		     is_active, is_allowed_for_patients, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT DO NOTHING`,
		id, orgID, catID, ev.GetName(), desc,
		ev.GetIsActive(), ev.GetIsAllowedForPatients(), createdAt, createdAt,
	).Error; err != nil {
		return wrapClassifier(err, "type", "insert", id)
	}
	return nil
}

// TypeDetailsUpdated updates name and description in projections.incident_types.
//
// See: docs/services/incident/Classifier.md
func TypeDetailsUpdated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentTypeDetailsUpdated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	var desc *string
	if sv := ev.GetDescription(); sv != nil {
		v := sv.GetValue()
		desc = &v
	}

	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetName(), desc, updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "type", "update_details", id)
	}
	return nil
}

// TypeMoved updates category_id in projections.incident_types.
//
// See: docs/services/incident/Classifier.md
func TypeMoved(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentTypeMoved,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	catID, err := parseUUID(ev.GetNewCategoryId(), "new_category_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET category_id = ?, updated_at = ?
		 WHERE id = ?`,
		catID, updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "type", "move", id)
	}
	return nil
}

// TypeDeactivated sets is_active=false in projections.incident_types.
//
// See: docs/services/incident/Classifier.md
func TypeDeactivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentTypeDeactivated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET is_active = FALSE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "type", "deactivate", id)
	}
	return nil
}

// TypeReactivated sets is_active=true in projections.incident_types.
//
// See: docs/services/incident/Classifier.md
func TypeReactivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentTypeReactivated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET is_active = TRUE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "type", "reactivate", id)
	}
	return nil
}

// TypeAllowedForPatients sets is_allowed_for_patients=true in projections.incident_types.
//
// See: docs/services/incident/Classifier.md
func TypeAllowedForPatients(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentTypeAllowedForPatients,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET is_allowed_for_patients = TRUE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "type", "allow_for_patients", id)
	}
	return nil
}

// TypeDisallowedForPatients sets is_allowed_for_patients=false in projections.incident_types.
//
// See: docs/services/incident/Classifier.md
func TypeDisallowedForPatients(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *classifierv1.IncidentTypeDisallowedForPatients,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.incident_types
		   SET is_allowed_for_patients = FALSE, updated_at = ?
		 WHERE id = ?`,
		updatedAt, id,
	).Error; err != nil {
		return wrapClassifier(err, "type", "disallow_for_patients", id)
	}
	return nil
}

// TypeDeleted removes a row from projections.incident_types.
//
// See: docs/services/incident/Classifier.md
func TypeDeleted(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	_ *classifierv1.IncidentTypeDeleted,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.incident_classifier", ErrCodeIncidentClassifierProjectionFailed)
	if err != nil {
		return err
	}

	if err := tx.Exec(
		`DELETE FROM projections.incident_types WHERE id = ?`, id,
	).Error; err != nil {
		return wrapClassifier(err, "type", "delete", id)
	}
	return nil
}

func wrapClassifier(err error, aggregate, action string, id uuid.UUID) error {
	return oops.In("projector.incident.classifier."+aggregate).
		Code(ErrCodeIncidentClassifierProjectionFailed).
		With("action", action).
		With(aggregate+"_id", id).
		Wrap(err)
}

// ── Projectors forwarding methods ────────────────────────────────────────────

func (p *Projectors) CategoryCreated(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentCategoryCreated) error {
	return CategoryCreated(tx, id, t, ev)
}

func (p *Projectors) CategoryDetailsUpdated(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentCategoryDetailsUpdated) error {
	return CategoryDetailsUpdated(tx, id, t, ev)
}

func (p *Projectors) CategoryMoved(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentCategoryMoved) error {
	return CategoryMoved(tx, id, t, ev)
}

func (p *Projectors) CategoryDeactivated(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentCategoryDeactivated) error {
	return CategoryDeactivated(tx, id, t, ev)
}

func (p *Projectors) CategoryReactivated(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentCategoryReactivated) error {
	return CategoryReactivated(tx, id, t, ev)
}

func (p *Projectors) CategoryDeleted(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentCategoryDeleted) error {
	return CategoryDeleted(tx, id, t, ev)
}

func (p *Projectors) TypeCreated(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentTypeCreated) error {
	return TypeCreated(tx, id, t, ev)
}

func (p *Projectors) TypeDetailsUpdated(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentTypeDetailsUpdated) error {
	return TypeDetailsUpdated(tx, id, t, ev)
}

func (p *Projectors) TypeMoved(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentTypeMoved) error {
	return TypeMoved(tx, id, t, ev)
}

func (p *Projectors) TypeDeactivated(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentTypeDeactivated) error {
	return TypeDeactivated(tx, id, t, ev)
}

func (p *Projectors) TypeReactivated(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentTypeReactivated) error {
	return TypeReactivated(tx, id, t, ev)
}

func (p *Projectors) TypeAllowedForPatients(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentTypeAllowedForPatients) error {
	return TypeAllowedForPatients(tx, id, t, ev)
}

func (p *Projectors) TypeDisallowedForPatients(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentTypeDisallowedForPatients) error {
	return TypeDisallowedForPatients(tx, id, t, ev)
}

func (p *Projectors) TypeDeleted(tx *gorm.DB, id string, t time.Time, ev *classifierv1.IncidentTypeDeleted) error {
	return TypeDeleted(tx, id, t, ev)
}
