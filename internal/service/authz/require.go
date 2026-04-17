package authz

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// ErrCodePermissionDenied is emitted whenever a check* call returns
// (ok=false, nil): the caller is known, the scope exists or does not,
// and either way they are not authorized. Maps to
// codes.PermissionDenied at the gRPC boundary.
const ErrCodePermissionDenied = "permission_denied"

// errPermissionDenied shapes the deny result with the scope that was
// being checked so server logs carry enough to triage failures. The
// Public message is deliberately scope-agnostic to avoid revealing
// which field mattered — the caller is always "permission denied".
// scopeField/scopeID are zero-valued for RequireSystemAdmin, in which
// case only caller_id is attached.
func errPermissionDenied(callerID, scopeField string, scopeID uuid.UUID) error {
	b := oops.In("service.authz").
		Code(ErrCodePermissionDenied).
		Public("Permission denied.").
		With("caller_id", callerID)
	if scopeField != "" {
		b = b.With(scopeField, scopeID)
	}
	return b.Errorf("permission denied")
}

func (a *Authz) RequireSystemAdmin(ctx context.Context, callerID string) error {
	ok, err := a.checkSystemAdmin(ctx, callerID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID, "", uuid.Nil)
	}
	return nil
}

func (a *Authz) RequireOrgAdmin(ctx context.Context, callerID string, orgID uuid.UUID) error {
	ok, err := a.checkOrgAccess(ctx, callerID, orgID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID, "organization_id", orgID)
	}
	return nil
}

func (a *Authz) RequireOrgAdminViaClinic(ctx context.Context, callerID string, clinicID uuid.UUID) error {
	ok, err := a.checkOrgAccessViaClinic(ctx, callerID, clinicID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID, "clinic_id", clinicID)
	}
	return nil
}

func (a *Authz) RequireOrgAdminViaDepartment(ctx context.Context, callerID string, deptID uuid.UUID) error {
	ok, err := a.checkOrgAccessViaDepartment(ctx, callerID, deptID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID, "department_id", deptID)
	}
	return nil
}

func (a *Authz) RequireOrgAdminViaEmployee(ctx context.Context, callerID string, empID uuid.UUID) error {
	ok, err := a.checkOrgAccessViaEmployee(ctx, callerID, empID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID, "employee_id", empID)
	}
	return nil
}

func (a *Authz) RequireOrgAdminViaVacation(ctx context.Context, callerID string, vacID uuid.UUID) error {
	ok, err := a.checkOrgAccessViaVacation(ctx, callerID, vacID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID, "vacation_id", vacID)
	}
	return nil
}

func (a *Authz) RequireOrgAdminViaCategory(ctx context.Context, callerID string, catID uuid.UUID) error {
	ok, err := a.checkOrgAccessViaCategory(ctx, callerID, catID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID, "category_id", catID)
	}
	return nil
}

func (a *Authz) RequireOrgAdminViaType(ctx context.Context, callerID string, typeID uuid.UUID) error {
	ok, err := a.checkOrgAccessViaType(ctx, callerID, typeID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID, "type_id", typeID)
	}
	return nil
}
