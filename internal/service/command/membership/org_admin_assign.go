package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// AssignOrganizationAdminCommand carries the identifiers needed to link
// an employee to an organization as its admin.
type AssignOrganizationAdminCommand struct {
	OrganizationID uuid.UUID `validate:"required"`
	EmployeeID     uuid.UUID `validate:"required"`
}

// AssignOrganizationAdmin links the employee to the organization as its
// admin. The employee must currently belong to that organization
// (organization_id is denormalized on employees). See spec §8.2 (OrgAdmin variant).
func (s *EmployeeService) AssignOrganizationAdmin(ctx context.Context, cmd AssignOrganizationAdminCommand) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireOrganizationExists(tx, scopeOrgAdmin, cmd.OrganizationID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeOrgAdmin, cmd.EmployeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInOrganization(scopeOrgAdmin, emp, cmd.OrganizationID); err != nil {
			return err
		}

		row := model.OrgAdmin{
			OrganizationID: cmd.OrganizationID,
			EmployeeID:     cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeOrganizationAdminAlreadyAssigned).
					Public("This employee is already an organization admin.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
		}

		return projector.OrgAdminAssigned(tx, &row)
	})
}
