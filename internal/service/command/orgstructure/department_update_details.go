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
)

type UpdateDepartmentDetailsCommand struct {
	ID          uuid.UUID
	Name        string
	Description *string
}

func (s *DepartmentService) UpdateDetails(
	ctx context.Context,
	cmd UpdateDepartmentDetailsCommand,
) error {
	var errs []error
	if err := validateDepartmentName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateDepartmentDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
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
		newDesc := null.StringFromPtr(cmd.Description)
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
