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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
)

type UpdateIncidentTypeDetailsCommand struct {
	TypeID      uuid.UUID
	Name        string
	Description *string
}

type UpdateIncidentTypeDetailsResult struct{}

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

		return projector.TypeUpdateDetails(tx, &row)
	})
	return UpdateIncidentTypeDetailsResult{}, err
}
