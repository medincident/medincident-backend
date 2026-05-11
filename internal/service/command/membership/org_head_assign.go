package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// AssignOrganizationHeadPayload carries the identifiers needed to link
// an employee to an organization as its head.
type AssignOrganizationHeadPayload struct {
	OrganizationID string `validate:"required,uuid"`
	EmployeeID     string `validate:"required,uuid"`
}

// AssignOrganizationHeadCommand = caller + payload.
type AssignOrganizationHeadCommand struct {
	Caller  authz.Caller
	Payload AssignOrganizationHeadPayload
}

// AssignOrganizationHead links the employee to the organization as its
// head. The employee must currently belong to that organization
// (organization_id is denormalized on employees). See spec §8.2 (OrgHead variant).
//
// See: docs/services/Membership.md
func (s *EmployeeService) AssignOrganizationHead(ctx context.Context, cmd AssignOrganizationHeadCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	organizationID := uuid.MustParse(cmd.Payload.OrganizationID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(organizationID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireOrganizationExists(tx, scopeOrgHead, organizationID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeOrgHead, employeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInOrganization(scopeOrgHead, emp, organizationID); err != nil {
			return err
		}

		row := model.OrgHead{
			OrganizationID: organizationID,
			EmployeeID:     employeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeOrgHead).
					Code(ErrCodeOrganizationHeadAlreadyAssigned).
					Public("This employee is already an organization head.").
					With("organization_id", organizationID).
					With("employee_id", employeeID).
					Wrap(err)
			}
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
		}

		env, err := buildOrgHeadAssignedEnvelope(&row)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.organization.v1.org_head_assigned", env)
	})
}

func buildOrgHeadAssignedEnvelope(row *model.OrgHead) (*eventv1.Envelope, error) {
	msg := &orgv1.OrgHeadAssigned{
		EmployeeId: row.EmployeeID.String(),
		AssignedAt: timestamppb.New(row.CreatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(row.CreatedAt),
		AggregateType: "organization",
		AggregateId:   row.OrganizationID.String(),
		Payload:       payload,
	}, nil
}
