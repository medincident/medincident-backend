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

// DeleteIncidentTypePayload identifies the incident type to delete.
type DeleteIncidentTypePayload struct {
	TypeID string `validate:"required,uuid"`
}

// DeleteIncidentTypeCommand = caller + payload.
type DeleteIncidentTypeCommand struct {
	Caller  authz.Caller
	Payload DeleteIncidentTypePayload
}

// DeleteIncidentTypeResult is empty.
type DeleteIncidentTypeResult struct{}

// Delete permanently removes an incident type.
//
// See: docs/services/incident/Classifier.md
func (s *IncidentTypeService) Delete(
	ctx context.Context,
	cmd DeleteIncidentTypeCommand,
) (DeleteIncidentTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DeleteIncidentTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.IncidentType(typeID)); err != nil {
		return DeleteIncidentTypeResult{}, err
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

		now := time.Now().UTC()
		env, err := buildIncidentTypeDeletedEnvelope(row.ID, now)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.incident_type.v1.deleted", env); err != nil {
			return err
		}
		if err := tx.Delete(&model.IncidentType{}, "id = ?", row.ID).Error; err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		return nil
	})
	return DeleteIncidentTypeResult{}, err
}

func buildIncidentTypeDeletedEnvelope(typeID uuid.UUID, deletedAt time.Time) (*eventv1.Envelope, error) {
	msg := &classifierv1.IncidentTypeDeleted{TypeId: typeID.String(), DeletedAt: timestamppb.New(deletedAt)}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.incident.classifier.type").Code(ErrCodeIncidentTypeSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{OccurredAt: timestamppb.New(deletedAt), AggregateType: "incident_type", AggregateId: typeID.String(), Payload: payload}, nil
}
