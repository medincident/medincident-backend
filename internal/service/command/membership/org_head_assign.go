package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
)

// AssignOrganizationHeadCommand carries the identifiers needed to link
// an employee to an organization as its head.
type AssignOrganizationHeadCommand struct {
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

// AssignOrganizationHead links the employee to the organization as its
// head. The employee must currently belong to that organization
// (organization_id is denormalized on employees). See spec §8.2 (OrgHead variant).
func (s *EmployeeService) AssignOrganizationHead(ctx context.Context, cmd AssignOrganizationHeadCommand) error {
	if err := validateOrgRoleKeys(scopeOrgHead, cmd.OrganizationID, cmd.EmployeeID); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireOrganizationExists(tx, scopeOrgHead, cmd.OrganizationID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeOrgHead, cmd.EmployeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInOrganization(scopeOrgHead, emp, cmd.OrganizationID); err != nil {
			return err
		}

		row := model.OrgHead{
			OrganizationID: cmd.OrganizationID,
			EmployeeID:     cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeOrgHead).
					Code(ErrCodeOrganizationHeadAlreadyAssigned).
					Public("This employee is already an organization head.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
		}

		return projector.OrgHeadAssigned(tx, &row)
	})
}
