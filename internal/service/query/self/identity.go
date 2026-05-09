package self

import (
	"context"

	"github.com/samber/oops"
)

// Error codes emitted by SelfReader identity and organization list methods.
const (
	ErrCodeSelfIdentityFailed     = "self_identity_failed"
	ErrCodeSelfOrgsFailed         = "self_orgs_failed"
	ErrCodeSelfClinicRoleNotFound = "self_clinic_role_not_found"
	ErrCodeSelfClinicRoleFailed   = "self_clinic_role_failed"
	ErrCodeSelfDeptRoleNotFound   = "self_dept_role_not_found"
	ErrCodeSelfDeptRoleFailed     = "self_dept_role_failed"
)

// IdentityView carries the global identity facts about the caller.
type IdentityView struct {
	IsSystemAdmin bool
}

// GetMyIdentity returns whether the caller is a system administrator.
// It never returns NOT_FOUND — a non-admin caller gets IsSystemAdmin=false.
//
// See: docs/services/self/Self.md
func (r *SelfReader) GetMyIdentity(ctx context.Context, callerID string) (*IdentityView, error) {
	var exists bool
	err := r.db.WithContext(ctx).Raw(
		`SELECT EXISTS(SELECT 1 FROM projections.system_admins WHERE zitadel_user_id = ?)`,
		callerID,
	).Scan(&exists).Error
	if err != nil {
		return nil, oops.In("reader.self").
			Code(ErrCodeSelfIdentityFailed).
			With("caller_id", callerID).
			Wrap(err)
	}
	return &IdentityView{IsSystemAdmin: exists}, nil
}

// OrganizationListItem is a minimal organization record for list display.
type OrganizationListItem struct {
	ID   string
	Name string
}

// ListMyOrganizations returns every organization where the caller has
// an active (non-terminated) employee record, ordered by name.
// Returns an empty slice when the caller is in no organizations.
//
// See: docs/services/self/Self.md
func (r *SelfReader) ListMyOrganizations(ctx context.Context, callerID string) ([]OrganizationListItem, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT o.id, o.name
		  FROM projections.organizations o
		  JOIN projections.employee_cards ec ON ec.organization_id = o.id
		 WHERE ec.zitadel_user_id = ?
		   AND ec.terminated_at IS NULL
		 ORDER BY o.name ASC`, callerID).Rows()
	if err != nil {
		return nil, oops.In("reader.self").
			Code(ErrCodeSelfOrgsFailed).
			With("caller_id", callerID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	var out []OrganizationListItem
	for rows.Next() {
		var item OrganizationListItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, oops.In("reader.self").
				Code(ErrCodeSelfOrgsFailed).
				With("caller_id", callerID).
				Wrap(err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.self").
			Code(ErrCodeSelfOrgsFailed).
			With("caller_id", callerID).
			Wrap(err)
	}
	return out, nil
}
