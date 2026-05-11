package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	classifierv1 "github.com/medincident/medincident-backend/pkg/event/incident/classifier/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// UpdateIncidentCategoryDetailsPayload carries the new name and
// (optional) description for an existing incident category.
type UpdateIncidentCategoryDetailsPayload struct {
	CategoryID  string  `validate:"required,uuid"`
	Name        string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// UpdateIncidentCategoryDetailsCommand = caller + payload.
type UpdateIncidentCategoryDetailsCommand struct {
	Caller  authz.Caller
	Payload UpdateIncidentCategoryDetailsPayload
}

// UpdateIncidentCategoryDetailsResult is empty — the event is the
// meaningful result.
type UpdateIncidentCategoryDetailsResult struct{}

// UpdateDetails changes the name and description of an incident category.
//
// See: docs/services/incident/Classifier.md
func (s *IncidentCategoryService) UpdateDetails(
	ctx context.Context,
	cmd UpdateIncidentCategoryDetailsCommand,
) (UpdateIncidentCategoryDetailsResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return UpdateIncidentCategoryDetailsResult{}, err
	}
	categoryID := uuid.MustParse(cmd.Payload.CategoryID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Category(categoryID)); err != nil {
		return UpdateIncidentCategoryDetailsResult{}, err
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cat model.IncidentCategory
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&cat, "id = ?", categoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryNotFound).
					Public("Incident category not found.").
					With("incident_category_id", categoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryLoadFailed).
				With("incident_category_id", categoryID).
				Wrap(err)
		}

		if err := lockClassifierOrg(tx, cat.OrganizationID); err != nil {
			return err
		}

		newName := strings.TrimSpace(cmd.Payload.Name)
		var newDescription null.String
		if cmd.Payload.Description != nil {
			newDescription = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if cat.Name == newName && cat.Description == newDescription {
			return nil
		}
		cat.Name = newName
		cat.Description = newDescription

		if err := tx.Save(&cat).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryNameConflict).
					Public("An active incident category with this name already exists.").
					With("incident_category_id", cat.ID).
					With("name", cat.Name).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategorySaveFailed).
				With("incident_category_id", cat.ID).
				Wrap(err)
		}

		env, err := buildIncidentCategoryDetailsUpdatedEnvelope(&cat)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.incident_category.v1.details_updated", env)
	})
	return UpdateIncidentCategoryDetailsResult{}, err
}

func buildIncidentCategoryDetailsUpdatedEnvelope(cat *model.IncidentCategory) (*eventv1.Envelope, error) {
	msg := &classifierv1.IncidentCategoryDetailsUpdated{
		CategoryId: cat.ID.String(),
		Name:       cat.Name,
		UpdatedAt:  timestamppb.New(cat.UpdatedAt),
	}
	if cat.Description.Valid {
		msg.Description = wrapperspb.String(cat.Description.String)
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.incident.classifier.category").Code(ErrCodeIncidentCategorySaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt: timestamppb.New(cat.UpdatedAt), AggregateType: "incident_category",
		AggregateId: cat.ID.String(), Payload: payload,
	}, nil
}
