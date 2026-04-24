package orgstructure

import (
	"context"
	"time"

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

// GetOrganization returns the full card for one organization.
func (h *OrgStructureQueryHandler) GetOrganization(
	ctx context.Context,
	req *orgqueryv1.GetOrganizationRequest,
) (*orgqueryv1.GetOrganizationResponse, error) {
	id, err := parseOrganizationID(req.GetId())
	if err != nil {
		return nil, err
	}
	view, err := h.orgReader.Get(ctx, id)
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
	items, err := h.orgReader.List(ctx, orgread.ListQuery{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*orgqueryv1.OrganizationListItem, 0, len(items))
	for _, item := range items {
		out = append(out, &orgqueryv1.OrganizationListItem{
			Id:   item.ID.String(),
			Name: item.Name,
		})
	}
	return &orgqueryv1.ListOrganizationsResponse{Items: out}, nil
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
	caller, err := authz.CallerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
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
	caller, err := authz.CallerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	items, err := h.clinReader.ListByOrganization(ctx, caller, orgID, orgread.ListQuery{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*orgqueryv1.ClinicListItem, 0, len(items))
	for _, item := range items {
		out = append(out, &orgqueryv1.ClinicListItem{
			Id:             item.ID.String(),
			OrganizationId: item.OrganizationID.String(),
			Name:           item.Name,
		})
	}
	return &orgqueryv1.ListClinicsByOrganizationResponse{Items: out}, nil
}

// GetDepartment returns the full card for one department.
func (h *OrgStructureQueryHandler) GetDepartment(
	ctx context.Context,
	req *orgqueryv1.GetDepartmentRequest,
) (*orgqueryv1.GetDepartmentResponse, error) {
	caller, err := authz.CallerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
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
	caller, err := authz.CallerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	clinicID, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	items, err := h.deptReader.ListByClinic(ctx, caller, clinicID, orgread.ListQuery{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*orgqueryv1.DepartmentListItem, 0, len(items))
	for _, item := range items {
		out = append(out, &orgqueryv1.DepartmentListItem{
			Id:       item.ID.String(),
			ClinicId: item.ClinicID.String(),
			Name:     item.Name,
		})
	}
	return &orgqueryv1.ListDepartmentsByClinicResponse{Items: out}, nil
}

// organizationToProto adapts an OrganizationDetails into the wire type.
func organizationToProto(v *orgread.OrganizationDetails) *orgqueryv1.Organization {
	return &orgqueryv1.Organization{
		Id:           v.ID.String(),
		Name:         v.Name,
		Description:  v.Description,
		LegalAddress: addressViewToProto(v.LegalAddress),
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
