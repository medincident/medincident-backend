package classifier

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// DeactivateIncidentTypePayload identifies the incident type to
// deactivate.
type DeactivateIncidentTypePayload struct {
	TypeID string `validate:"required,uuid"`
}

// DeactivateIncidentTypeCommand = caller + payload.
type DeactivateIncidentTypeCommand struct {
	Caller  authz.Caller
	Payload DeactivateIncidentTypePayload
}

// DeactivateIncidentTypeResult is empty.
type DeactivateIncidentTypeResult struct{}

// Deactivate marks an incident type as inactive.
//
// See: docs/services/incident/Classifier.md
func (s *IncidentTypeService) Deactivate(
	ctx context.Context,
	cmd DeactivateIncidentTypeCommand,
) (DeactivateIncidentTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DeactivateIncidentTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.IncidentType(typeID)); err != nil {
		return DeactivateIncidentTypeResult{}, err
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

		if !row.IsActive {
			return nil
		}
		// Update the row and return updated_at in a single round-trip
		// so the outbox event's OccurredAt matches the DB clock.
		var updatedAt time.Time
		if err := tx.Raw(`
			UPDATE domain.incident_types
			SET is_active = FALSE, updated_at = now()
			WHERE id = ?
			RETURNING updated_at`, row.ID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		return appendTypeDeactivatedEvent(tx, row.ID, updatedAt)
	})
	return DeactivateIncidentTypeResult{}, err
}
