package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
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

// CreateDepartmentCommand is the input of DepartmentService.Create.
type CreateDepartmentCommand struct {
	ClinicID    uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=4,max=256"`
	Description *string   `validate:"omitnil,min=8,max=2048"`
}

// CreateDepartmentResult is the output of DepartmentService.Create.
type CreateDepartmentResult struct {
	ID uuid.UUID
}

// Create persists a new Department under the given clinic.
func (s *DepartmentService) Create(
	ctx context.Context,
	cmd CreateDepartmentCommand,
) (CreateDepartmentResult, error) {
	if err := validation.Struct(cmd); err != nil {
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
		ClinicID: cmd.ClinicID,
		Name:     strings.TrimSpace(cmd.Name),
	}
	if cmd.Description != nil {
		dept.Description = null.StringFrom(strings.TrimSpace(*cmd.Description))
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

		if err := projector.DepartmentCreated(tx, &dept); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
