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
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// DisallowIncidentTypeForPatientsPayload identifies the incident type to
// remove from the patient-submission allow list.
type DisallowIncidentTypeForPatientsPayload struct {
	TypeID string `validate:"required,uuid"`
}

// DisallowIncidentTypeForPatientsCommand = caller + payload.
type DisallowIncidentTypeForPatientsCommand struct {
	Caller  authz.Caller
	Payload DisallowIncidentTypeForPatientsPayload
}

// DisallowIncidentTypeForPatientsResult is empty.
type DisallowIncidentTypeForPatientsResult struct{}

// DisallowForPatients revokes the patient-submission permission for the
// given type. Idempotent: a no-op if the flag is already FALSE. Requires
// the type to be is_active = TRUE.
func (s *IncidentTypeService) DisallowForPatients(
	ctx context.Context,
	cmd DisallowIncidentTypeForPatientsCommand,
) (DisallowIncidentTypeForPatientsResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DisallowIncidentTypeForPatientsResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.IncidentType(typeID)); err != nil {
		return DisallowIncidentTypeForPatientsResult{}, err
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
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeInactive).
				Public("Cannot mutate an inactive incident type.").
				With("incident_type_id", row.ID).
				Errorf("type inactive")
		}

		// Idempotent: nothing to do if already disallowed.
		if !row.IsAllowedForPatients {
			return nil
		}

		var updatedAt time.Time
		if err := tx.Raw(`
			UPDATE domain.incident_types
			SET is_allowed_for_patients = FALSE, updated_at = now()
			WHERE id = ?
			RETURNING updated_at`, row.ID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		return projector.TypeDisallowForPatients(tx, row.ID, updatedAt)
	})
	return DisallowIncidentTypeForPatientsResult{}, err
}
