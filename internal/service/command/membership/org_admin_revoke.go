package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

// RevokeOrganizationAdminPayload carries the identifiers needed to
// remove an employee's organization admin role.
type RevokeOrganizationAdminPayload struct {
	OrganizationID string `validate:"required,uuid"`
	EmployeeID     string `validate:"required,uuid"`
}

// RevokeOrganizationAdminCommand = caller + payload.
type RevokeOrganizationAdminCommand struct {
	Caller  authz.Caller
	Payload RevokeOrganizationAdminPayload
}

// RevokeOrganizationAdmin removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
//
// See: docs/services/Membership.md
func (s *EmployeeService) RevokeOrganizationAdmin(ctx context.Context, cmd RevokeOrganizationAdminCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	organizationID := uuid.MustParse(cmd.Payload.OrganizationID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(organizationID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishOrgAdminDeputyRemoved(tx, organizationID, employeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.OrgAdmin{}, "organization_id = ? AND employee_id = ?",
			organizationID, employeeID).Error; err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminDeleteFailed).Wrap(err)
		}

		return publishOrgAdminRevoked(tx, organizationID, employeeID, now)
	})
}

// publishOrgAdminRevoked is a shared helper for Revoke and
// cascade-on-terminate. It does NOT delete the domain row; the caller
// owns that.
func publishOrgAdminRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	env, err := buildOrgAdminRevokedEnvelope(organizationID, employeeID, now)
	if err != nil {
		return err
	}
	return outbox.Append(tx, "medincident.event.organization.v1.org_admin_revoked", env)
}

// publishOrgAdminDeputyRemoved is a shared helper; it clears the
// projection's deputy slot. Caller owns the domain-row update.
func publishOrgAdminDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	env, err := buildOrgAdminDeputyRemovedEnvelope(organizationID, employeeID, now)
	if err != nil {
		return err
	}
	return outbox.Append(tx, "medincident.event.organization.v1.org_admin_deputy_removed", env)
}

func buildOrgAdminRevokedEnvelope(organizationID, employeeID uuid.UUID, now time.Time) (*eventv1.Envelope, error) {
	msg := &orgv1.OrgAdminRevoked{EmployeeId: employeeID.String()}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminDeleteFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "organization",
		AggregateId:   organizationID.String(),
		Payload:       payload,
	}, nil
}

func buildOrgAdminDeputyRemovedEnvelope(organizationID, employeeID uuid.UUID, now time.Time) (*eventv1.Envelope, error) {
	msg := &orgv1.OrgAdminDeputyRemoved{EmployeeId: employeeID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "organization",
		AggregateId:   organizationID.String(),
		Payload:       payload,
	}, nil
}
