package di

import (
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"

	orghandler "github.com/medincident/medincident-command-service/internal/handler/orgstructure"
	orgsvc "github.com/medincident/medincident-command-service/internal/services/orgstructure"
)

// ProvideOrgStructureHandler wires the gRPC handler for OrgStructureService.
func ProvideOrgStructureHandler(injector do.Injector) (*orghandler.OrgStructureHandler, error) {
	orgSvc, err := do.Invoke[*orgsvc.OrganizationService](injector)
	if err != nil {
		return nil, err
	}
	clinSvc, err := do.Invoke[*orgsvc.ClinicService](injector)
	if err != nil {
		return nil, err
	}
	deptSvc, err := do.Invoke[*orgsvc.DepartmentService](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return orghandler.NewOrgStructureHandler(orgSvc, clinSvc, deptSvc, logger), nil
}
