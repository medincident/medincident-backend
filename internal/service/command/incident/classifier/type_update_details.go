package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/outbox"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	typeeventv1 "github.com/medincident/medincident-command-service/pkg/event/incident/type/v1"
)

type UpdateIncidentTypeDetailsCommand struct {
	TypeID      uuid.UUID
	Name        string
	Description *string
}

type UpdateIncidentTypeDetailsResult struct{}

func buildIncidentTypeDetailsChangedEvent(t *model.IncidentType) *typeeventv1.IncidentTypeDetailsChanged {
	return &typeeventv1.IncidentTypeDetailsChanged{
		Name:        t.Name,
		Description: t.Description.Ptr(),
	}
}

func (s *IncidentTypeService) UpdateDetails(
	ctx context.Context,
	cmd UpdateIncidentTypeDetailsCommand,
) (UpdateIncidentTypeDetailsResult, error) {
	var errs []error
	if cmd.TypeID == uuid.Nil {
		errs = append(errs, oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeIDEmpty).
			Public("Incident type ID is required.").
			With("field", "type_id").
			Errorf("type id empty"))
	}
	if err := validateIncidentTypeName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateIncidentTypeDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return UpdateIncidentTypeDetailsResult{}, errors.Join(errs...)
	}

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

		if err := lockClassifierOrg(tx, row.OrganizationID); err != nil {
			return err
		}

		newName := strings.TrimSpace(cmd.Name)
		newDescription := null.StringFromPtr(trimmedStringPtr(cmd.Description))
		if row.Name == newName && row.Description == newDescription {
			return nil
		}
		row.Name = newName
		row.Description = newDescription
		if err := tx.Save(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNameConflict).
					Public("An active incident type with this name already exists.").
					With("incident_type_id", row.ID).
					With("name", row.Name).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}

		if err := projector.TypeUpdateDetails(tx, &row); err != nil {
			return err
		}

		event := buildIncidentTypeDetailsChangedEvent(&row)
		return outbox.Publish(tx, SubjectIncidentTypeDetailsChanged, AggregateTypeIncidentType, row.ID.String(), row.UpdatedAt, event)
	})
	return UpdateIncidentTypeDetailsResult{}, err
}
