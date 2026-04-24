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

// ErrCodeIncidentTypeInactive is returned when a mutation method is
// invoked on a type whose is_active is FALSE. Shared with
// type_disallow_for_patients.go.
const ErrCodeIncidentTypeInactive = "incident_type_inactive"

// AllowIncidentTypeForPatientsPayload identifies the incident type to
// mark as allowed for patient submission.
type AllowIncidentTypeForPatientsPayload struct {
	TypeID string `validate:"required,uuid"`
}

// AllowIncidentTypeForPatientsCommand = caller + payload.
type AllowIncidentTypeForPatientsCommand struct {
	Caller  authz.Caller
	Payload AllowIncidentTypeForPatientsPayload
}

// AllowIncidentTypeForPatientsResult is empty.
type AllowIncidentTypeForPatientsResult struct{}

// AllowForPatients marks a type as allowed for patients to use when
// submitting incidents. Idempotent: a no-op if the flag is already TRUE.
// Requires the type to be is_active = TRUE.
func (s *IncidentTypeService) AllowForPatients(
	ctx context.Context,
	cmd AllowIncidentTypeForPatientsCommand,
) (AllowIncidentTypeForPatientsResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return AllowIncidentTypeForPatientsResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.IncidentType(typeID)); err != nil {
		return AllowIncidentTypeForPatientsResult{}, err
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

		// Idempotent: nothing to do if already allowed.
		if row.IsAllowedForPatients {
			return nil
		}

		var updatedAt time.Time
		if err := tx.Raw(`
			UPDATE domain.incident_types
			SET is_allowed_for_patients = TRUE, updated_at = now()
			WHERE id = ?
			RETURNING updated_at`, row.ID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		return projector.TypeAllowForPatients(tx, row.ID, updatedAt)
	})
	return AllowIncidentTypeForPatientsResult{}, err
}
