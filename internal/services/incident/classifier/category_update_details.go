package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	categoryeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/category/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

type UpdateIncidentCategoryDetailsCommand struct {
	CategoryID  uuid.UUID
	Name        string
	Description *string
}

type UpdateIncidentCategoryDetailsResult struct{}

func buildIncidentCategoryDetailsChangedEvent(c *model.IncidentCategory) *categoryeventv1.IncidentCategoryDetailsChanged {
	ev := &categoryeventv1.IncidentCategoryDetailsChanged{Name: c.Name}
	if c.Description.Valid {
		d := c.Description.String
		ev.Description = &d
	}
	return ev
}

func (s *IncidentCategoryService) UpdateDetails(
	ctx context.Context,
	cmd UpdateIncidentCategoryDetailsCommand,
) (UpdateIncidentCategoryDetailsResult, error) {
	var errs []error
	if err := validateIncidentCategoryName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateIncidentCategoryDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return UpdateIncidentCategoryDetailsResult{}, errors.Join(errs...)
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cat model.IncidentCategory
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&cat, "id = ?", cmd.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryNotFound).
					Public("Incident category not found.").
					With("incident_category_id", cmd.CategoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryLoadFailed).
				With("incident_category_id", cmd.CategoryID).
				Wrap(err)
		}

		cat.Name = strings.TrimSpace(cmd.Name)
		cat.Description = null.StringFromPtr(trimmedStringPtr(cmd.Description))

		if err := tx.Save(&cat).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
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

		event := buildIncidentCategoryDetailsChangedEvent(&cat)
		payload, err := anypb.New(event)
		if err != nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryEventBuildFailed).
				With("incident_category_id", cat.ID).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(cat.UpdatedAt),
			AggregateType: AggregateTypeIncidentCategory,
			AggregateId:   cat.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectIncidentCategoryDetailsChanged, envelope, nil)
	})
	return UpdateIncidentCategoryDetailsResult{}, err
}
