package membership

import (
	"context"
	"strings"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	memberread "github.com/medincident/medincident-backend/internal/service/query/membership"
	membershipqueryv1 "github.com/medincident/medincident-backend/pkg/query/membership/v1"
)

// ListCandidatesForHire returns Zitadel users not yet active employees
// of the given organization.
//
// See: docs/services/Membership.md
func (h *MembershipQueryHandler) ListCandidatesForHire(
	ctx context.Context,
	req *membershipqueryv1.ListCandidatesForHireRequest,
) (*membershipqueryv1.ListCandidatesForHireResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	items, next, err := h.candidateReader.ForHire(
		ctx,
		caller,
		orgID,
		strings.TrimSpace(req.GetQuery()),
		req.GetAfter(),
		int(req.GetLimit()),
	)
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListCandidatesForHireResponse{
		Items:      zitadelUsersToProto(items),
		NextCursor: next,
	}, nil
}

// ListCandidatesForSystemAdmin returns Zitadel users not yet system admins.
//
// See: docs/services/Membership.md
func (h *MembershipQueryHandler) ListCandidatesForSystemAdmin(
	ctx context.Context,
	req *membershipqueryv1.ListCandidatesForSystemAdminRequest,
) (*membershipqueryv1.ListCandidatesForSystemAdminResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	items, next, err := h.candidateReader.ForSystemAdmin(
		ctx,
		caller,
		strings.TrimSpace(req.GetQuery()),
		req.GetAfter(),
		int(req.GetLimit()),
	)
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListCandidatesForSystemAdminResponse{
		Items:      zitadelUsersToProto(items),
		NextCursor: next,
	}, nil
}

// ListCandidatesForOrgAdmin returns active org employees not yet assigned
// as org admin.
//
// See: docs/services/Membership.md
func (h *MembershipQueryHandler) ListCandidatesForOrgAdmin(
	ctx context.Context,
	req *membershipqueryv1.ListCandidatesForOrgAdminRequest,
) (*membershipqueryv1.ListCandidatesForOrgAdminResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	items, next, err := h.candidateReader.ForOrgAdmin(
		ctx,
		caller,
		orgID,
		strings.TrimSpace(req.GetQuery()),
		req.GetAfter(),
		int(req.GetLimit()),
	)
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListCandidatesForOrgAdminResponse{
		Items:      employeeCardsToProto(items),
		NextCursor: next,
	}, nil
}

// ListCandidatesForOrgHead returns active org employees not yet assigned
// as org head.
//
// See: docs/services/Membership.md
func (h *MembershipQueryHandler) ListCandidatesForOrgHead(
	ctx context.Context,
	req *membershipqueryv1.ListCandidatesForOrgHeadRequest,
) (*membershipqueryv1.ListCandidatesForOrgHeadResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	items, next, err := h.candidateReader.ForOrgHead(
		ctx,
		caller,
		orgID,
		strings.TrimSpace(req.GetQuery()),
		req.GetAfter(),
		int(req.GetLimit()),
	)
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListCandidatesForOrgHeadResponse{
		Items:      employeeCardsToProto(items),
		NextCursor: next,
	}, nil
}

// ListCandidatesForOrgDispatcher returns active org employees not yet
// assigned as org dispatcher.
//
// See: docs/services/Membership.md
func (h *MembershipQueryHandler) ListCandidatesForOrgDispatcher(
	ctx context.Context,
	req *membershipqueryv1.ListCandidatesForOrgDispatcherRequest,
) (*membershipqueryv1.ListCandidatesForOrgDispatcherResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	items, next, err := h.candidateReader.ForOrgDispatcher(
		ctx,
		caller,
		orgID,
		strings.TrimSpace(req.GetQuery()),
		req.GetAfter(),
		int(req.GetLimit()),
	)
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListCandidatesForOrgDispatcherResponse{
		Items:      employeeCardsToProto(items),
		NextCursor: next,
	}, nil
}

// ListCandidatesForClinicHead returns active clinic employees not yet
// assigned as clinic head.
//
// See: docs/services/Membership.md
func (h *MembershipQueryHandler) ListCandidatesForClinicHead(
	ctx context.Context,
	req *membershipqueryv1.ListCandidatesForClinicHeadRequest,
) (*membershipqueryv1.ListCandidatesForClinicHeadResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	clinicID, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	items, next, err := h.candidateReader.ForClinicHead(
		ctx,
		caller,
		clinicID,
		strings.TrimSpace(req.GetQuery()),
		req.GetAfter(),
		int(req.GetLimit()),
	)
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListCandidatesForClinicHeadResponse{
		Items:      employeeCardsToProto(items),
		NextCursor: next,
	}, nil
}

// ListCandidatesForDeptResponsible returns active department employees
// not yet assigned as department responsible.
//
// See: docs/services/Membership.md
func (h *MembershipQueryHandler) ListCandidatesForDeptResponsible(
	ctx context.Context,
	req *membershipqueryv1.ListCandidatesForDeptResponsibleRequest,
) (*membershipqueryv1.ListCandidatesForDeptResponsibleResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	deptID, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	items, next, err := h.candidateReader.ForDeptResponsible(
		ctx,
		caller,
		deptID,
		strings.TrimSpace(req.GetQuery()),
		req.GetAfter(),
		int(req.GetLimit()),
	)
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListCandidatesForDeptResponsibleResponse{
		Items:      employeeCardsToProto(items),
		NextCursor: next,
	}, nil
}

// zitadelUsersToProto converts a slice of ZitadelUserView to proto messages.
func zitadelUsersToProto(views []memberread.ZitadelUserView) []*membershipqueryv1.ZitadelUserView {
	out := make([]*membershipqueryv1.ZitadelUserView, 0, len(views))
	for _, v := range views {
		out = append(out, &membershipqueryv1.ZitadelUserView{
			ZitadelUserId: v.ZitadelUserID,
			FirstName:     v.FirstName,
			LastName:      v.LastName,
			DisplayName:   v.DisplayName,
			Email:         v.Email,
		})
	}
	return out
}
