package classifier

import (
	"context"
	"errors"

	"github.com/google/uuid"
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

type DeleteIncidentTypeCommand struct {
	TypeID uuid.UUID
}

type DeleteIncidentTypeResult struct{}

func (s *IncidentTypeService) Delete(
	ctx context.Context,
	cmd DeleteIncidentTypeCommand,
) (DeleteIncidentTypeResult, error) {
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

		payload, err := anypb.New(&typeeventv1.IncidentTypeDeleted{})
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
		if err := outbox.AppendEvent(tx, SubjectIncidentTypeDeleted, envelope, nil); err != nil {
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
