package self

import (
	"context"
	"time"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	memberread "github.com/medincident/medincident-backend/internal/service/query/membership"
	selfread "github.com/medincident/medincident-backend/internal/service/query/self"
	membershipqueryv1 "github.com/medincident/medincident-backend/pkg/query/membership/v1"
	orgqueryv1 "github.com/medincident/medincident-backend/pkg/query/orgstructure/v1"
	selfqueryv1 "github.com/medincident/medincident-backend/pkg/query/self/v1"
)

// SelfQueryHandler implements selfqueryv1.SelfQueryServiceServer.
type SelfQueryHandler struct {
	selfqueryv1.UnimplementedSelfQueryServiceServer
	reader *selfread.SelfReader
}

// NewSelfQueryHandler wires the handler with the reader.
func NewSelfQueryHandler(reader *selfread.SelfReader) *SelfQueryHandler {
	return &SelfQueryHandler{reader: reader}
}

// GetMyIdentity returns whether the caller is a system administrator.
func (h *SelfQueryHandler) GetMyIdentity(
	ctx context.Context,
	_ *selfqueryv1.GetMyIdentityRequest,
) (*selfqueryv1.GetMyIdentityResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	view, err := h.reader.GetMyIdentity(ctx, callerID)
	if err != nil {
		return nil, err
	}
	return &selfqueryv1.GetMyIdentityResponse{IsSystemAdmin: view.IsSystemAdmin}, nil
}

// ListMyOrganizations returns organizations where the caller is an active employee.
func (h *SelfQueryHandler) ListMyOrganizations(
	ctx context.Context,
	_ *selfqueryv1.ListMyOrganizationsRequest,
) (*selfqueryv1.ListMyOrganizationsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListMyOrganizations(ctx, callerID)
	if err != nil {
		return nil, err
	}
	out := make([]*orgqueryv1.OrganizationListItem, 0, len(items))
	for _, item := range items {
		out = append(out, &orgqueryv1.OrganizationListItem{
			Id:   item.ID,
			Name: item.Name,
		})
	}
	return &selfqueryv1.ListMyOrganizationsResponse{Items: out}, nil
}

// GetMyEmployment returns the caller's employee card in the given org.
func (h *SelfQueryHandler) GetMyEmployment(
	ctx context.Context,
	req *selfqueryv1.GetMyEmploymentRequest,
) (*selfqueryv1.GetMyEmploymentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	view, err := h.reader.GetMyEmployment(ctx, callerID, orgID)
	if err != nil {
		return nil, err
	}
	return &selfqueryv1.GetMyEmploymentResponse{Employee: employeeCardToProto(view)}, nil
}

// GetMyOrganizationRole returns the caller's named roles in the given org.
func (h *SelfQueryHandler) GetMyOrganizationRole(
	ctx context.Context,
	req *selfqueryv1.GetMyOrganizationRoleRequest,
) (*selfqueryv1.GetMyOrganizationRoleResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	view, err := h.reader.GetMyOrganizationRole(ctx, callerID, orgID)
	if err != nil {
		return nil, err
	}
	return &selfqueryv1.GetMyOrganizationRoleResponse{
		IsOrgAdmin:      view.IsOrgAdmin,
		IsOrgHead:       view.IsOrgHead,
		IsOrgDispatcher: view.IsOrgDispatcher,
	}, nil
}

// GetMyClinicRole returns the caller's clinic-head status in the given clinic.
func (h *SelfQueryHandler) GetMyClinicRole(
	ctx context.Context,
	req *selfqueryv1.GetMyClinicRoleRequest,
) (*selfqueryv1.GetMyClinicRoleResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	clinicID, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	view, err := h.reader.GetMyClinicRole(ctx, callerID, clinicID)
	if err != nil {
		return nil, err
	}
	return &selfqueryv1.GetMyClinicRoleResponse{IsClinicHead: view.IsClinicHead}, nil
}

// GetMyDepartmentRole returns the caller's dept-responsible status in the given dept.
func (h *SelfQueryHandler) GetMyDepartmentRole(
	ctx context.Context,
	req *selfqueryv1.GetMyDepartmentRoleRequest,
) (*selfqueryv1.GetMyDepartmentRoleResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	deptID, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	view, err := h.reader.GetMyDepartmentRole(ctx, callerID, deptID)
	if err != nil {
		return nil, err
	}
	return &selfqueryv1.GetMyDepartmentRoleResponse{IsDepartmentResponsible: view.IsDepartmentResponsible}, nil
}

// employeeCardToProto adapts a membership.EmployeeCardView to its proto form.
// Mirrors the adapter in internal/handler/query/membership/query.go.
func employeeCardToProto(v *memberread.EmployeeCardView) *membershipqueryv1.EmployeeCardView {
	out := &membershipqueryv1.EmployeeCardView{
		EmployeeId:       v.EmployeeID.String(),
		ZitadelUserId:    v.ZitadelUserID,
		FirstName:        v.FirstName,
		LastName:         v.LastName,
		DisplayName:      v.DisplayName,
		Email:            v.Email,
		OrganizationId:   v.OrganizationID.String(),
		OrganizationName: v.OrganizationName,
		ClinicName:       v.ClinicName,
		DepartmentId:     v.DepartmentID.String(),
		DepartmentName:   v.DepartmentName,
		Position:         v.Position,
	}
	if v.ClinicID != nil {
		s := v.ClinicID.String()
		out.ClinicId = &s
	}
	if v.TerminatedAt != nil {
		s := v.TerminatedAt.UTC().Format(time.RFC3339Nano)
		out.TerminatedAt = &s
	}
	if v.CurrentVacationEndsAt != nil {
		s := v.CurrentVacationEndsAt.UTC().Format(time.RFC3339Nano)
		out.CurrentVacationEndsAt = &s
	}
	if v.NextVacationStartsAt != nil {
		s := v.NextVacationStartsAt.UTC().Format(time.RFC3339Nano)
		out.NextVacationStartsAt = &s
	}
	return out
}
