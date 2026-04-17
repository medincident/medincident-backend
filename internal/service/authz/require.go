package authz

import (
	"context"

	"github.com/google/uuid"
)

func (a *Authz) RequireSystemAdmin(ctx context.Context, callerID string) error {
	ok, err := a.checkSystemAdmin(ctx, callerID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID)
	}
	return nil
}

func (a *Authz) RequireOrgAdmin(ctx context.Context, callerID string, orgID uuid.UUID) error {
	ok, err := a.checkOrgAdmin(ctx, callerID, orgID)
	if err != nil {
		return err
	}
	if !ok {
		return errPermissionDenied(callerID)
	}
	return nil
}

func (a *Authz) RequireOrgAdminViaClinic(ctx context.Context, callerID string, clinicID uuid.UUID) error {
	orgID, err := a.resolveOrgByClinic(ctx, clinicID)
	if err != nil {
		return err
	}
	return a.RequireOrgAdmin(ctx, callerID, orgID)
}

func (a *Authz) RequireOrgAdminViaDepartment(ctx context.Context, callerID string, deptID uuid.UUID) error {
	orgID, err := a.resolveOrgByDepartment(ctx, deptID)
	if err != nil {
		return err
	}
	return a.RequireOrgAdmin(ctx, callerID, orgID)
}

func (a *Authz) RequireOrgAdminViaEmployee(ctx context.Context, callerID string, empID uuid.UUID) error {
	orgID, err := a.resolveOrgByEmployee(ctx, empID)
	if err != nil {
		return err
	}
	return a.RequireOrgAdmin(ctx, callerID, orgID)
}

func (a *Authz) RequireOrgAdminViaVacation(ctx context.Context, callerID string, vacID uuid.UUID) error {
	orgID, err := a.resolveOrgByVacation(ctx, vacID)
	if err != nil {
		return err
	}
	return a.RequireOrgAdmin(ctx, callerID, orgID)
}

func (a *Authz) RequireOrgAdminViaCategory(ctx context.Context, callerID string, catID uuid.UUID) error {
	orgID, err := a.resolveOrgByCategory(ctx, catID)
	if err != nil {
		return err
	}
	return a.RequireOrgAdmin(ctx, callerID, orgID)
}

func (a *Authz) RequireOrgAdminViaType(ctx context.Context, callerID string, typeID uuid.UUID) error {
	orgID, err := a.resolveOrgByType(ctx, typeID)
	if err != nil {
		return err
	}
	return a.RequireOrgAdmin(ctx, callerID, orgID)
}
