package self

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Error codes emitted by GetMyOrganizationRole.
const (
	ErrCodeSelfOrgRoleNotFound = "self_org_role_not_found"
	ErrCodeSelfOrgRoleFailed   = "self_org_role_failed"
)

// OrgRoleView carries the caller's named roles within an organization.
type OrgRoleView struct {
	IsOrgAdmin      bool
	IsOrgHead       bool
	IsOrgDispatcher bool
}

// GetMyOrganizationRole returns the caller's named roles in the given
// organization. Returns ErrCodeSelfOrgRoleNotFound when the caller is
// not an active employee of that organization.
//
// See: docs/services/self/Self.md
func (r *SelfReader) GetMyOrganizationRole(
	ctx context.Context,
	callerID string,
	organizationID uuid.UUID,
) (*OrgRoleView, error) {
	var out OrgRoleView
	err := r.db.WithContext(ctx).Raw(`
		SELECT oa.employee_id IS NOT NULL AS is_org_admin,
		       oh.employee_id IS NOT NULL AS is_org_head,
		       od.employee_id IS NOT NULL AS is_org_dispatcher
		  FROM projections.employee_cards ec
		  LEFT JOIN projections.org_admins oa ON oa.employee_id = ec.employee_id
		  LEFT JOIN projections.org_heads  oh ON oh.employee_id = ec.employee_id
		  LEFT JOIN projections.org_dispatchers od ON od.employee_id = ec.employee_id
		 WHERE ec.zitadel_user_id = ?
		   AND ec.organization_id = ?
		   AND ec.terminated_at IS NULL`,
		callerID, organizationID,
	).Row().Scan(&out.IsOrgAdmin, &out.IsOrgHead, &out.IsOrgDispatcher)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.self").
				Code(ErrCodeSelfOrgRoleNotFound).
				Public("Employment record not found in this organization.").
				With("caller_id", callerID).
				With("organization_id", organizationID).
				Errorf("not found")
		}
		return nil, oops.In("reader.self").
			Code(ErrCodeSelfOrgRoleFailed).
			With("caller_id", callerID).
			With("organization_id", organizationID).
			Wrap(err)
	}
	return &out, nil
}
