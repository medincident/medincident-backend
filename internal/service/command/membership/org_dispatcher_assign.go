package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// AssignOrganizationDispatcherPayload carries the identifiers needed to link
// an employee to an organization as its dispatcher.
type AssignOrganizationDispatcherPayload struct {
	OrganizationID string `validate:"required,uuid"`
	EmployeeID     string `validate:"required,uuid"`
}

// AssignOrganizationDispatcherCommand = caller + payload.
type AssignOrganizationDispatcherCommand struct {
	Caller  authz.Caller
	Payload AssignOrganizationDispatcherPayload
}

// AssignOrganizationDispatcher links the employee to the organization as its
// dispatcher. The employee must currently belong to that organization
// (organization_id is denormalized on employees). See spec §8.2 (OrgDispatcher variant).
func (s *EmployeeService) AssignOrganizationDispatcher(ctx context.Context, cmd AssignOrganizationDispatcherCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	organizationID := uuid.MustParse(cmd.Payload.OrganizationID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(organizationID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireOrganizationExists(tx, scopeOrgDispatcher, organizationID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeOrgDispatcher, employeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInOrganization(scopeOrgDispatcher, emp, organizationID); err != nil {
			return err
		}

		row := model.OrgDispatcher{
			OrganizationID: organizationID,
			EmployeeID:     employeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeOrganizationDispatcherAlreadyAssigned).
					Public("This employee is already an organization dispatcher.").
					With("organization_id", organizationID).
					With("employee_id", employeeID).
					Wrap(err)
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherSaveFailed).Wrap(err)
		}

		return projector.OrgDispatcherAssigned(tx, &row)
	})
}
