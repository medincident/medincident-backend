package classifier

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

const (
	ErrCodeIncidentTypeReactivateInactiveAncestor = "incident_type_reactivate_inactive_ancestor"
	ErrCodeIncidentTypeReactivateNameConflict     = "incident_type_reactivate_name_conflict"
)

// ReactivateIncidentTypePayload identifies the incident type to
// reactivate.
type ReactivateIncidentTypePayload struct {
	TypeID string `validate:"required,uuid"`
}

// ReactivateIncidentTypeCommand = caller + payload.
type ReactivateIncidentTypeCommand struct {
	Caller  authz.Caller
	Payload ReactivateIncidentTypePayload
}

// ReactivateIncidentTypeResult is empty.
type ReactivateIncidentTypeResult struct{}

// typeHasInactiveAncestor walks from the type's owning category up to
// the root and returns the first inactive category id found (either
// the category itself or any ancestor). Returns uuid.Nil when the
// whole chain is active.
func typeHasInactiveAncestor(tx *gorm.DB, categoryID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	const q = `
		WITH RECURSIVE chain AS (
			SELECT id, parent_category_id, is_active, 0 AS depth
			FROM domain.incident_categories
			WHERE id = ?
			UNION ALL
			SELECT c.id, c.parent_category_id, c.is_active, ch.depth + 1
			FROM domain.incident_categories c
			JOIN chain ch ON ch.parent_category_id = c.id
		)
		SELECT id FROM chain
		WHERE is_active = FALSE
		ORDER BY depth ASC
		LIMIT 1;
	`
	row := tx.Raw(q, categoryID).Row()
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}
	return id, nil
}

// Reactivate marks a deactivated incident type as active. Fails if the owning category is inactive.
//
// See: docs/services/incident/Classifier.md
func (s *IncidentTypeService) Reactivate(
	ctx context.Context,
	cmd ReactivateIncidentTypeCommand,
) (ReactivateIncidentTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return ReactivateIncidentTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.IncidentType(typeID)); err != nil {
		return ReactivateIncidentTypeResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.IncidentType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&row, "id = ?", typeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNotFound).
					Public("Incident type not found.").
					With("incident_type_id", typeID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_type_id", typeID).
				Wrap(err)
		}

		if err := lockClassifierOrg(tx, row.OrganizationID); err != nil {
			return err
		}

		if row.IsActive {
			return nil
		}

		blocker, err := typeHasInactiveAncestor(tx, row.CategoryID)
		if err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		if blocker != uuid.Nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeReactivateInactiveAncestor).
				Public("Cannot reactivate: the owning category or one of its ancestors is inactive.").
				With("incident_type_id", row.ID).
				With("inactive_category_id", blocker).
				Errorf("inactive ancestor blocks type reactivation")
		}

		var updatedAt time.Time
		if err := tx.Raw(`
			UPDATE domain.incident_types
			SET is_active = TRUE, updated_at = now()
			WHERE id = ?
			RETURNING updated_at`, row.ID,
		).Row().Scan(&updatedAt); err != nil {
			// Raw SQL bypasses GORM's TranslateError — check pgconn directly.
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeReactivateNameConflict).
					Public("Cannot reactivate: another active incident type uses this name.").
					With("incident_type_id", row.ID).
					With("name", row.Name).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}

		return projector.TypeReactivate(tx, row.ID, updatedAt)
	})
	return ReactivateIncidentTypeResult{}, err
}
