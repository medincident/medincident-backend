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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// UpdateOrganizationDetailsPayload carries the new name and (optional)
// description for an existing organization.
type UpdateOrganizationDetailsPayload struct {
	ID          string  `validate:"required,uuid"`
	Name        string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// UpdateOrganizationDetailsCommand = caller + payload.
type UpdateOrganizationDetailsCommand struct {
	Caller  authz.Caller
	Payload UpdateOrganizationDetailsPayload
}

// UpdateDetails changes an organization's name and description. If
// neither field actually changes, returns nil without writing anything.
func (s *OrganizationService) UpdateDetails(
	ctx context.Context,
	cmd UpdateOrganizationDetailsCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(id)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var org model.Organization
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&org, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", id).
					Errorf("organization not found")
			}
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", id).
				Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Payload.Name)
		var newDesc null.String
		if cmd.Payload.Description != nil {
			newDesc = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if org.Name == newName && org.Description == newDesc {
			return nil
		}

		org.Name = newName
		org.Description = newDesc
		if err := tx.Save(&org).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", id).
				Wrap(err)
		}

		return projector.OrganizationDetailsChanged(tx, &org)
	})
}
