package incident

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// CreateIncidentPayload is the validated client-facing payload.
type CreateIncidentPayload struct {
	DepartmentID string    `validate:"required,uuid"`
	CategoryID   string    `validate:"required,uuid"`
	TypeID       string    `validate:"required,uuid"`
	Description  *string   `validate:"omitnil,no_extra_ws,min=1,max=10000"`
	OccurredAt   time.Time `validate:"required"`
}

type CreateIncidentCommand struct {
	Caller  authz.Caller
	Payload CreateIncidentPayload
}

type CreateIncidentResult struct {
	ID uuid.UUID
}

// Create persists a new incident and the initial status history row.
// Authorization: caller must be an employee of the department's
// organization (MemberOf).
//
// See: docs/services/incident/Incidents.md
func (s *IncidentService) Create(
	ctx context.Context, cmd *CreateIncidentCommand,
) (CreateIncidentResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateIncidentResult{}, err
	}
	deptID := uuid.MustParse(cmd.Payload.DepartmentID)
	categoryID := uuid.MustParse(cmd.Payload.CategoryID)
	typeID := uuid.MustParse(cmd.Payload.TypeID)

	now := time.Now()
	if cmd.Payload.OccurredAt.After(now) {
		return CreateIncidentResult{}, oops.In(scope).
			Code(ErrCodeIncidentOccurredAtFuture).
			Public("occurred_at must not be in the future.").
			With("occurred_at", cmd.Payload.OccurredAt).
			Errorf("future occurred_at")
	}
	if now.Sub(cmd.Payload.OccurredAt) > incidentMaxOccurredAtAge {
		return CreateIncidentResult{}, oops.In(scope).
			Code(ErrCodeIncidentOccurredAtTooOld).
			Public("occurred_at exceeds the 48h registration window.").
			With("occurred_at", cmd.Payload.OccurredAt).
			With("max_age_hours", int(incidentMaxOccurredAtAge.Hours())).
			Errorf("occurred_at too old")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateIncidentResult{}, oops.In(scope).
			Code(ErrCodeIncidentIDGenerationFailed).
			Public("Failed to create incident.").
			Wrap(err)
	}

	var result CreateIncidentResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Resolve dept -> clinic -> organization in a single JOIN.
		var deptClinic struct {
			ClinicID       uuid.UUID
			OrganizationID uuid.UUID
		}
		if err := tx.Raw(`SELECT c.id AS clinic_id, c.organization_id
			FROM domain.departments d
			JOIN domain.clinics c ON c.id = d.clinic_id
			WHERE d.id = ?`, deptID).Scan(&deptClinic).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentLoadFailed).Wrap(err)
		}
		if deptClinic.ClinicID == uuid.Nil {
			return oops.In(scope).
				Code(ErrCodeIncidentDeptNotFound).
				Public("Department not found.").
				With("department_id", deptID).
				Errorf("department not found")
		}
		orgID := deptClinic.OrganizationID

		// Authorization: caller is an employee of this org.
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			authz.MemberOf.Organization(orgID)); err != nil {
			return err
		}

		// Validate category + type belong to org and are active.
		var cat model.IncidentCategory
		if err := tx.First(&cat, "id = ?", categoryID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentCategoryNotFound).
				With("category_id", categoryID).Wrap(err)
		}
		if cat.OrganizationID != orgID {
			return oops.In(scope).Code(ErrCodeIncidentCategoryNotFound).
				Public("Category does not belong to the incident's organization.").
				With("category_id", categoryID).Errorf("org mismatch")
		}
		if !cat.IsActive {
			return oops.In(scope).Code(ErrCodeIncidentCategoryInactive).
				Public("Category is inactive.").
				With("category_id", categoryID).Errorf("inactive")
		}
		var typ model.IncidentType
		if err := tx.First(&typ, "id = ?", typeID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentTypeNotFound).
				With("type_id", typeID).Wrap(err)
		}
		if typ.OrganizationID != orgID {
			return oops.In(scope).Code(ErrCodeIncidentTypeOrgMismatch).
				Public("Type does not belong to the incident's organization.").
				With("type_id", typeID).Errorf("org mismatch")
		}
		if typ.CategoryID != categoryID {
			return oops.In(scope).Code(ErrCodeIncidentTypeCategoryMismatch).
				Public("Type does not belong to the chosen category.").
				With("type_id", typeID).With("category_id", categoryID).
				Errorf("category mismatch")
		}
		if !typ.IsActive {
			return oops.In(scope).Code(ErrCodeIncidentTypeInactive).
				Public("Type is inactive.").
				With("type_id", typeID).Errorf("inactive")
		}

		// Resolve registrar employee + display_name.
		reg, err := s.loadRegistrarSnapshot(tx, cmd.Caller.ZitadelUserID, orgID)
		if err != nil {
			return err
		}

		desc := null.String{}
		if cmd.Payload.Description != nil {
			desc = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}

		inc := model.Incident{
			ID:                  id,
			OrganizationID:      orgID,
			ClinicID:            deptClinic.ClinicID,
			DepartmentID:        deptID,
			CategoryID:          categoryID,
			TypeID:              typeID,
			Status:              model.IncidentStatusPending,
			Priority:            model.IncidentPriorityNormal,
			Description:         desc,
			OccurredAt:          cmd.Payload.OccurredAt,
			RegistrarEmployeeID: reg.EmployeeID,
			CreatedAt:           now,
			UpdatedAt:           now,
		}
		if err := tx.Create(&inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).
				With("incident_id", id).Wrap(err)
		}
		if err := projector.IncidentCreated(tx, &inc, &reg); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}

