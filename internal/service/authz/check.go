package authz

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

func (a *Authz) checkSystemAdmin(ctx context.Context, callerID string) (bool, error) {
	var exists bool
	if err := a.db.WithContext(ctx).
		Raw(`SELECT EXISTS(SELECT 1 FROM domain.system_admins WHERE zitadel_user_id = ?)`, callerID).
		Scan(&exists).Error; err != nil {
		return false, oops.In("service.authz").Code(ErrCodeAuthzQueryFailed).
			With("caller_id", callerID).Wrap(err)
	}
	return exists, nil
}

func (a *Authz) checkOrgAdmin(ctx context.Context, callerID string, orgID uuid.UUID) (bool, error) {
	// 1. System admin has access to everything.
	isSysAdmin, err := a.checkSystemAdmin(ctx, callerID)
	if err != nil {
		return false, err
	}
	if isSysAdmin {
		return true, nil
	}

	// 2. Find employee in this organization.
	var employeeIDStr *string
	if err := a.db.WithContext(ctx).
		Raw(`SELECT id::text FROM domain.employees WHERE zitadel_user_id = ? AND organization_id = ?`, callerID, orgID).
		Scan(&employeeIDStr).Error; err != nil {
		return false, oops.In("service.authz").Code(ErrCodeAuthzQueryFailed).
			With("caller_id", callerID).With("organization_id", orgID).Wrap(err)
	}
	if employeeIDStr == nil {
		return false, nil
	}
	employeeID, err := uuid.Parse(*employeeIDStr)
	if err != nil {
		return false, oops.In("service.authz").Code(ErrCodeAuthzQueryFailed).Wrap(err)
	}

	// 3. Direct org admin?
	var isAdmin bool
	if err := a.db.WithContext(ctx).
		Raw(`SELECT EXISTS(SELECT 1 FROM domain.org_admins WHERE organization_id = ? AND employee_id = ?)`, orgID, employeeID).
		Scan(&isAdmin).Error; err != nil {
		return false, oops.In("service.authz").Code(ErrCodeAuthzQueryFailed).
			With("caller_id", callerID).With("organization_id", orgID).Wrap(err)
	}
	if isAdmin {
		return true, nil
	}

	// 4. Deputy of an org admin who is on active vacation?
	var isDeputy bool
	if err := a.db.WithContext(ctx).
		Raw(`SELECT EXISTS(
			SELECT 1 FROM domain.org_admins oa
			JOIN domain.employee_vacations v ON v.employee_id = oa.employee_id
			WHERE oa.organization_id = ?
			  AND oa.deputy_employee_id = ?
			  AND v.starts_at <= now()
			  AND (v.ends_at IS NULL OR v.ends_at > now())
		)`, orgID, employeeID).
		Scan(&isDeputy).Error; err != nil {
		return false, oops.In("service.authz").Code(ErrCodeAuthzQueryFailed).
			With("caller_id", callerID).With("organization_id", orgID).Wrap(err)
	}
	return isDeputy, nil
}
