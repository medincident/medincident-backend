package incident

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

// CreateIncidentPayload is the validated client-facing payload.
type CreateIncidentPayload struct {
	DepartmentID string  `validate:"required,uuid"`
	CategoryID   string  `validate:"required,uuid"`
	TypeID       string  `validate:"required,uuid"`
	Description  *string `validate:"omitnil,no_extra_ws,min=1,max=10000"`
	OccurredAt   string  `validate:"required"`
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
	occurred, err := time.Parse(time.RFC3339Nano, cmd.Payload.OccurredAt)
	if err != nil {
		return CreateIncidentResult{}, oops.In(scope).
			Code(ErrCodeIncidentOccurredAtInvalid).
			Public("occurred_at is not a valid RFC3339 timestamp.").
			With("occurred_at", cmd.Payload.OccurredAt).Wrap(err)
	}

	now := time.Now()
	if occurred.After(now) {
		return CreateIncidentResult{}, oops.In(scope).
			Code(ErrCodeIncidentOccurredAtFuture).
			Public("occurred_at must not be in the future.").
			With("occurred_at", occurred).
			Errorf("future occurred_at")
	}
	if now.Sub(occurred) > incidentMaxOccurredAtAge {
		return CreateIncidentResult{}, oops.In(scope).
			Code(ErrCodeIncidentOccurredAtTooOld).
			Public("occurred_at exceeds the 48h registration window.").
			With("occurred_at", occurred).
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
		// Resolve dept -> clinic -> organization.
		var dept model.Department
		if err := tx.First(&dept, "id = ?", deptID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).
					Code(ErrCodeIncidentDeptNotFound).
					Public("Department not found.").
					With("department_id", deptID).
					Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeIncidentLoadFailed).Wrap(err)
		}
		var clinic model.Clinic
		if err := tx.First(&clinic, "id = ?", dept.ClinicID).Error; err != nil {
			return oops.In(scope).
				Code(ErrCodeIncidentClinicNotFound).
				With("clinic_id", dept.ClinicID).
				Wrap(err)
		}
		orgID := clinic.OrganizationID

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

		// Resolve registrar employee and their department (for clinic denorm).
		registrarEmp, err := s.loadRegistrarEmployee(tx, cmd.Caller.ZitadelUserID, orgID)
		if err != nil {
			return err
		}
		var regDept model.Department
		if err := tx.First(&regDept, "id = ?", registrarEmp.DepartmentID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentLoadFailed).Wrap(err)
		}

		inc := model.Incident{
			ID:                  id,
			OrganizationID:      orgID,
			ClinicID:            clinic.ID,
			DepartmentID:        deptID,
			CategoryID:          categoryID,
			TypeID:              typeID,
			Status:              model.IncidentStatusPending,
			Priority:            model.IncidentPriorityNormal,
			OccurredAt:          occurred,
			RegistrarEmployeeID: registrarEmp.ID,
			CreatedAt:           now,
			UpdatedAt:           now,
		}
		if cmd.Payload.Description != nil {
			inc.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if err := tx.Create(&inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).
				With("incident_id", id).Wrap(err)
		}
		env, err := buildIncidentCreatedEnvelope(&inc, cmd.Caller.ZitadelUserID, registrarEmp, regDept.ClinicID)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.incident.v1.created", env); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}

// loadRegistrarEmployee loads the caller's employee row in the given org.
// The event carries the Zitadel user ID; display_name is resolved query-side.
func (s *IncidentService) loadRegistrarEmployee(
	tx *gorm.DB, callerZitadelID string, orgID uuid.UUID,
) (*model.Employee, error) {
	var emp model.Employee
	if err := tx.Where("zitadel_user_id = ? AND organization_id = ?",
		callerZitadelID, orgID).First(&emp).Error; err != nil {
		return nil, oops.In(scope).
			Code(ErrCodeIncidentEmployeeNotFound).
			Public("Caller is not an employee of this organization.").
			With("organization_id", orgID).
			Wrap(err)
	}
	return &emp, nil
}

func buildIncidentCreatedEnvelope(
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
		return nil, oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(inc.CreatedAt),
		AggregateType: "incident",
		AggregateId:   inc.ID.String(),
		Payload:       payload,
	}, nil
}
