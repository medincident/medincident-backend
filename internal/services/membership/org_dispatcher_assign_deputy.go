package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	organizationv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/organization/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// AssignOrganizationDispatcherDeputyCommand carries the identifiers needed
// to set the deputy slot on an existing OrgDispatcher role.
type AssignOrganizationDispatcherDeputyCommand struct {
	OrganizationID   uuid.UUID
	EmployeeID       uuid.UUID
	DeputyEmployeeID uuid.UUID
}

// AssignOrganizationDispatcherDeputy sets the deputy slot on an existing
// OrgDispatcher role. See spec §8.5 (OrgDispatcher variant).
func (s *EmployeeService) AssignOrganizationDispatcherDeputy(ctx context.Context, cmd AssignOrganizationDispatcherDeputyCommand) error {
	var errs []error
	if cmd.OrganizationID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationIDEmpty).
			Public("Organization ID is required.").Errorf("organization id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgDispatcher).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if cmd.DeputyEmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgDispatcher).Code(ErrCodeDeputyEmployeeIDEmpty).
			Public("Deputy employee ID is required.").Errorf("deputy id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.OrgDispatcher
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", cmd.OrganizationID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeOrganizationDispatcherNotFound).
					Public("Organization dispatcher not found.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", cmd.DeputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", cmd.DeputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if deputy.OrganizationID != cmd.OrganizationID {
			return oops.In(scopeOrgDispatcher).
				Code(ErrCodeDeputyNotInOrganization).
				Public("Deputy employee does not belong to this organization.").
				With("deputy_employee_id", cmd.DeputyEmployeeID).
				With("deputy_organization_id", deputy.OrganizationID).
				With("target_organization_id", cmd.OrganizationID).
				Errorf("deputy not in target organization")
		}

		if cmd.DeputyEmployeeID == cmd.EmployeeID {
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

		deputyID := cmd.DeputyEmployeeID
		row.DeputyEmployeeID = null.ValueFrom(deputyID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherSaveFailed).Wrap(err)
		}

		ev := &organizationv1.OrganizationDispatcherDeputyAssigned{
			EmployeeId:       cmd.EmployeeID.String(),
			DeputyEmployeeId: cmd.DeputyEmployeeID.String(),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(time.Now().UTC()),
			AggregateType: AggregateTypeOrganization,
			AggregateId:   cmd.OrganizationID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectOrganizationDispatcherDeputyAssigned, envelope, nil)
	})
}
