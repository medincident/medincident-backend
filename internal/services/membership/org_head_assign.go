package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

const scopeOrgHead = "services.membership.org_head"

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
	var errs []error
	if cmd.OrganizationID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgHead).
			Code(ErrCodeOrganizationIDEmpty).
			Public("Organization ID is required.").
			Errorf("organization id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgHead).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var orgExists int64
		if err := tx.Raw(`SELECT count(*) FROM domain.organizations WHERE id = ?`, cmd.OrganizationID).
			Scan(&orgExists).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadLoadFailed).Wrap(err)
		}
		if orgExists == 0 {
			return oops.In(scopeOrgHead).
				Code(ErrCodeOrganizationNotFound).
				Public("Organization not found.").
				With("organization_id", cmd.OrganizationID).
				Errorf("organization not found")
		}

		var emp model.Employee
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", cmd.EmployeeID).
			First(&emp).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgHead).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", cmd.EmployeeID).
					Errorf("employee not found")
			}
			return oops.In(scopeOrgHead).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if emp.OrganizationID != cmd.OrganizationID {
			return oops.In(scopeOrgHead).
				Code(ErrCodeEmployeeNotInOrganization).
				Public("Employee does not belong to this organization.").
				With("employee_id", cmd.EmployeeID).
				With("employee_organization_id", emp.OrganizationID).
				With("target_organization_id", cmd.OrganizationID).
				Errorf("employee not in target organization")
		}

		row := model.OrgHead{
			OrganizationID: cmd.OrganizationID,
			EmployeeID:     cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
				return oops.In(scopeOrgHead).
					Code(ErrCodeOrganizationHeadAlreadyAssigned).
					Public("This employee is already an organization head.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
		}

		ev := &organizationv1.OrganizationHeadAssigned{
			EmployeeId: cmd.EmployeeID.String(),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(row.CreatedAt),
			AggregateType: AggregateTypeOrganization,
			AggregateId:   cmd.OrganizationID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectOrganizationHeadAssigned, envelope, nil)
	})
}
