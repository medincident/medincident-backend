package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

// RevokeOrganizationAdminCommand carries the identifiers needed to
// remove an employee's organization admin role.
type RevokeOrganizationAdminCommand struct {
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

// RevokeOrganizationAdmin removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
func (s *EmployeeService) RevokeOrganizationAdmin(ctx context.Context, cmd RevokeOrganizationAdminCommand) error {
	var errs []error
	if cmd.OrganizationID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationIDEmpty).
			Public("Organization ID is required.").Errorf("organization id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeOrgAdmin).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishOrgAdminDeputyRemoved(tx, cmd.OrganizationID, cmd.EmployeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.OrgAdmin{}, "organization_id = ? AND employee_id = ?",
			cmd.OrganizationID, cmd.EmployeeID).Error; err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminDeleteFailed).Wrap(err)
		}

		return publishOrgAdminRevoked(tx, cmd.OrganizationID, cmd.EmployeeID, now)
	})
}

// publishOrgAdminRevoked is a shared helper for Revoke,
// cascade-on-terminate. It does NOT delete the row; the caller is
// responsible for that.
func publishOrgAdminRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	ev := &organizationv1.OrganizationAdminRevoked{
		EmployeeId: employeeID.String(),
	}
	payload, err := anypb.New(ev)
	if err != nil {
		return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminEventBuildFailed).Wrap(err)
	}
	envelope := &envelopev1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: AggregateTypeOrganization,
		AggregateId:   organizationID.String(),
		Payload:       payload,
	}
	return outbox.AppendEvent(tx, SubjectOrganizationAdminRevoked, envelope, nil)
}

// publishOrgAdminDeputyRemoved is a shared helper; it publishes the
// event only and does NOT update the row.
func publishOrgAdminDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	ev := &organizationv1.OrganizationAdminDeputyRemoved{
		EmployeeId: employeeID.String(),
	}
	payload, err := anypb.New(ev)
	if err != nil {
		return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminEventBuildFailed).Wrap(err)
	}
	envelope := &envelopev1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: AggregateTypeOrganization,
		AggregateId:   organizationID.String(),
		Payload:       payload,
	}
	return outbox.AppendEvent(tx, SubjectOrganizationAdminDeputyRemoved, envelope, nil)
}
