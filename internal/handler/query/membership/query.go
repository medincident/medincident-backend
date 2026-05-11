package membership

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	memberread "github.com/medincident/medincident-backend/internal/service/query/membership"
	membershipqueryv1 "github.com/medincident/medincident-backend/pkg/query/membership/v1"
)

// MembershipQueryHandler implements
// membershipqueryv1.MembershipQueryServiceServer.
type MembershipQueryHandler struct {
	membershipqueryv1.UnimplementedMembershipQueryServiceServer

	empReader  *memberread.EmployeeReader
	roleReader *memberread.RoleReader
}

// NewMembershipQueryHandler wires the handler with the readers.
func NewMembershipQueryHandler(
	empReader *memberread.EmployeeReader,
	roleReader *memberread.RoleReader,
) *MembershipQueryHandler {
	return &MembershipQueryHandler{
		empReader:  empReader,
		roleReader: roleReader,
	}
}

// afterPtr converts an empty proto string to nil, treating empty as
// "start from the beginning".
func afterPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// GetEmployee returns the denormalised employee card.
func (h *MembershipQueryHandler) GetEmployee(
	ctx context.Context,
	req *membershipqueryv1.GetEmployeeRequest,
) (*membershipqueryv1.GetEmployeeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseEmployeeID(req.GetId())
	if err != nil {
		return nil, err
	}
	view, err := h.empReader.Get(ctx, caller, id)
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.GetEmployeeResponse{Employee: employeeCardToProto(view)}, nil
}

