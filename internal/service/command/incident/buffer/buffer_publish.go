package buffer

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

type PublishPayload struct {
	BufferID     string  `validate:"required,uuid"`
	DepartmentID string  `validate:"required,uuid"`
	CategoryID   string  `validate:"required,uuid"`
	TypeID       string  `validate:"required,uuid"`
	Description  *string `validate:"omitnil,no_extra_ws,min=1,max=10000"`
}

type PublishCommand struct {
	Caller  authz.Caller
	Payload PublishPayload
}

type PublishResult struct {
	IncidentID uuid.UUID
}

// Publish materialises a buffer entry into a real incident. The
// dispatcher selects the department, finalises category + type,
// optionally rewrites the description, and the original patient
// description is preserved on the incident as patient_original_description.
//
// See: docs/services/incident/Buffer.md
func (s *BufferService) Publish(
	ctx context.Context, cmd *PublishCommand,
) (PublishResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return PublishResult{}, err
	}
	bufID := uuid.MustParse(cmd.Payload.BufferID)
	deptID := uuid.MustParse(cmd.Payload.DepartmentID)
	categoryID := uuid.MustParse(cmd.Payload.CategoryID)
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	now := time.Now()

	incidentID, err := uuid.NewV7()
	if err != nil {
		return PublishResult{}, oops.In(scope).
			Code(ErrCodeBufferIDGenerationFailed).Wrap(err)
	}

	var result PublishResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Load buffer to read organization_id, then authorize BEFORE any
		// state / cross-row checks so an unauthorized caller cannot
		// distinguish pending vs non-pending vs not-found vs org-mismatch.
		b, err := s.loadBuffer(tx, bufID)
		if err != nil {
			return err
		}
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			authz.AnyOf(
				authz.SystemAdmin,
				authz.OrgAdminOf.Organization(b.OrganizationID),
				authz.OrgDispatcherOf.Organization(b.OrganizationID),
			),
		); err != nil {
			return err
		}
		if b.Status != model.BufferStatusPending {
			return oops.In(scope).Code(ErrCodeBufferNotPending).
				Public("Buffer entry is not pending.").
				With("status", b.Status).Errorf("not pending")
		}

		// Department -> clinic -> organization in a single JOIN.
		var deptClinic struct {
			ClinicID       uuid.UUID
			OrganizationID uuid.UUID
		}
		if err := tx.Raw(`SELECT c.id AS clinic_id, c.organization_id
			FROM domain.departments d
			JOIN domain.clinics c ON c.id = d.clinic_id
			WHERE d.id = ?`, deptID).Scan(&deptClinic).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferLoadFailed).Wrap(err)
		}
		if deptClinic.ClinicID == uuid.Nil {
			return oops.In(scope).Code(ErrCodeBufferDeptNotFound).
				Public("Department not found.").
				With("department_id", deptID).
				Errorf("department not found")
		}
		if deptClinic.OrganizationID != b.OrganizationID {
			return oops.In(scope).Code(ErrCodeBufferDeptNotFound).
				Public("Department does not belong to this organization.").
				With("department_id", deptID).Errorf("org mismatch")
		}

		// Validate dispatcher's chosen category/type (NOT restricted to patient-allowed).
		var cat model.IncidentCategory
		if err := tx.First(&cat, "id = ?", categoryID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferCategoryNotFound).Wrap(err)
		}
		if cat.OrganizationID != b.OrganizationID || !cat.IsActive {
			return oops.In(scope).Code(ErrCodeBufferCategoryNotFound).
				Public("Category is not available.").Errorf("invalid category")
		}
		var typ model.IncidentType
		if err := tx.First(&typ, "id = ?", typeID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferTypeNotFound).Wrap(err)
		}
		if typ.OrganizationID != b.OrganizationID || !typ.IsActive ||
			typ.CategoryID != categoryID {
			return oops.In(scope).Code(ErrCodeBufferTypeNotFound).
				Public("Type is not available for this category.").
				Errorf("invalid type")
		}

		// Resolve dispatcher snapshot for registrar denormalisation.
		reg, err := s.loadDispatcherSnapshot(tx, cmd.Caller.ZitadelUserID, b.OrganizationID)
		if err != nil {
			return err
		}

		// Description: dispatcher-supplied if present, else patient's text.
		desc := b.Description
		if cmd.Payload.Description != nil {
			desc = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}

		inc := model.Incident{
			ID:                         incidentID,
			OrganizationID:             b.OrganizationID,
			ClinicID:                   deptClinic.ClinicID,
			DepartmentID:               deptID,
			CategoryID:                 categoryID,
			TypeID:                     typeID,
			Status:                     model.IncidentStatusPending,
			Priority:                   model.IncidentPriorityNormal,
			Description:                desc,
			PatientOriginalDescription: b.Description,
			OccurredAt:                 occurredAtOrNow(b.OccurredAt, now),
			RegistrarEmployeeID:        reg.EmployeeID,
			SourcePatientZitadelUserID: null.StringFrom(b.PatientZitadelUserID),
			SourceBufferID:             uuid.NullUUID{UUID: b.ID, Valid: true},
			CreatedAt:                  now,
			UpdatedAt:                  now,
		}
		if err := tx.Create(&inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}
		if err := projector.IncidentCreated(tx, &inc, &reg); err != nil {
			return err
		}

		// Buffer flips to published, with a back-pointer.
		b.Status = model.BufferStatusPublished
		b.PublishedIncidentID = uuid.NullUUID{UUID: incidentID, Valid: true}
		b.CategoryID = uuid.NullUUID{UUID: categoryID, Valid: true}
		b.TypeID = uuid.NullUUID{UUID: typeID, Valid: true}
		b.UpdatedAt = now
		if err := tx.Save(b).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}
		if err := projector.BufferUpdated(tx, b); err != nil {
			return err
		}
		result.IncidentID = incidentID
		return nil
	})
	return result, err
}

// loadDispatcherSnapshot mirrors loadRegistrarSnapshot from the
// incident package but is local to keep packages decoupled.
func (s *BufferService) loadDispatcherSnapshot(
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
			Code(ErrCodeBufferDispatcherNotFound).
			Public("Dispatcher employee record not found.").Wrap(err)
	}
	if snap.EmployeeID == uuid.Nil {
		return projector.IncidentRegistrarSnapshot{}, oops.In(scope).
			Code(ErrCodeBufferDispatcherNotFound).
			Public("Dispatcher employee record not found.").
			Errorf("employee not found")
	}
	var displayName string
	if err := tx.Raw(
		`SELECT display_name FROM projections.users WHERE id = ?`,
		callerZitadelID,
	).Scan(&displayName).Error; err != nil || displayName == "" {
		return projector.IncidentRegistrarSnapshot{}, oops.In(scope).
			Code(ErrCodeBufferDispatcherNotFound).
			Public("Dispatcher user record missing.").Errorf("display_name lookup failed")
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
