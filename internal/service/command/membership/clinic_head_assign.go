package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// AssignClinicHeadPayload carries the identifiers needed to link an
// employee to a clinic as its head.
type AssignClinicHeadPayload struct {
	ClinicID   string `validate:"required,uuid"`
	EmployeeID string `validate:"required,uuid"`
}

// AssignClinicHeadCommand = caller + payload.
type AssignClinicHeadCommand struct {
	Caller  authz.Caller
	Payload AssignClinicHeadPayload
}

// AssignClinicHead links the employee to the clinic as its head. The
// employee must currently work in a department that belongs to the
// clinic. See spec §8.2 (ClinicHead variant).
func (s *EmployeeService) AssignClinicHead(ctx context.Context, cmd AssignClinicHeadCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	clinicID := uuid.MustParse(cmd.Payload.ClinicID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireClinicExists(tx, scopeClinicHead, clinicID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeClinicHead, employeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInClinic(tx, scopeClinicHead, emp, clinicID); err != nil {
			return err
		}

		row := model.ClinicHead{
			ClinicID:   clinicID,
			EmployeeID: employeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeClinicHead).
					Code(ErrCodeClinicHeadAlreadyAssigned).
					Public("This employee is already a clinic head.").
					With("clinic_id", clinicID).
					With("employee_id", employeeID).
					Wrap(err)
			}
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
		}

		return projector.ClinicHeadAssigned(tx, &row)
	})
}
