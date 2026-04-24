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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
	"github.com/medincident/medincident-backend/internal/service/zitadel"
)

// HireEmployeePayload carries everything the service needs to create
// a new Employee row. The organisation is derived from the
// department's parent lineage; callers pass only the department.
type HireEmployeePayload struct {
	ZitadelUserID string  `validate:"required,no_extra_ws"`
	DepartmentID  string  `validate:"required,uuid"`
	Position      *string `validate:"omitnil,no_extra_ws,min=2,max=256"`
}

// HireEmployeeCommand = caller + payload.
type HireEmployeeCommand struct {
	Caller  authz.Caller
	Payload HireEmployeePayload
}

// HireEmployeeResult holds the identifiers of the newly created employee.
type HireEmployeeResult struct {
	ID uuid.UUID
}

// Hire creates a new Employee row in the given department, deriving
// the organisation via a JOIN on domain.departments → domain.clinics.
// Invariants and error codes are spelled out in the spec.
func (s *EmployeeService) Hire(ctx context.Context, cmd HireEmployeeCommand) (HireEmployeeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return HireEmployeeResult{}, err
	}
	departmentID := uuid.MustParse(cmd.Payload.DepartmentID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(departmentID)); err != nil {
		return HireEmployeeResult{}, err
	}

	zitadelUserID := strings.TrimSpace(cmd.Payload.ZitadelUserID)
	var position null.String
	if cmd.Payload.Position != nil {
		position = null.StringFrom(strings.TrimSpace(*cmd.Payload.Position))
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
			WHERE d.id = ?`, departmentID).Row().Scan(&orgID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
				return oops.In(scopeEmployee).
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", departmentID).
					Errorf("department not found")
			}
			return oops.In(scopeEmployee).
				Code(ErrCodeDepartmentLookupFailed).
				With("department_id", departmentID).
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
			DepartmentID:   departmentID,
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

		if err := projector.EmployeeHired(tx, &emp); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
