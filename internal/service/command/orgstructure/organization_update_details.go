package orgstructure

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
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// UpdateOrganizationDetailsCommand carries the new name and (optional)
// description for an existing organization.
type UpdateOrganizationDetailsCommand struct {
	ID          uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=4,max=256"`
	Description *string   `validate:"omitnil,min=8,max=2048"`
}

// UpdateDetails changes an organization's name and description. If
// neither field actually changes, returns nil without writing anything.
func (s *OrganizationService) UpdateDetails(
	ctx context.Context,
	cmd UpdateOrganizationDetailsCommand,
) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var org model.Organization
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&org, "id = ?", cmd.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", cmd.ID).
					Errorf("organization not found")
			}
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", cmd.ID).
				Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Name)
		var newDesc null.String
		if cmd.Description != nil {
			newDesc = null.StringFrom(strings.TrimSpace(*cmd.Description))
		}
		if org.Name == newName && org.Description == newDesc {
			return nil
		}

		org.Name = newName
		org.Description = newDesc
		if err := tx.Save(&org).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", cmd.ID).
				Wrap(err)
		}

		return projector.OrganizationDetailsChanged(tx, &org)
	})
}
