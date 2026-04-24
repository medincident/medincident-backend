package membership

import (
	"context"
	"errors"

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

// AssignOrganizationDispatcherDeputyPayload carries the identifiers needed
// to set the deputy slot on an existing OrgDispatcher role.
type AssignOrganizationDispatcherDeputyPayload struct {
	OrganizationID   string `validate:"required,uuid"`
	EmployeeID       string `validate:"required,uuid"`
	DeputyEmployeeID string `validate:"required,uuid"`
}

// AssignOrganizationDispatcherDeputyCommand = caller + payload.
type AssignOrganizationDispatcherDeputyCommand struct {
	Caller  authz.Caller
	Payload AssignOrganizationDispatcherDeputyPayload
}

// AssignOrganizationDispatcherDeputy sets the deputy slot on an existing
// OrgDispatcher role. See spec §8.5 (OrgDispatcher variant).
func (s *EmployeeService) AssignOrganizationDispatcherDeputy(ctx context.Context, cmd AssignOrganizationDispatcherDeputyCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	organizationID := uuid.MustParse(cmd.Payload.OrganizationID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	deputyEmployeeID := uuid.MustParse(cmd.Payload.DeputyEmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(organizationID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.OrgDispatcher
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", organizationID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeOrganizationDispatcherNotFound).
					Public("Organization dispatcher not found.").
					With("organization_id", organizationID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", deputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", deputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if deputy.OrganizationID != organizationID {
			return oops.In(scopeOrgDispatcher).
				Code(ErrCodeDeputyNotInOrganization).
				Public("Deputy employee does not belong to this organization.").
				With("deputy_employee_id", deputyEmployeeID).
				With("deputy_organization_id", deputy.OrganizationID).
				With("target_organization_id", organizationID).
				Errorf("deputy not in target organization")
		}

		if deputyEmployeeID == employeeID {
			return oops.In(scopeOrgDispatcher).
				Code(ErrCodeDeputyIsHolder).
				Public("Deputy cannot be the same employee as the role holder.").
				Errorf("deputy is holder")
		}

		if row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgDispatcher).
				Code(ErrCodeDeputyAlreadyAssigned).
				Public("Deputy slot is already occupied.").
				With("current_deputy_employee_id", row.DeputyEmployeeID.V).
				Errorf("deputy already assigned")
		}

		row.DeputyEmployeeID = null.ValueFrom(deputyEmployeeID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherSaveFailed).Wrap(err)
		}

		return projector.OrgDispatcherDeputyAssigned(tx, &row)
	})
}
