package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
	clinicv1 "github.com/medincident/medincident-command-service/pkg/event/clinic/v1"
)

// AssignClinicHeadCommand carries the identifiers needed to link an
// employee to a clinic as its head.
type AssignClinicHeadCommand struct {
	ClinicID   uuid.UUID
	EmployeeID uuid.UUID
}

// AssignClinicHead links the employee to the clinic as its head. The
// employee must currently work in a department that belongs to the
// clinic. See spec §8.2 (ClinicHead variant).
func (s *EmployeeService) AssignClinicHead(ctx context.Context, cmd AssignClinicHeadCommand) error {
	var errs []error
	if cmd.ClinicID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).
			Code(ErrCodeClinicIDEmpty).
			Public("Clinic ID is required.").
			Errorf("clinic id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireClinicExists(tx, scopeClinicHead, cmd.ClinicID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeClinicHead, cmd.EmployeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInClinic(tx, scopeClinicHead, emp, cmd.ClinicID); err != nil {
			return err
		}

		row := model.ClinicHead{
			ClinicID:   cmd.ClinicID,
			EmployeeID: cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeClinicHead).
					Code(ErrCodeClinicHeadAlreadyAssigned).
					Public("This employee is already a clinic head.").
					With("clinic_id", cmd.ClinicID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
		}

		ev := &clinicv1.ClinicHeadAssigned{
			EmployeeId: cmd.EmployeeID.String(),
		}
		return outbox.Publish(tx, SubjectClinicHeadAssigned, AggregateTypeClinic, cmd.ClinicID.String(), row.CreatedAt, ev)
	})
}
