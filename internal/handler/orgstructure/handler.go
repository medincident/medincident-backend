// Package orgstructure is the gRPC transport for the OrgStructureService.
package orgstructure

import (
	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/service/authz"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/service/orgstructure/v1"
)

// Error codes emitted by handler-layer request parsing.
const (
	ErrCodeHandlerInvalidOrganizationID = "handler_invalid_organization_id"
	ErrCodeHandlerInvalidClinicID       = "handler_invalid_clinic_id"
	ErrCodeHandlerInvalidDepartmentID   = "handler_invalid_department_id"
)

// OrgStructureHandler implements orgstructurev1.OrgStructureServiceServer.
type OrgStructureHandler struct {
	orgstructurev1.UnimplementedOrgStructureServiceServer

	authz   *authz.Authz
	orgSvc  *orgsvc.OrganizationService
	clinSvc *orgsvc.ClinicService
	deptSvc *orgsvc.DepartmentService
}

// NewOrgStructureHandler wires the handler with the three services.
func NewOrgStructureHandler(
	orgSvc *orgsvc.OrganizationService,
	clinSvc *orgsvc.ClinicService,
	deptSvc *orgsvc.DepartmentService,
	az *authz.Authz,
) *OrgStructureHandler {
	return &OrgStructureHandler{
		authz:   az,
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

// parseOrganizationID is shared by handler files that take an
// organization_id from the proto request.
func parseOrganizationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.orgstructure").
			Code(ErrCodeHandlerInvalidOrganizationID).
			Public("Invalid organization id.").
			With("organization_id", raw).
			Wrap(err)
	}
	return id, nil
}

// parseClinicID is shared by handler files that take a clinic_id.
func parseClinicID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.orgstructure").
			Code(ErrCodeHandlerInvalidClinicID).
			Public("Invalid clinic id.").
			With("clinic_id", raw).
			Wrap(err)
	}
	return id, nil
}

// parseDepartmentID is shared by handler files that take a department_id.
func parseDepartmentID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.orgstructure").
			Code(ErrCodeHandlerInvalidDepartmentID).
			Public("Invalid department id.").
			With("department_id", raw).
			Wrap(err)
	}
	return id, nil
}
