package di

import (
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/services/membership"
	orgsvc "github.com/medincident/medincident-command-service/internal/services/orgstructure"
	"github.com/medincident/medincident-command-service/internal/zitadel"
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

// ProvideEmployeeService wires EmployeeService.
func ProvideEmployeeService(injector do.Injector) (*membership.EmployeeService, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	verifier, err := do.Invoke[zitadel.UserVerifier](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return membership.NewEmployeeService(db, verifier, logger), nil
}
