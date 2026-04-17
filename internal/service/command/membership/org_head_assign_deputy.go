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

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/outbox"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	organizationv1 "github.com/medincident/medincident-command-service/pkg/event/organization/v1"
)

// AssignOrganizationHeadDeputyCommand carries the identifiers needed
// to set the deputy slot on an existing OrgHead role.
type AssignOrganizationHeadDeputyCommand struct {
	OrganizationID   uuid.UUID
	EmployeeID       uuid.UUID
	DeputyEmployeeID uuid.UUID
}

// AssignOrganizationHeadDeputy sets the deputy slot on an existing
// OrgHead role. See spec §8.5 (OrgHead variant).
func (s *EmployeeService) AssignOrganizationHeadDeputy(ctx context.Context, cmd AssignOrganizationHeadDeputyCommand) error {
	var errs []error
	if cmd.OrganizationID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgHead).Code(ErrCodeOrganizationIDEmpty).
			Public("Organization ID is required.").Errorf("organization id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgHead).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if cmd.DeputyEmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgHead).Code(ErrCodeDeputyEmployeeIDEmpty).
			Public("Deputy employee ID is required.").Errorf("deputy id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var row model.OrgHead
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", cmd.OrganizationID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgHead).
					Code(ErrCodeOrganizationHeadNotFound).
					Public("Organization head not found.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", cmd.DeputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgHead).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", cmd.DeputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeOrgHead).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if deputy.OrganizationID != cmd.OrganizationID {
			return oops.In(scopeOrgHead).
				Code(ErrCodeDeputyNotInOrganization).
				Public("Deputy employee does not belong to this organization.").
				With("deputy_employee_id", cmd.DeputyEmployeeID).
				With("deputy_organization_id", deputy.OrganizationID).
				With("target_organization_id", cmd.OrganizationID).
				Errorf("deputy not in target organization")
		}

		if cmd.DeputyEmployeeID == cmd.EmployeeID {
			return oops.In(scopeOrgHead).
				Code(ErrCodeDeputyIsHolder).
				Public("Deputy cannot be the same employee as the role holder.").
				Errorf("deputy is holder")
		}

		if row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgHead).
				Code(ErrCodeDeputyAlreadyAssigned).
				Public("Deputy slot is already occupied.").
				With("current_deputy_employee_id", row.DeputyEmployeeID.V).
				Errorf("deputy already assigned")
		}

		deputyID := cmd.DeputyEmployeeID
		row.DeputyEmployeeID = null.ValueFrom(deputyID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
		}

		if err := projector.OrgHeadDeputyAssigned(tx, &row); err != nil {
			return err
		}

		ev := &organizationv1.OrganizationHeadDeputyAssigned{
			EmployeeId:       cmd.EmployeeID.String(),
			DeputyEmployeeId: cmd.DeputyEmployeeID.String(),
		}
		return outbox.Publish(tx, SubjectOrganizationHeadDeputyAssigned, AggregateTypeOrganization, cmd.OrganizationID.String(), now, ev)
	})
}
