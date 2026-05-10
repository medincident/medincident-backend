package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// Error codes emitted by Department-aggregate commands that are not
// primitive validation.
const (
	ErrCodeDepartmentIDGenerationFailed = "department_id_generation_failed"
	ErrCodeDepartmentSaveFailed         = "department_save_failed"
	ErrCodeDepartmentLoadFailed         = "department_load_failed"
	ErrCodeDepartmentNotFound           = "department_not_found"
	ErrCodeDepartmentClinicNotFound     = "department_clinic_not_found"
)

// CreateDepartmentPayload is the validated client-facing payload of
// CreateDepartment.
type CreateDepartmentPayload struct {
	ClinicID    string  `validate:"required,uuid"`
	Name        string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// CreateDepartmentCommand = caller + payload.
type CreateDepartmentCommand struct {
	Caller  authz.Caller
	Payload CreateDepartmentPayload
}

// CreateDepartmentResult is the output of DepartmentService.Create.
type CreateDepartmentResult struct {
	ID uuid.UUID
}

func buildDepartmentCreatedEnvelope(d *model.Department) (*eventv1.Envelope, error) {
	var desc string
	if d.Description.Valid {
		desc = d.Description.String
	}
	msg := &deptv1.DepartmentCreated{
		ClinicId:    d.ClinicID.String(),
		Name:        d.Name,
		Description: desc,
		CreatedAt:   timestamppb.New(d.CreatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(d.CreatedAt),
		AggregateType: "department",
		AggregateId:   d.ID.String(),
		Payload:       payload,
	}, nil
}

// Create persists a new Department under the given clinic.
//
// See: docs/services/OrgStructure.md
func (s *DepartmentService) Create(
	ctx context.Context,
	cmd CreateDepartmentCommand,
) (CreateDepartmentResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateDepartmentResult{}, err
	}
	clinicID := uuid.MustParse(cmd.Payload.ClinicID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return CreateDepartmentResult{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateDepartmentResult{}, oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentIDGenerationFailed).
			Public("Failed to create department.").
			Wrap(err)
	}

	dept := model.Department{
		ID:       id,
		ClinicID: clinicID,
		Name:     strings.TrimSpace(cmd.Payload.Name),
	}
	if cmd.Payload.Description != nil {
		dept.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
	}

	var result CreateDepartmentResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dept).Error; err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", clinicID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentSaveFailed).
				With("department_id", id).
				Wrap(err)
		}

		env, err := buildDepartmentCreatedEnvelope(&dept)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.department.v1.created", env); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
