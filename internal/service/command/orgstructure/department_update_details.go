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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// UpdateDepartmentDetailsPayload carries the new name and (optional)
// description for an existing department.
type UpdateDepartmentDetailsPayload struct {
	ID          string  `validate:"required,uuid"`
	Name        string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// UpdateDepartmentDetailsCommand = caller + payload.
type UpdateDepartmentDetailsCommand struct {
	Caller  authz.Caller
	Payload UpdateDepartmentDetailsPayload
}

// UpdateDetails changes a department's name and description.
//
// See: docs/services/OrgStructure.md
func (s *DepartmentService) UpdateDetails(
	ctx context.Context,
	cmd UpdateDepartmentDetailsCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(id)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dept model.Department
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&dept, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", id).
					Errorf("department not found")
			}
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("department_id", id).
				Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Payload.Name)
		var newDesc null.String
		if cmd.Payload.Description != nil {
			newDesc = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if dept.Name == newName && dept.Description == newDesc {
			return nil
		}
		dept.Name = newName
		dept.Description = newDesc
		if err := tx.Save(&dept).Error; err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentSaveFailed).
				With("department_id", id).
				Wrap(err)
		}

		return projector.DepartmentDetailsChanged(tx, &dept)
	})
}
