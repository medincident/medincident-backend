package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"gorm.io/gorm"

	organizationv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/organization/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/pgerr"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
)

// AssignOrganizationDispatcherCommand carries the identifiers needed to link
// an employee to an organization as its dispatcher.
type AssignOrganizationDispatcherCommand struct {
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

// AssignOrganizationDispatcher links the employee to the organization as its
// dispatcher. The employee must currently belong to that organization
// (organization_id is denormalized on employees). See spec §8.2 (OrgDispatcher variant).
func (s *EmployeeService) AssignOrganizationDispatcher(ctx context.Context, cmd AssignOrganizationDispatcherCommand) error {
	if err := validateOrgRoleKeys(scopeOrgDispatcher, cmd.OrganizationID, cmd.EmployeeID); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireOrganizationExists(tx, scopeOrgDispatcher, cmd.OrganizationID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeOrgDispatcher, cmd.EmployeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInOrganization(scopeOrgDispatcher, emp, cmd.OrganizationID); err != nil {
			return err
		}

		row := model.OrgDispatcher{
			OrganizationID: cmd.OrganizationID,
			EmployeeID:     cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerr.CodeUniqueViolation {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeOrganizationDispatcherAlreadyAssigned).
					Public("This employee is already an organization dispatcher.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherSaveFailed).Wrap(err)
		}

		ev := &organizationv1.OrganizationDispatcherAssigned{
			EmployeeId: cmd.EmployeeID.String(),
		}
		return outbox.Publish(tx, SubjectOrganizationDispatcherAssigned, AggregateTypeOrganization, cmd.OrganizationID.String(), row.CreatedAt, ev)
	})
}
