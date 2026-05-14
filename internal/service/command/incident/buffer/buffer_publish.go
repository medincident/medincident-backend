package buffer

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	incidentv1 "github.com/medincident/medincident-backend/pkg/event/incident/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
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

		// Department -> clinic; clinic must belong to buffer's org.
		var dept model.Department
		if err := tx.First(&dept, "id = ?", deptID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).Code(ErrCodeBufferDeptNotFound).
					Public("Department not found.").
					With("department_id", deptID).Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeBufferLoadFailed).Wrap(err)
		}
		var clinic model.Clinic
		if err := tx.First(&clinic, "id = ?", dept.ClinicID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferDeptNotFound).Wrap(err)
		}
		if clinic.OrganizationID != b.OrganizationID {
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

		// Resolve dispatcher's employee and their department (for clinic denorm).
		dispatcherEmp, err := s.loadDispatcherEmployee(tx, cmd.Caller.ZitadelUserID, b.OrganizationID)
		if err != nil {
			return err
		}
		var dispDept model.Department
		if err := tx.First(&dispDept, "id = ?", dispatcherEmp.DepartmentID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}

		// Description: dispatcher-supplied if present, else patient's text.
		desc := b.Description
		if cmd.Payload.Description != nil {
			desc = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}

		inc := model.Incident{
			ID:                         incidentID,
			OrganizationID:             b.OrganizationID,
			ClinicID:                   clinic.ID,
			DepartmentID:               deptID,
			CategoryID:                 categoryID,
			TypeID:                     typeID,
			Status:                     model.IncidentStatusPending,
			Priority:                   model.IncidentPriorityNormal,
			Description:                desc,
			PatientOriginalDescription: b.Description,
			OccurredAt:                 occurredAtOrNow(b.OccurredAt, now),
			RegistrarEmployeeID:        dispatcherEmp.ID,
			SourcePatientZitadelUserID: null.StringFrom(b.PatientZitadelUserID),
			SourceBufferID:             uuid.NullUUID{UUID: b.ID, Valid: true},
			CreatedAt:                  now,
			UpdatedAt:                  now,
		}
		if err := tx.Create(&inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}

		// Emit IncidentCreated
		incEnv, err := buildBufferIncidentCreatedEnvelope(&inc, cmd.Caller.ZitadelUserID, dispatcherEmp, dispDept.ClinicID)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.incident.v1.created", incEnv); err != nil {
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
		bufEnv, err := buildPatientIncidentBufferUpdatedEnvelope(b)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.patient_incident_buffer.v1.updated", bufEnv); err != nil {
			return err
		}
		result.IncidentID = incidentID
		return nil
	})
	return result, err
}

// loadDispatcherEmployee returns the dispatcher's employee row in the given org.
// Display name is resolved query-side from the event's registrar_zitadel_user_id.
func (s *BufferService) loadDispatcherEmployee(
	tx *gorm.DB, callerZitadelID string, orgID uuid.UUID,
) (*model.Employee, error) {
	var emp model.Employee
	if err := tx.Where("zitadel_user_id = ? AND organization_id = ?",
		callerZitadelID, orgID).First(&emp).Error; err != nil {
		return nil, oops.In(scope).
			Code(ErrCodeBufferDispatcherNotFound).
			Public("Dispatcher employee record not found.").Wrap(err)
	}
	return &emp, nil
}

// buildBufferIncidentCreatedEnvelope builds an IncidentCreated envelope
// for an incident materialised from a buffer entry. Defined locally
// to keep the buffer package independent of the incident package.
func buildBufferIncidentCreatedEnvelope(
	inc *model.Incident,
	registrarZitadelUserID string,
	registrarEmp *model.Employee,
	registrarClinicID uuid.UUID,
) (*eventv1.Envelope, error) {
	msg := &incidentv1.IncidentCreated{
		IncidentId:              inc.ID.String(),
		OrganizationId:          inc.OrganizationID.String(),
		ClinicId:                inc.ClinicID.String(),
		DepartmentId:            inc.DepartmentID.String(),
		CategoryId:              inc.CategoryID.String(),
		TypeId:                  inc.TypeID.String(),
		Status:                  string(inc.Status),
		Priority:                string(inc.Priority),
		OccurredAt:              timestamppb.New(inc.OccurredAt),
		RegistrarZitadelUserId:  registrarZitadelUserID,
		RegistrarEmployeeId:     inc.RegistrarEmployeeID.String(),
		RegistrarOrganizationId: registrarEmp.OrganizationID.String(),
		RegistrarClinicId:       registrarClinicID.String(),
		RegistrarDepartmentId:   registrarEmp.DepartmentID.String(),
		CreatedAt:               timestamppb.New(inc.CreatedAt),
	}
	if registrarEmp.Position.Valid {
		msg.RegistrarPosition = wrapperspb.String(registrarEmp.Position.String)
	}
	if inc.Description.Valid {
		msg.Description = wrapperspb.String(inc.Description.String)
	}
	if inc.PatientOriginalDescription.Valid {
		msg.PatientOriginalDescription = wrapperspb.String(inc.PatientOriginalDescription.String)
	}
	if inc.SourcePatientZitadelUserID.Valid {
		msg.SourcePatientZitadelUserId = wrapperspb.String(inc.SourcePatientZitadelUserID.String)
	}
	if inc.SourceBufferID.Valid {
		msg.SourceBufferId = wrapperspb.String(inc.SourceBufferID.UUID.String())
	}
	if inc.ReopenedFromIncidentID.Valid {
		msg.ReopenedFromIncidentId = wrapperspb.String(inc.ReopenedFromIncidentID.UUID.String())
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(inc.CreatedAt),
		AggregateType: "incident",
		AggregateId:   inc.ID.String(),
		Payload:       payload,
	}, nil
}
