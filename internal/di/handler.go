package di

import (
	"github.com/samber/do/v2"

	classifierhandler "github.com/medincident/medincident-command-service/internal/handler/incident/classifier"
	membershiphandler "github.com/medincident/medincident-command-service/internal/handler/membership"
	orghandler "github.com/medincident/medincident-command-service/internal/handler/orgstructure"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/command/orgstructure"
)

func provideOrgStructureHandler(injector do.Injector) (*orghandler.OrgStructureHandler, error) {
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
	az, err := do.Invoke[*authz.Authz](injector)
	if err != nil {
		return nil, err
	}
	return orghandler.NewOrgStructureHandler(orgSvc, clinSvc, deptSvc, az), nil
}

func provideIncidentClassifierHandler(injector do.Injector) (*classifierhandler.IncidentClassifierHandler, error) {
	categorySvc, err := do.Invoke[*classifiersvc.IncidentCategoryService](injector)
	if err != nil {
		return nil, err
	}
	typeSvc, err := do.Invoke[*classifiersvc.IncidentTypeService](injector)
	if err != nil {
		return nil, err
	}
	az, err := do.Invoke[*authz.Authz](injector)
	if err != nil {
		return nil, err
	}
	return classifierhandler.NewIncidentClassifierHandler(categorySvc, typeSvc, az), nil
}

func provideMembershipHandler(injector do.Injector) (*membershiphandler.MembershipHandler, error) {
	empSvc, err := do.Invoke[*membership.EmployeeService](injector)
	if err != nil {
		return nil, err
	}
	az, err := do.Invoke[*authz.Authz](injector)
	if err != nil {
		return nil, err
	}
	return membershiphandler.NewMembershipHandler(empSvc, az), nil
}
