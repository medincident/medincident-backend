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
	"github.com/medincident/medincident-command-service/internal/pgerr"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// AssignOrganizationAdminCommand carries the identifiers needed to link
// an employee to an organization as its admin.
type AssignOrganizationAdminCommand struct {
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

// AssignOrganizationAdmin links the employee to the organization as its
// admin. The employee must currently belong to that organization
// (organization_id is denormalized on employees). See spec §8.2 (OrgAdmin variant).
func (s *EmployeeService) AssignOrganizationAdmin(ctx context.Context, cmd AssignOrganizationAdminCommand) error {
	var errs []error
	if cmd.OrganizationID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgAdmin).
			Code(ErrCodeOrganizationIDEmpty).
			Public("Organization ID is required.").
			Errorf("organization id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgAdmin).
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
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationLookupFailed).Wrap(err)
		}
		if orgExists == 0 {
			return oops.In(scopeOrgAdmin).
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
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", cmd.EmployeeID).
					Errorf("employee not found")
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if emp.OrganizationID != cmd.OrganizationID {
			return oops.In(scopeOrgAdmin).
				Code(ErrCodeEmployeeNotInOrganization).
				Public("Employee does not belong to this organization.").
				With("employee_id", cmd.EmployeeID).
				With("employee_organization_id", emp.OrganizationID).
				With("target_organization_id", cmd.OrganizationID).
				Errorf("employee not in target organization")
		}

		row := model.OrgAdmin{
			OrganizationID: cmd.OrganizationID,
			EmployeeID:     cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerr.CodeUniqueViolation {
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeOrganizationAdminAlreadyAssigned).
					Public("This employee is already an organization admin.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
		}

		ev := &organizationv1.OrganizationAdminAssigned{
			EmployeeId: cmd.EmployeeID.String(),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(row.CreatedAt),
			AggregateType: AggregateTypeOrganization,
			AggregateId:   cmd.OrganizationID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectOrganizationAdminAssigned, envelope)
	})
}
