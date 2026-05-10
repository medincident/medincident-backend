package self

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Error codes emitted by GetMyDepartmentRole.
const (
	ErrCodeSelfDeptRoleNotFound = "self_dept_role_not_found"
	ErrCodeSelfDeptRoleFailed   = "self_dept_role_failed"
)

// DepartmentRoleView carries the caller's named role within a department.
type DepartmentRoleView struct {
	IsDepartmentResponsible bool
}

// GetMyDepartmentRole returns whether the caller holds the
// department-responsible role for the given department. Returns
// ErrCodeSelfDeptRoleNotFound when the caller is not an active
// employee of that department.
//
// See: docs/services/self/Self.md
func (r *SelfReader) GetMyDepartmentRole(
	ctx context.Context,
	callerID string,
	departmentID uuid.UUID,
) (*DepartmentRoleView, error) {
	var out DepartmentRoleView
	err := r.db.WithContext(ctx).Raw(`
		SELECT dr.employee_id IS NOT NULL AS is_department_responsible
		  FROM projections.employee_cards ec
		  LEFT JOIN projections.department_responsibles dr
		         ON dr.employee_id = ec.employee_id
		        AND dr.department_id = ?
		 WHERE ec.zitadel_user_id = ?
		   AND ec.department_id = ?
		   AND ec.terminated_at IS NULL`,
		departmentID, callerID, departmentID,
	).Row().Scan(&out.IsDepartmentResponsible)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.self").
				Code(ErrCodeSelfDeptRoleNotFound).
				Public("Employment record not found in this department.").
				With("caller_id", callerID).
				With("department_id", departmentID).
				Errorf("not found")
		}
		return nil, oops.In("reader.self").
			Code(ErrCodeSelfDeptRoleFailed).
			With("caller_id", callerID).
			With("department_id", departmentID).
			Wrap(err)
	}
	return &out, nil
}
