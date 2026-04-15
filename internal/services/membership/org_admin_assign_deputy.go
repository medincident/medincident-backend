package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	organizationv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/organization/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// AssignOrganizationAdminDeputyCommand carries the identifiers needed
// to set the deputy slot on an existing OrgAdmin role.
type AssignOrganizationAdminDeputyCommand struct {
	OrganizationID   uuid.UUID
	EmployeeID       uuid.UUID
	DeputyEmployeeID uuid.UUID
}

// AssignOrganizationAdminDeputy sets the deputy slot on an existing
// OrgAdmin role. See spec §8.5 (OrgAdmin variant).
func (s *EmployeeService) AssignOrganizationAdminDeputy(ctx context.Context, cmd AssignOrganizationAdminDeputyCommand) error {
	var errs []error
	if cmd.OrganizationID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationIDEmpty).
			Public("Organization ID is required.").Errorf("organization id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgAdmin).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if cmd.DeputyEmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgAdmin).Code(ErrCodeDeputyEmployeeIDEmpty).
			Public("Deputy employee ID is required.").Errorf("deputy id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.OrgAdmin
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", cmd.OrganizationID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeOrganizationAdminNotFound).
					Public("Organization admin not found.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", cmd.DeputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", cmd.DeputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if deputy.OrganizationID != cmd.OrganizationID {
			return oops.In(scopeOrgAdmin).
				Code(ErrCodeDeputyNotInOrganization).
				Public("Deputy employee does not belong to this organization.").
				With("deputy_employee_id", cmd.DeputyEmployeeID).
				With("deputy_organization_id", deputy.OrganizationID).
				With("target_organization_id", cmd.OrganizationID).
				Errorf("deputy not in target organization")
		}

		if cmd.DeputyEmployeeID == cmd.EmployeeID {
			return oops.In(scopeOrgAdmin).
				Code(ErrCodeDeputyIsHolder).
				Public("Deputy cannot be the same employee as the role holder.").
				Errorf("deputy is holder")
		}

		if row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgAdmin).
				Code(ErrCodeDeputyAlreadyAssigned).
				Public("Deputy slot is already occupied.").
				With("current_deputy_employee_id", row.DeputyEmployeeID.V).
				Errorf("deputy already assigned")
		}

		deputyID := cmd.DeputyEmployeeID
		row.DeputyEmployeeID = null.ValueFrom(deputyID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
		}

		ev := &organizationv1.OrganizationAdminDeputyAssigned{
			EmployeeId:       cmd.EmployeeID.String(),
			DeputyEmployeeId: cmd.DeputyEmployeeID.String(),
		}
		return outbox.Publish(tx, SubjectOrganizationAdminDeputyAssigned, AggregateTypeOrganization, cmd.OrganizationID.String(), time.Now().UTC(), ev)
	})
}
