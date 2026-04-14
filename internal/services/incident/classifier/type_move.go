package classifier

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	typeeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/type/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

const (
	ErrCodeIncidentTypeMoveOrganizationMismatch = "incident_type_move_organization_mismatch"
)

type MoveIncidentTypeCommand struct {
	TypeID        uuid.UUID
	NewCategoryID uuid.UUID
}

type MoveIncidentTypeResult struct{}

func (s *IncidentTypeService) Move(
	ctx context.Context,
	cmd MoveIncidentTypeCommand,
) (MoveIncidentTypeResult, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var moving model.IncidentType
		if err := tx.First(&moving, "id = ?", cmd.TypeID).Error; err != nil {
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

		var newCategory model.IncidentCategory
		if err := tx.First(&newCategory, "id = ?", cmd.NewCategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeCategoryNotFound).
					Public("New incident category not found.").
					With("incident_category_id", cmd.NewCategoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_category_id", cmd.NewCategoryID).
				Wrap(err)
		}

		if newCategory.OrganizationID != moving.OrganizationID {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeMoveOrganizationMismatch).
				Public("New incident category belongs to a different organization.").
				With("incident_type_id", moving.ID).
				With("incident_category_id", newCategory.ID).
				Errorf("organization mismatch on type move")
		}

		if err := tx.Model(&model.IncidentType{}).
			Where("id = ?", moving.ID).
			Update("category_id", newCategory.ID).Error; err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", moving.ID).
				Wrap(err)
		}

		payload, err := anypb.New(&typeeventv1.IncidentTypeMoved{
			NewCategoryId: newCategory.ID.String(),
		})
		if err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeEventBuildFailed).
				With("incident_type_id", moving.ID).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.Now(),
			AggregateType: AggregateTypeIncidentType,
			AggregateId:   moving.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectIncidentTypeMoved, envelope, nil)
	})
	return MoveIncidentTypeResult{}, err
}
