package orgstructure

import (
	"context"
	"strings"
	"time"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgread "github.com/medincident/medincident-backend/internal/service/query/orgstructure"
	orgqueryv1 "github.com/medincident/medincident-backend/pkg/query/orgstructure/v1"
)

// OrgStructureQueryHandler implements orgqueryv1.OrgStructureQueryServiceServer
// by adapting the per-domain readers to proto messages. Time values
// are emitted as RFC3339Nano strings; optional columns are wired
// through proto's `optional` scalar representation (pointer in Go).
type OrgStructureQueryHandler struct {
	orgqueryv1.UnimplementedOrgStructureQueryServiceServer

	orgReader  *orgread.OrganizationReader
	clinReader *orgread.ClinicReader
	deptReader *orgread.DepartmentReader
}

// NewOrgStructureQueryHandler wires the handler with the three readers.
func NewOrgStructureQueryHandler(
	orgReader *orgread.OrganizationReader,
	clinReader *orgread.ClinicReader,
	deptReader *orgread.DepartmentReader,
) *OrgStructureQueryHandler {
	return &OrgStructureQueryHandler{
		orgReader:  orgReader,
		clinReader: clinReader,
		deptReader: deptReader,
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

// GetOrganization returns the full card for one organization.
func (h *OrgStructureQueryHandler) GetOrganization(
	ctx context.Context,
	req *orgqueryv1.GetOrganizationRequest,
) (*orgqueryv1.GetOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetId())
	if err != nil {
		return nil, err
	}
	view, err := h.orgReader.Get(ctx, caller, id)
	if err != nil {
		return nil, err
	}
	return &orgqueryv1.GetOrganizationResponse{
		Organization: organizationToProto(view),
	}, nil
}

// ListOrganizations returns a paginated list of organizations.
func (h *OrgStructureQueryHandler) ListOrganizations(
	ctx context.Context,
	req *orgqueryv1.ListOrganizationsRequest,
) (*orgqueryv1.ListOrganizationsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	result, err := h.orgReader.List(ctx, caller, req.GetIncludeDeactivated(), orgread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*orgqueryv1.OrganizationListItem, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, &orgqueryv1.OrganizationListItem{
			Id:       item.ID.String(),
			Name:     item.Name,
			IsActive: item.IsActive,
		})
	}
	return &orgqueryv1.ListOrganizationsResponse{
		Items:      out,
		NextCursor: result.NextCursor,
	}, nil
}

// SearchOrganizations returns organizations whose name matches the
// given case-insensitive substring. Empty query degenerates to
// ListOrganizations semantics. When include_deactivated is true the
// caller must be a system admin; otherwise only active organizations
// are searched.
func (h *OrgStructureQueryHandler) SearchOrganizations(
	ctx context.Context,
	req *orgqueryv1.SearchOrganizationsRequest,
) (*orgqueryv1.SearchOrganizationsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	result, err := h.orgReader.Search(ctx, caller, strings.TrimSpace(req.GetQuery()), req.GetIncludeDeactivated(), orgread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*orgqueryv1.OrganizationListItem, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, &orgqueryv1.OrganizationListItem{
			Id:       item.ID.String(),
			Name:     item.Name,
			IsActive: item.IsActive,
		})
	}
	return &orgqueryv1.SearchOrganizationsResponse{
		Items:      out,
		NextCursor: result.NextCursor,
	}, nil
}

// CountOrganizations returns the total organizations count.
func (h *OrgStructureQueryHandler) CountOrganizations(
	ctx context.Context,
	_ *orgqueryv1.CountOrganizationsRequest,
) (*orgqueryv1.CountOrganizationsResponse, error) {
	total, err := h.orgReader.Count(ctx)
	if err != nil {
		return nil, err
	}
	return &orgqueryv1.CountOrganizationsResponse{Total: total}, nil
}

// GetClinic returns the full card for one clinic.
func (h *OrgStructureQueryHandler) GetClinic(
	ctx context.Context,
	req *orgqueryv1.GetClinicRequest,
) (*orgqueryv1.GetClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseClinicID(req.GetId())
	if err != nil {
		return nil, err
	}
	view, err := h.clinReader.Get(ctx, caller, id)
	if err != nil {
		return nil, err
	}
	return &orgqueryv1.GetClinicResponse{Clinic: clinicToProto(view)}, nil
}

