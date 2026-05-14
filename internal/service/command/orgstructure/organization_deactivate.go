package orgstructure

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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	"github.com/medincident/medincident-backend/internal/service/validation"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// DeactivateOrganizationPayload identifies the organization to deactivate.
type DeactivateOrganizationPayload struct {
	ID string `validate:"required,uuid"`
}

// DeactivateOrganizationCommand = caller + payload.
type DeactivateOrganizationCommand struct {
	Caller  authz.Caller
	Payload DeactivateOrganizationPayload
}

// Deactivate performs cascading deactivation: the organization, all its
// clinics, all departments in those clinics, and all employees in those
// departments. Each employee loses all roles. One outbox event is written
// per row actually toggled, all inside a single transaction.
//
// See: docs/services/OrgStructure.md
func (s *OrganizationService) Deactivate(
	ctx context.Context,
	cmd DeactivateOrganizationCommand,
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

		// Collect active clinic IDs.
		var clinicIDs []uuid.UUID
		if err := tx.Raw(
			`SELECT id FROM domain.clinics WHERE organization_id = ? AND is_active FOR UPDATE`,
			orgID,
		).Scan(&clinicIDs).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", orgID).
				Wrap(err)
		}

		if len(clinicIDs) > 0 {
			// Collect active department IDs in those clinics.
			var deptIDs []uuid.UUID
			if err := tx.Raw(
				`SELECT id FROM domain.departments WHERE clinic_id = ANY(?) AND is_active FOR UPDATE`,
				clinicIDs,
			).Scan(&deptIDs).Error; err != nil {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationLoadFailed).
					With("organization_id", orgID).
					Wrap(err)
			}

			if len(deptIDs) > 0 {
				// Collect active employee IDs in those departments.
				var empIDs []uuid.UUID
				if err := tx.Raw(
					`SELECT id FROM domain.employees WHERE department_id = ANY(?) AND is_active FOR UPDATE`,
					deptIDs,
				).Scan(&empIDs).Error; err != nil {
					return oops.In("services.orgstructure.organization").
						Code(ErrCodeOrganizationLoadFailed).
						With("organization_id", orgID).
						Wrap(err)
				}

				for _, empID := range empIDs {
					if err := membership.RevokeAllEmployeeRoles(tx, empID, now); err != nil {
						return err
					}
				}

				// Deactivate employees.
				var deactivatedEmpIDs []uuid.UUID
				if err := tx.Raw(
					`UPDATE domain.employees SET is_active = FALSE, updated_at = now() WHERE id = ANY(?) AND is_active RETURNING id`,
					empIDs,
				).Scan(&deactivatedEmpIDs).Error; err != nil {
					return oops.In("services.orgstructure.organization").
						Code(ErrCodeOrganizationSaveFailed).
						With("organization_id", orgID).
						Wrap(err)
				}
				for _, id := range deactivatedEmpIDs {
					if err := appendOrgCascadeEmployeeDeactivatedEvent(tx, id, now); err != nil {
						return err
					}
				}

				// Deactivate departments.
				var deactivatedDeptIDs []uuid.UUID
				if err := tx.Raw(
					`UPDATE domain.departments SET is_active = FALSE, updated_at = now() WHERE id = ANY(?) AND is_active RETURNING id`,
					deptIDs,
				).Scan(&deactivatedDeptIDs).Error; err != nil {
					return oops.In("services.orgstructure.organization").
						Code(ErrCodeOrganizationSaveFailed).
						With("organization_id", orgID).
						Wrap(err)
				}
				for _, id := range deactivatedDeptIDs {
					if err := appendOrgCascadeDepartmentDeactivatedEvent(tx, id, now); err != nil {
						return err
					}
				}
			}

			// Deactivate clinics.
			var deactivatedClinicIDs []uuid.UUID
			if err := tx.Raw(
				`UPDATE domain.clinics SET is_active = FALSE, updated_at = now() WHERE id = ANY(?) AND is_active RETURNING id`,
				clinicIDs,
			).Scan(&deactivatedClinicIDs).Error; err != nil {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationSaveFailed).
					With("organization_id", orgID).
					Wrap(err)
			}
			for _, id := range deactivatedClinicIDs {
				if err := appendOrgCascadeClinicDeactivatedEvent(tx, id, now); err != nil {
					return err
				}
			}
		}

		// Deactivate the organization itself (idempotent).
		var updatedAt time.Time
		if err := tx.Raw(
			`UPDATE domain.organizations SET is_active = FALSE, updated_at = now() WHERE id = ? RETURNING updated_at`,
			orgID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", orgID).
				Wrap(err)
		}
		return appendOrganizationDeactivatedEvent(tx, orgID, updatedAt)
	})
}

func appendOrganizationDeactivatedEvent(tx *gorm.DB, orgID uuid.UUID, now time.Time) error {
	msg := &orgv1.OrganizationDeactivated{OrganizationId: orgID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.organization").Code(ErrCodeOrganizationSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "organization",
		AggregateId:   orgID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.organization.v1.deactivated", env)
}

func appendOrgCascadeClinicDeactivatedEvent(tx *gorm.DB, clinicID uuid.UUID, now time.Time) error {
	msg := &clinicv1.ClinicDeactivated{ClinicId: clinicID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.organization").Code(ErrCodeOrganizationSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "clinic",
		AggregateId:   clinicID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.clinic.v1.deactivated", env)
}

func appendOrgCascadeDepartmentDeactivatedEvent(tx *gorm.DB, deptID uuid.UUID, now time.Time) error {
	msg := &deptv1.DepartmentDeactivated{DepartmentId: deptID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.organization").Code(ErrCodeOrganizationSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "department",
		AggregateId:   deptID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.department.v1.deactivated", env)
}

func appendOrgCascadeEmployeeDeactivatedEvent(tx *gorm.DB, empID uuid.UUID, now time.Time) error {
	msg := &empv1.EmployeeDeactivated{EmployeeId: empID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.organization").Code(ErrCodeOrganizationSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "employee",
		AggregateId:   empID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.employee.v1.deactivated", env)
}
