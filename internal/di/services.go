package di

import (
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	"github.com/medincident/medincident-backend/internal/service/zitadel"
)

func provideOrganizationService(injector do.Injector) (*orgsvc.OrganizationService, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return orgsvc.NewOrganizationService(db, logger), nil
}

func provideClinicService(injector do.Injector) (*orgsvc.ClinicService, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return orgsvc.NewClinicService(db, logger), nil
}

func provideDepartmentService(injector do.Injector) (*orgsvc.DepartmentService, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return orgsvc.NewDepartmentService(db, logger), nil
}

func provideIncidentCategoryService(injector do.Injector) (*classifiersvc.IncidentCategoryService, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return classifiersvc.NewIncidentCategoryService(db, logger), nil
}

func provideIncidentTypeService(injector do.Injector) (*classifiersvc.IncidentTypeService, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return classifiersvc.NewIncidentTypeService(db, logger), nil
}

func provideEmployeeService(injector do.Injector) (*membership.EmployeeService, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	verifier, err := do.Invoke[*zitadel.Service](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return membership.NewEmployeeService(db, verifier, logger), nil
}