// ListClinicsByOrganization returns clinics under one organization.
func (h *OrgStructureQueryHandler) ListClinicsByOrganization(
	ctx context.Context,
	req *orgqueryv1.ListClinicsByOrganizationRequest,
) (*orgqueryv1.ListClinicsByOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.clinReader.ListByOrganization(ctx, caller, orgID, req.GetIncludeDeactivated(), orgread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*orgqueryv1.ClinicListItem, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, &orgqueryv1.ClinicListItem{
			Id:             item.ID.String(),
			OrganizationId: item.OrganizationID.String(),
			Name:           item.Name,
			IsActive:       item.IsActive,
		})
	}
	return &orgqueryv1.ListClinicsByOrganizationResponse{
		Items:      out,
		NextCursor: result.NextCursor,
	}, nil
}

// CountClinicsByOrganization returns the total clinics count for one
// organization.
func (h *OrgStructureQueryHandler) CountClinicsByOrganization(
	ctx context.Context,
	req *orgqueryv1.CountClinicsByOrganizationRequest,
) (*orgqueryv1.CountClinicsByOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	total, err := h.clinReader.CountByOrganization(ctx, caller, orgID)
	if err != nil {
		return nil, err
	}
	return &orgqueryv1.CountClinicsByOrganizationResponse{Total: total}, nil
}

// GetDepartment returns the full card for one department.
func (h *OrgStructureQueryHandler) GetDepartment(
	ctx context.Context,
	req *orgqueryv1.GetDepartmentRequest,
) (*orgqueryv1.GetDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseDepartmentID(req.GetId())
	if err != nil {
		return nil, err
	}
	view, err := h.deptReader.Get(ctx, caller, id)
	if err != nil {
		return nil, err
	}
	return &orgqueryv1.GetDepartmentResponse{Department: departmentToProto(view)}, nil
}

// ListDepartmentsByClinic returns departments under one clinic.
func (h *OrgStructureQueryHandler) ListDepartmentsByClinic(
	ctx context.Context,
	req *orgqueryv1.ListDepartmentsByClinicRequest,
) (*orgqueryv1.ListDepartmentsByClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	clinicID, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	result, err := h.deptReader.ListByClinic(ctx, caller, clinicID, orgread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*orgqueryv1.DepartmentListItem, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, &orgqueryv1.DepartmentListItem{
			Id:       item.ID.String(),
			ClinicId: item.ClinicID.String(),
			Name:     item.Name,
			IsActive: item.IsActive,
		})
	}
	return &orgqueryv1.ListDepartmentsByClinicResponse{
		Items:      out,
		NextCursor: result.NextCursor,
	}, nil
}

// CountDepartmentsByClinic returns the total departments count for one
// clinic.
func (h *OrgStructureQueryHandler) CountDepartmentsByClinic(
	ctx context.Context,
	req *orgqueryv1.CountDepartmentsByClinicRequest,
) (*orgqueryv1.CountDepartmentsByClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	clinicID, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	total, err := h.deptReader.CountByClinic(ctx, caller, clinicID)
	if err != nil {
		return nil, err
	}
	return &orgqueryv1.CountDepartmentsByClinicResponse{Total: total}, nil
}

// organizationToProto adapts an OrganizationDetails into the wire type.
func organizationToProto(v *orgread.OrganizationDetails) *orgqueryv1.Organization {
	return &orgqueryv1.Organization{
		Id:           v.ID.String(),
		Name:         v.Name,
		Description:  v.Description,
		LegalAddress: addressViewToProto(v.LegalAddress),
		IsActive:     v.IsActive,
		CreatedAt:    v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:    v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

// clinicToProto adapts a ClinicDetails into the wire type.
func clinicToProto(v *orgread.ClinicDetails) *orgqueryv1.Clinic {
	return &orgqueryv1.Clinic{
		Id:              v.ID.String(),
		OrganizationId:  v.OrganizationID.String(),
		Name:            v.Name,
		Description:     v.Description,
		PhysicalAddress: addressViewToProto(v.PhysicalAddress),
		IsActive:        v.IsActive,
		CreatedAt:       v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:       v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

// departmentToProto adapts a DepartmentDetails into the wire type.
func departmentToProto(v *orgread.DepartmentDetails) *orgqueryv1.Department {
	return &orgqueryv1.Department{
		Id:          v.ID.String(),
		ClinicId:    v.ClinicID.String(),
		Name:        v.Name,
		Description: v.Description,
		IsActive:    v.IsActive,
		CreatedAt:   v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

// addressViewToProto adapts an AddressView into the wire type.
func addressViewToProto(v orgread.AddressView) *orgqueryv1.Address {
	out := &orgqueryv1.Address{Text: v.Text}
	if v.Point != nil {
		out.Point = &orgqueryv1.Point{
			Longitude: v.Point.Longitude,
			Latitude:  v.Point.Latitude,
		}
	}
	return out
}
