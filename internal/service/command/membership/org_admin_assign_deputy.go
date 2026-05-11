package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// AssignOrganizationAdminDeputyPayload carries the identifiers needed
// to set the deputy slot on an existing OrgAdmin role.
type AssignOrganizationAdminDeputyPayload struct {
	OrganizationID   string `validate:"required,uuid"`
	EmployeeID       string `validate:"required,uuid"`
	DeputyEmployeeID string `validate:"required,uuid"`
}

// AssignOrganizationAdminDeputyCommand = caller + payload.
type AssignOrganizationAdminDeputyCommand struct {
	Caller  authz.Caller
	Payload AssignOrganizationAdminDeputyPayload
}

// AssignOrganizationAdminDeputy sets the deputy slot on an existing
// OrgAdmin role. See spec §8.5 (OrgAdmin variant).
//
// See: docs/services/Membership.md
func (s *EmployeeService) AssignOrganizationAdminDeputy(ctx context.Context, cmd AssignOrganizationAdminDeputyCommand) error {
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
		var row model.OrgAdmin
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", organizationID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeOrganizationAdminNotFound).
					Public("Organization admin not found.").
					With("organization_id", organizationID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", deputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", deputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if deputy.OrganizationID != organizationID {
			return oops.In(scopeOrgAdmin).
				Code(ErrCodeDeputyNotInOrganization).
				Public("Deputy employee does not belong to this organization.").
				With("deputy_employee_id", deputyEmployeeID).
				With("deputy_organization_id", deputy.OrganizationID).
				With("target_organization_id", organizationID).
				Errorf("deputy not in target organization")
		}

		if deputyEmployeeID == employeeID {
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

		row.DeputyEmployeeID = null.ValueFrom(deputyEmployeeID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
		}

		env, err := buildOrgAdminDeputyAssignedEnvelope(&row)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.organization.v1.org_admin_deputy_assigned", env)
	})
}

func buildOrgAdminDeputyAssignedEnvelope(row *model.OrgAdmin) (*eventv1.Envelope, error) {
	msg := &orgv1.OrgAdminDeputyAssigned{
		EmployeeId:       row.EmployeeID.String(),
		DeputyEmployeeId: row.DeputyEmployeeID.V.String(),
		UpdatedAt:        timestamppb.New(row.UpdatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(row.UpdatedAt),
		AggregateType: "organization",
		AggregateId:   row.OrganizationID.String(),
		Payload:       payload,
	}, nil
}
