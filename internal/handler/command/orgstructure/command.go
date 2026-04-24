// Package orgstructure is the gRPC transport for the OrgStructureCommandService.
package orgstructure

import (
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

// OrgStructureHandler implements orgstructurev1.OrgStructureCommandServiceServer.
// Pure transport: translates proto → service Command + Payload, wraps
// the authenticated caller in an authz.Caller. Authorization happens
// inside each service method, so the handler does not depend on
// *authz.Authz, and UUID-parse guards are no longer needed at the
// transport boundary.
type OrgStructureHandler struct {
	orgstructurev1.UnimplementedOrgStructureCommandServiceServer

	orgSvc  *orgsvc.OrganizationService
	clinSvc *orgsvc.ClinicService
	deptSvc *orgsvc.DepartmentService
}

// NewOrgStructureHandler wires the handler with the three services.
func NewOrgStructureHandler(
	orgSvc *orgsvc.OrganizationService,
	clinSvc *orgsvc.ClinicService,
	deptSvc *orgsvc.DepartmentService,
) *OrgStructureHandler {
	return &OrgStructureHandler{
		orgSvc:  orgSvc,
		clinSvc: clinSvc,
		deptSvc: deptSvc,
	}
}

// addressInputFromProto converts a proto AddressInput into a service
// AddressInput. A nil proto pointer becomes a zero AddressInput (Text
// will then fail the required check in the service validator).
func addressInputFromProto(in *orgstructurev1.AddressInput) orgsvc.AddressInput {
	if in == nil {
		return orgsvc.AddressInput{}
	}
	out := orgsvc.AddressInput{Text: in.GetText()}
	if p := in.GetPoint(); p != nil {
		out.Point = &orgsvc.PointInput{
			Longitude: p.GetLongitude(),
			Latitude:  p.GetLatitude(),
		}
	}
	return out
}