// ListEmployeesByDepartment returns the cards under a department.
func (h *MembershipQueryHandler) ListEmployeesByDepartment(
	ctx context.Context,
	req *membershipqueryv1.ListEmployeesByDepartmentRequest,
) (*membershipqueryv1.ListEmployeesByDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	result, err := h.empReader.ListByDepartment(ctx, caller, id, memberread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	}, memberread.EmployeeFilter{
		IncludeTerminated: req.GetIncludeTerminated(),
		OnVacation:        req.GetOnVacation(),
		Position:          strings.TrimSpace(req.GetPosition()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListEmployeesByDepartmentResponse{
		Items:      employeeCardsToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// ListEmployeesByClinic returns the cards under a clinic.
func (h *MembershipQueryHandler) ListEmployeesByClinic(
	ctx context.Context,
	req *membershipqueryv1.ListEmployeesByClinicRequest,
) (*membershipqueryv1.ListEmployeesByClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	result, err := h.empReader.ListByClinic(ctx, caller, id, memberread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	}, memberread.EmployeeFilter{
		IncludeTerminated: req.GetIncludeTerminated(),
		OnVacation:        req.GetOnVacation(),
		Position:          strings.TrimSpace(req.GetPosition()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListEmployeesByClinicResponse{
		Items:      employeeCardsToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// ListEmployeesByOrganization returns the cards under an organization.
func (h *MembershipQueryHandler) ListEmployeesByOrganization(
	ctx context.Context,
	req *membershipqueryv1.ListEmployeesByOrganizationRequest,
) (*membershipqueryv1.ListEmployeesByOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.empReader.ListByOrganization(ctx, caller, id, memberread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	}, memberread.EmployeeFilter{
		IncludeTerminated: req.GetIncludeTerminated(),
		OnVacation:        req.GetOnVacation(),
		Position:          strings.TrimSpace(req.GetPosition()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListEmployeesByOrganizationResponse{
		Items:      employeeCardsToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// SearchEmployeesByOrganization returns cards under an organization
// whose name/email fields match the given fuzzy query. Access control
// is identical to ListEmployeesByOrganization.
func (h *MembershipQueryHandler) SearchEmployeesByOrganization(
	ctx context.Context,
	req *membershipqueryv1.SearchEmployeesByOrganizationRequest,
) (*membershipqueryv1.SearchEmployeesByOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.empReader.SearchByOrganization(ctx, caller, id,
		strings.TrimSpace(req.GetQuery()),
		memberread.ListQuery{
			Limit: int(req.GetLimit()),
			After: afterPtr(req.GetAfter()),
		},
		memberread.EmployeeFilter{
			IncludeTerminated: req.GetIncludeTerminated(),
			OnVacation:        req.GetOnVacation(),
			Position:          strings.TrimSpace(req.GetPosition()),
		})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.SearchEmployeesByOrganizationResponse{
		Items:      employeeCardsToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// CountEmployeesByDepartment returns the employee count under a department.
func (h *MembershipQueryHandler) CountEmployeesByDepartment(
	ctx context.Context,
	req *membershipqueryv1.CountEmployeesByDepartmentRequest,
) (*membershipqueryv1.CountEmployeesByDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	total, err := h.empReader.CountByDepartment(ctx, caller, id, memberread.EmployeeFilter{
		IncludeTerminated: req.GetIncludeTerminated(),
		OnVacation:        req.GetOnVacation(),
		Position:          strings.TrimSpace(req.GetPosition()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.CountEmployeesByDepartmentResponse{Total: total}, nil
}

// CountEmployeesByClinic returns the employee count under a clinic.
func (h *MembershipQueryHandler) CountEmployeesByClinic(
	ctx context.Context,
	req *membershipqueryv1.CountEmployeesByClinicRequest,
) (*membershipqueryv1.CountEmployeesByClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	total, err := h.empReader.CountByClinic(ctx, caller, id, memberread.EmployeeFilter{
		IncludeTerminated: req.GetIncludeTerminated(),
		OnVacation:        req.GetOnVacation(),
		Position:          strings.TrimSpace(req.GetPosition()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.CountEmployeesByClinicResponse{Total: total}, nil
}

// CountEmployeesByOrganization returns the employee count under an org.
func (h *MembershipQueryHandler) CountEmployeesByOrganization(
	ctx context.Context,
	req *membershipqueryv1.CountEmployeesByOrganizationRequest,
) (*membershipqueryv1.CountEmployeesByOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	total, err := h.empReader.CountByOrganization(ctx, caller, id, memberread.EmployeeFilter{
		IncludeTerminated: req.GetIncludeTerminated(),
		OnVacation:        req.GetOnVacation(),
		Position:          strings.TrimSpace(req.GetPosition()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.CountEmployeesByOrganizationResponse{Total: total}, nil
}

// CountVacationsByEmployee returns the vacation count for an employee.
func (h *MembershipQueryHandler) CountVacationsByEmployee(
	ctx context.Context,
	req *membershipqueryv1.CountVacationsByEmployeeRequest,
) (*membershipqueryv1.CountVacationsByEmployeeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	total, err := h.empReader.CountVacationsByEmployee(ctx, caller, id, req.GetState())
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.CountVacationsByEmployeeResponse{Total: total}, nil
}

// ListVacationsByEmployee returns the vacation rows for an employee.
func (h *MembershipQueryHandler) ListVacationsByEmployee(
	ctx context.Context,
	req *membershipqueryv1.ListVacationsByEmployeeRequest,
) (*membershipqueryv1.ListVacationsByEmployeeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	result, err := h.empReader.ListVacationsByEmployee(ctx, caller, id, req.GetState(), memberread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*membershipqueryv1.VacationView, 0, len(result.Items))
	for i := range result.Items {
		v := &result.Items[i]
		vi := &membershipqueryv1.VacationView{
			Id:         v.ID.String(),
			EmployeeId: v.EmployeeID.String(),
			State:      v.State,
			StartsAt:   v.StartsAt.UTC().Format(time.RFC3339Nano),
			CreatedAt:  v.CreatedAt.UTC().Format(time.RFC3339Nano),
			UpdatedAt:  v.UpdatedAt.UTC().Format(time.RFC3339Nano),
		}
		if v.EndsAt != nil {
			s := v.EndsAt.UTC().Format(time.RFC3339Nano)
			vi.EndsAt = &s
		}
		out = append(out, vi)
	}
	return &membershipqueryv1.ListVacationsByEmployeeResponse{
		Items:      out,
		NextCursor: result.NextCursor,
	}, nil
}

// GetClinicHead returns the clinic head assignment or a nil holder.
func (h *MembershipQueryHandler) GetClinicHead(
	ctx context.Context,
	req *membershipqueryv1.GetClinicHeadRequest,
) (*membershipqueryv1.GetClinicHeadResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	a, err := h.roleReader.GetClinicHead(ctx, caller, id)
	if err != nil {
		if errors.Is(err, memberread.ErrRoleVacant) {
			return &membershipqueryv1.GetClinicHeadResponse{}, nil
		}
		return nil, err
	}
	return &membershipqueryv1.GetClinicHeadResponse{Assignment: roleAssignmentToProto(a)}, nil
}

// GetDepartmentResponsible returns the responsible assignment or nil.
func (h *MembershipQueryHandler) GetDepartmentResponsible(
	ctx context.Context,
	req *membershipqueryv1.GetDepartmentResponsibleRequest,
) (*membershipqueryv1.GetDepartmentResponsibleResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	a, err := h.roleReader.GetDepartmentResponsible(ctx, caller, id)
	if err != nil {
		if errors.Is(err, memberread.ErrRoleVacant) {
			return &membershipqueryv1.GetDepartmentResponsibleResponse{}, nil
		}
		return nil, err
	}
	return &membershipqueryv1.GetDepartmentResponsibleResponse{Assignment: roleAssignmentToProto(a)}, nil
}

// ListOrgAdmins returns every org-admin holder.
func (h *MembershipQueryHandler) ListOrgAdmins(
	ctx context.Context,
	req *membershipqueryv1.ListOrgAdminsRequest,
) (*membershipqueryv1.ListOrgAdminsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.roleReader.ListOrgAdmins(ctx, caller, id, memberread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListOrgAdminsResponse{
		Items:      roleAssignmentsToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// ListOrgDispatchers returns every org-dispatcher holder.
func (h *MembershipQueryHandler) ListOrgDispatchers(
	ctx context.Context,
	req *membershipqueryv1.ListOrgDispatchersRequest,
) (*membershipqueryv1.ListOrgDispatchersResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.roleReader.ListOrgDispatchers(ctx, caller, id, memberread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListOrgDispatchersResponse{
		Items:      roleAssignmentsToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// ListOrgHeads returns every org-head holder.
func (h *MembershipQueryHandler) ListOrgHeads(
	ctx context.Context,
	req *membershipqueryv1.ListOrgHeadsRequest,
) (*membershipqueryv1.ListOrgHeadsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.roleReader.ListOrgHeads(ctx, caller, id, memberread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	return &membershipqueryv1.ListOrgHeadsResponse{
		Items:      roleAssignmentsToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// ListSystemAdmins returns every system-admin row.
func (h *MembershipQueryHandler) ListSystemAdmins(
	ctx context.Context,
	req *membershipqueryv1.ListSystemAdminsRequest,
) (*membershipqueryv1.ListSystemAdminsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	result, err := h.roleReader.ListSystemAdmins(ctx, caller, memberread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*membershipqueryv1.SystemAdminView, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, &membershipqueryv1.SystemAdminView{
			ZitadelUserId: item.ZitadelUserID,
			CreatedAt:     item.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}
	return &membershipqueryv1.ListSystemAdminsResponse{
		Items:      out,
		NextCursor: result.NextCursor,
	}, nil
}

// employeeCardToProto adapts one card.
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

// employeeCardsToProto adapts a list of cards.
func employeeCardsToProto(views []memberread.EmployeeCardView) []*membershipqueryv1.EmployeeCardView {
	out := make([]*membershipqueryv1.EmployeeCardView, 0, len(views))
	for i := range views {
		out = append(out, employeeCardToProto(&views[i]))
	}
	return out
}

// roleAssignmentToProto adapts a (possibly nil) RoleAssignmentView. The
// deputy card reuses employeeCardToProto so the nil case is handled
// uniformly with the rest of the adapter layer.
func roleAssignmentToProto(v *memberread.RoleAssignmentView) *membershipqueryv1.RoleAssignment {
	if v == nil {
		return nil
	}
	out := &membershipqueryv1.RoleAssignment{Holder: employeeCardToProto(&v.Holder)}
	if v.Deputy != nil {
		out.Deputy = employeeCardToProto(v.Deputy)
	}
	return out
}

// roleAssignmentsToProto adapts a slice of RoleAssignmentView.
func roleAssignmentsToProto(views []memberread.RoleAssignmentView) []*membershipqueryv1.RoleAssignment {
	out := make([]*membershipqueryv1.RoleAssignment, 0, len(views))
	for i := range views {
		if p := roleAssignmentToProto(&views[i]); p != nil {
			out = append(out, p)
		}
	}
	return out
}