// loadRegistrarSnapshot reads the caller's employee row in the given
// org and the matching projections.users display_name.
func (s *IncidentService) loadRegistrarSnapshot(
	tx *gorm.DB, callerZitadelID string, orgID uuid.UUID,
) (projector.IncidentRegistrarSnapshot, error) {
	var snap struct {
		EmployeeID     uuid.UUID
		Position       null.String
		OrganizationID uuid.UUID
		ClinicID       uuid.UUID
		DepartmentID   uuid.UUID
	}
	if err := tx.Raw(`SELECT e.id AS employee_id, e.position, e.organization_id,
			d.clinic_id, e.department_id
		FROM domain.employees e
		JOIN domain.departments d ON d.id = e.department_id
		WHERE e.zitadel_user_id = ? AND e.organization_id = ?`,
		callerZitadelID, orgID).Scan(&snap).Error; err != nil {
		return projector.IncidentRegistrarSnapshot{}, oops.In(scope).
			Code(ErrCodeIncidentEmployeeNotFound).
			Public("Caller is not an employee of this organization.").
			With("organization_id", orgID).
			Wrap(err)
	}
	if snap.EmployeeID == uuid.Nil {
		return projector.IncidentRegistrarSnapshot{}, oops.In(scope).
			Code(ErrCodeIncidentEmployeeNotFound).
			Public("Caller is not an employee of this organization.").
			With("organization_id", orgID).
			Errorf("employee not found")
	}
	// projections.users.display_name is the canonical denormalized name.
	var displayName string
	if err := tx.Raw(
		`SELECT display_name FROM projections.users WHERE id = ?`,
		callerZitadelID,
	).Scan(&displayName).Error; err != nil || displayName == "" {
		return projector.IncidentRegistrarSnapshot{}, oops.In(scope).
			Code(ErrCodeIncidentRegistrarUserNotFound).
			Public("Registrar user record is missing.").
			With("zitadel_user_id", callerZitadelID).
			Errorf("display_name lookup failed")
	}
	return projector.IncidentRegistrarSnapshot{
		EmployeeID:     snap.EmployeeID,
		DisplayName:    displayName,
		Position:       snap.Position,
		OrganizationID: snap.OrganizationID,
		ClinicID:       snap.ClinicID,
		DepartmentID:   snap.DepartmentID,
	}, nil
}
