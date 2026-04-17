package membership

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
	"github.com/medincident/medincident-command-service/internal/service/zitadel"
)

// HireEmployeeCommand carries everything the service needs to create
// a new Employee row. The organisation is derived from the
// department's parent lineage; callers pass only the department.
type HireEmployeeCommand struct {
	ZitadelUserID string
	DepartmentID  uuid.UUID
	Position      *string
}

// HireEmployeeResult holds the identifiers of the newly created employee.
type HireEmployeeResult struct {
	ID uuid.UUID
}

// Hire creates a new Employee row in the given department, deriving
// the organisation via a JOIN on domain.departments → domain.clinics.
// Invariants and error codes are spelled out in the spec.
func (s *EmployeeService) Hire(ctx context.Context, cmd HireEmployeeCommand) (HireEmployeeResult, error) {
	// Phase 1: multi-error local validation.
	var errs []error
	if strings.TrimSpace(cmd.ZitadelUserID) == "" {
		errs = append(errs, oops.In(scopeEmployee).
			Code(ErrCodeEmployeeZitadelUserIDEmpty).
			Public("Zitadel user ID is required.").
			Errorf("zitadel user id is empty"))
	}
	if cmd.DepartmentID == uuid.Nil {
		errs = append(errs, oops.In(scopeEmployee).
			Code(ErrCodeEmployeeDepartmentIDEmpty).
			Public("Department ID is required.").
			Errorf("department id is empty"))
	}
	if err := validatePosition(cmd.Position); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return HireEmployeeResult{}, errors.Join(errs...)
	}

	zitadelUserID := strings.TrimSpace(cmd.ZitadelUserID)
	var position null.String
	if cmd.Position != nil {
		position = null.StringFrom(strings.TrimSpace(*cmd.Position))
	}

	// Phase 2: Zitadel verify (outside tx, fail-fast).
	if err := s.verifier.Verify(ctx, zitadelUserID); err != nil {
		if errors.Is(err, zitadel.ErrUserNotFound) {
			return HireEmployeeResult{}, oops.In(scopeEmployee).
				Code(ErrCodeZitadelUserNotFound).
				Public("User does not exist in identity provider.").
				With("zitadel_user_id", zitadelUserID).
				Wrap(err)
		}
		return HireEmployeeResult{}, oops.In(scopeEmployee).
			Code(zitadel.ErrCodeZitadelVerifyFailed).
			With("zitadel_user_id", zitadelUserID).
			Wrap(err)
	}

	// Phase 3: transaction.
	var result HireEmployeeResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var orgID uuid.UUID
		err := tx.Raw(`
			SELECT c.organization_id
			FROM domain.departments d
			JOIN domain.clinics c ON c.id = d.clinic_id
			WHERE d.id = ?`, cmd.DepartmentID).Row().Scan(&orgID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
				return oops.In(scopeEmployee).
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", cmd.DepartmentID).
					Errorf("department not found")
			}
			return oops.In(scopeEmployee).
				Code(ErrCodeDepartmentLookupFailed).
				With("department_id", cmd.DepartmentID).
				Wrap(err)
		}

		id, err := uuid.NewV7()
		if err != nil {
			return oops.In(scopeEmployee).
				Code(ErrCodeEmployeeIDGenerationFailed).
				Wrap(err)
		}

		emp := model.Employee{
			ID:             id,
			ZitadelUserID:  zitadelUserID,
			OrganizationID: orgID,
			DepartmentID:   cmd.DepartmentID,
			Position:       position,
		}
		if err := tx.Create(&emp).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeEmployee).
					Code(ErrCodeEmployeeAlreadyHired).
					Public("Employee is already hired in this organization.").
					With("organization_id", orgID).
					With("zitadel_user_id", zitadelUserID).
					Wrap(err)
			}
			return oops.In(scopeEmployee).
				Code(ErrCodeEmployeeSaveFailed).
				Wrap(err)
		}

		ev := &employeev1.EmployeeHired{
			ZitadelUserId:  zitadelUserID,
			OrganizationId: orgID.String(),
			DepartmentId:   cmd.DepartmentID.String(),
			Position:       position.Ptr(),
		}
		if err := outbox.Publish(tx, SubjectEmployeeHired, AggregateTypeEmployee, id.String(), emp.UpdatedAt, ev); err != nil {
			return err
		}

		result.ID = id
		return nil
	})
	return result, err
}
