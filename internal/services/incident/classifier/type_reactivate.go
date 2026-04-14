package classifier

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	typeeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/type/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

const (
	ErrCodeIncidentTypeReactivateInactiveAncestor = "incident_type_reactivate_inactive_ancestor"
	ErrCodeIncidentTypeReactivateNameConflict     = "incident_type_reactivate_name_conflict"
)

type ReactivateIncidentTypeCommand struct {
	TypeID uuid.UUID
}

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

func (s *IncidentTypeService) Reactivate(
	ctx context.Context,
	cmd ReactivateIncidentTypeCommand,
) (ReactivateIncidentTypeResult, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.IncidentType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&row, "id = ?", cmd.TypeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNotFound).
					Public("Incident type not found.").
					With("incident_type_id", cmd.TypeID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_type_id", cmd.TypeID).
				Wrap(err)
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

		if err := tx.Model(&model.IncidentType{}).
			Where("id = ?", row.ID).
			Updates(map[string]any{"is_active": true}).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
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

		payload, err := anypb.New(&typeeventv1.IncidentTypeReactivated{})
		if err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeEventBuildFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.Now(),
			AggregateType: AggregateTypeIncidentType,
			AggregateId:   row.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectIncidentTypeReactivated, envelope)
	})
	return ReactivateIncidentTypeResult{}, err
}
