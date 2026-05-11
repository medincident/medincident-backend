package classifier

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	classifierv1 "github.com/medincident/medincident-backend/pkg/event/incident/classifier/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
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
//
// See: docs/services/incident/Classifier.md
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
		env, err := buildIncidentTypeAllowedForPatientsEnvelope(row.ID, updatedAt)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.incident_type.v1.allowed_for_patients", env)
	})
	return AllowIncidentTypeForPatientsResult{}, err
}

func buildIncidentTypeAllowedForPatientsEnvelope(typeID uuid.UUID, updatedAt time.Time) (*eventv1.Envelope, error) {
	msg := &classifierv1.IncidentTypeAllowedForPatients{TypeId: typeID.String(), UpdatedAt: timestamppb.New(updatedAt)}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.incident.classifier.type").Code(ErrCodeIncidentTypeSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{OccurredAt: timestamppb.New(updatedAt), AggregateType: "incident_type", AggregateId: typeID.String(), Payload: payload}, nil
}
