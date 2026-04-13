package di

import (
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	orgsvc "github.com/medincident/medincident-command-service/internal/services/orgstructure"
)

// ProvideOrganizationService wires OrganizationService.
func ProvideOrganizationService(injector do.Injector) (*orgsvc.OrganizationService, error) {
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

// ProvideClinicService wires ClinicService.
func ProvideClinicService(injector do.Injector) (*orgsvc.ClinicService, error) {
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

// ProvideDepartmentService wires DepartmentService.
func ProvideDepartmentService(injector do.Injector) (*orgsvc.DepartmentService, error) {
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
