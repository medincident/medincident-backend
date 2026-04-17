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

	organizationv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/organization/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
)

// UpdateOrganizationDetailsCommand carries the new name and (optional)
// description for an existing organization.
type UpdateOrganizationDetailsCommand struct {
	ID          uuid.UUID
	Name        string
	Description *string // nil = clear description
}

// buildOrganizationDetailsChangedEvent assembles the details-changed
// event from the updated model. Only name and description are included;
// address is carried by a separate event type.
func buildOrganizationDetailsChangedEvent(org *model.Organization) *organizationv1.OrganizationDetailsChanged {
	ev := &organizationv1.OrganizationDetailsChanged{Name: org.Name}
	if org.Description.Valid {
		desc := org.Description.String
		ev.Description = &desc
	}
	return ev
}

// UpdateDetails changes an organization's name and description. If
// neither field actually changes, returns nil without writing anything.
func (s *OrganizationService) UpdateDetails(
	ctx context.Context,
	cmd UpdateOrganizationDetailsCommand,
) error {
	var errs []error
	if err := validateOrganizationName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateOrganizationDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
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
		newDesc := null.StringFromPtr(cmd.Description)
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

		event := buildOrganizationDetailsChangedEvent(&org)
		return outbox.Publish(tx, SubjectOrganizationDetailsChanged, AggregateTypeOrganization, org.ID.String(), org.UpdatedAt, event)
	})
}
