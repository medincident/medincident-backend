package orgstructure

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/outbox"
	departmentv1 "github.com/medincident/medincident-command-service/pkg/event/department/v1"
)

const (
	departmentMinNameLen = 4
	departmentMaxNameLen = 256
	departmentMinDescLen = 8
	departmentMaxDescLen = 2048
)

const (
	ErrCodeDepartmentNameEmpty           = "department_name_empty"
	ErrCodeDepartmentNameTooShort        = "department_name_too_short"
	ErrCodeDepartmentNameTooLong         = "department_name_too_long"
	ErrCodeDepartmentDescriptionTooShort = "department_description_too_short"
	ErrCodeDepartmentDescriptionTooLong  = "department_description_too_long"

	ErrCodeDepartmentIDGenerationFailed = "department_id_generation_failed"
	ErrCodeDepartmentSaveFailed         = "department_save_failed"
	ErrCodeDepartmentLoadFailed         = "department_load_failed"
	ErrCodeDepartmentNotFound           = "department_not_found"
	ErrCodeDepartmentEventBuildFailed   = "department_event_build_failed"
	ErrCodeDepartmentClinicNotFound     = "department_clinic_not_found"
)

const (
	SubjectDepartmentCreated        = "medincident.event.department.v1.created"
	SubjectDepartmentDetailsChanged = "medincident.event.department.v1.details_changed"
)

type CreateDepartmentCommand struct {
	ClinicID    uuid.UUID
	Name        string
	Description *string
}

type CreateDepartmentResult struct {
	ID uuid.UUID
}

func validateDepartmentName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentNameEmpty).
			Public("Department name is required.").
			With("field", "name").
			Errorf("name is empty")
	}
	n := utf8.RuneCountInString(trimmed)
	if n < departmentMinNameLen {
		return oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentNameTooShort).
			Public("Department name is too short.").
			With("field", "name").
			With("actual_length", n).
			With("min_length", departmentMinNameLen).
			Errorf("name too short")
	}
	if n > departmentMaxNameLen {
		return oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentNameTooLong).
			Public("Department name is too long.").
			With("field", "name").
			With("actual_length", n).
			With("max_length", departmentMaxNameLen).
			Errorf("name too long")
	}
	return nil
}

func validateDepartmentDescription(desc *string) error {
	if desc == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*desc)
	n := utf8.RuneCountInString(trimmed)
	if n < departmentMinDescLen {
		return oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentDescriptionTooShort).
			Public("Department description is too short.").
			With("field", "description").
			With("actual_length", n).
			With("min_length", departmentMinDescLen).
			Errorf("description too short")
	}
	if n > departmentMaxDescLen {
		return oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentDescriptionTooLong).
			Public("Department description is too long.").
			With("field", "description").
			With("actual_length", n).
			With("max_length", departmentMaxDescLen).
			Errorf("description too long")
	}
	return nil
}

func buildDepartmentCreatedEvent(d *model.Department) *departmentv1.DepartmentCreated {
	return &departmentv1.DepartmentCreated{
		ClinicId:    d.ClinicID.String(),
		Name:        d.Name,
		Description: d.Description.Ptr(),
	}
}

func (s *DepartmentService) Create(
	ctx context.Context,
	cmd CreateDepartmentCommand,
) (CreateDepartmentResult, error) {
	var errs []error
	if err := validateDepartmentName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateDepartmentDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return CreateDepartmentResult{}, errors.Join(errs...)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateDepartmentResult{}, oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentIDGenerationFailed).
			Public("Failed to create department.").
			Wrap(err)
	}

	dept := model.Department{
		ID:          id,
		ClinicID:    cmd.ClinicID,
		Name:        strings.TrimSpace(cmd.Name),
		Description: null.StringFromPtr(cmd.Description),
	}

	var result CreateDepartmentResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dept).Error; err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", cmd.ClinicID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentSaveFailed).
				With("department_id", id).
				Wrap(err)
		}

		event := buildDepartmentCreatedEvent(&dept)
		if err := outbox.Publish(tx, SubjectDepartmentCreated, AggregateTypeDepartment, dept.ID.String(), dept.UpdatedAt, event); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
