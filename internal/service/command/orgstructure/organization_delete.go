package orgstructure

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	"github.com/medincident/medincident-backend/internal/service/validation"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// DeleteOrganizationPayload identifies the organization to delete.
type DeleteOrganizationPayload struct {
	ID string `validate:"required,uuid"`
}

// DeleteOrganizationCommand = caller + payload.
type DeleteOrganizationCommand struct {
	Caller  authz.Caller
	Payload DeleteOrganizationPayload
}

// Delete hard-deletes an organization and all its children in the correct
// FK order: roles → employees → departments → clinics → organization.
// If the organization has dependents (incidents, service_requests) the DB
// raises a FK violation which is surfaced as organization_delete_has_dependents.
//
// See: docs/services/OrgStructure.md
func (s *OrganizationService) Delete(
	ctx context.Context,
	cmd DeleteOrganizationCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	orgID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var org model.Organization
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&org, "id = ?", orgID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", orgID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", orgID).
				Wrap(err)
		}

		// Lock all employees in this org and collect their IDs.
		var empIDs []uuid.UUID
		if err := tx.Raw(
			`SELECT id FROM domain.employees WHERE organization_id = ? FOR UPDATE`,
			orgID,
		).Scan(&empIDs).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", orgID).
				Wrap(err)
		}

		// Revoke all roles, then hard-delete each employee.
		for _, empID := range empIDs {
			if err := membership.RevokeAllEmployeeRoles(tx, empID, now); err != nil {
				return err
			}
		}
		if len(empIDs) > 0 {
			if err := tx.Exec(
				`DELETE FROM domain.employees WHERE id = ANY(?)`, empIDs,
			).Error; err != nil {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationDeleteFailed).
					With("organization_id", orgID).
					Wrap(err)
			}
		}

		// Delete departments.
		if err := tx.Exec(
			`DELETE FROM domain.departments WHERE clinic_id IN (SELECT id FROM domain.clinics WHERE organization_id = ?)`,
			orgID,
		).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationDeleteFailed).
				With("organization_id", orgID).
				Wrap(err)
		}

		// Delete clinics.
		if err := tx.Exec(
			`DELETE FROM domain.clinics WHERE organization_id = ?`, orgID,
		).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationDeleteFailed).
				With("organization_id", orgID).
				Wrap(err)
		}

		// Delete the organization itself.
		if err := tx.Delete(&model.Organization{}, "id = ?", orgID).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationDeleteHasDependents).
					Public("Organization has dependents and cannot be deleted.").
					With("organization_id", orgID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationDeleteFailed).
				With("organization_id", orgID).
				Wrap(err)
		}

		return appendOrganizationDeletedEvent(tx, orgID, now)
	})
}

func appendOrganizationDeletedEvent(tx *gorm.DB, orgID uuid.UUID, now time.Time) error {
	msg := &orgv1.OrganizationDeleted{OrganizationId: orgID.String(), DeletedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.organization").Code(ErrCodeOrganizationDeleteFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "organization",
		AggregateId:   orgID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.organization.v1.deleted", env)
}
