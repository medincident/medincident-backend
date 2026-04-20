package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// UpdateDepartmentDetailsCommand carries the new name and (optional)
// description for an existing department.
type UpdateDepartmentDetailsCommand struct {
	ID          uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=4,max=256"`
	Description *string   `validate:"omitnil,min=8,max=2048"`
}

// UpdateDetails changes a department's name and description.
func (s *DepartmentService) UpdateDetails(
	ctx context.Context,
	cmd UpdateDepartmentDetailsCommand,
) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dept model.Department
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&dept, "id = ?", cmd.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", cmd.ID).
					Errorf("department not found")
			}
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("department_id", cmd.ID).
				Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Name)
		var newDesc null.String
		if cmd.Description != nil {
			newDesc = null.StringFrom(strings.TrimSpace(*cmd.Description))
		}
		if dept.Name == newName && dept.Description == newDesc {
			return nil
		}
		dept.Name = newName
		dept.Description = newDesc
		if err := tx.Save(&dept).Error; err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentSaveFailed).
				With("department_id", cmd.ID).
				Wrap(err)
		}

		return projector.DepartmentDetailsChanged(tx, &dept)
	})
}
