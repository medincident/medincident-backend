package authz

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// resolveOrgID runs a query that returns a single organization_id::text
// and parses it into uuid.UUID. Returns ErrCodeScopeResolveFailed with
// publicMsg if no row is found.
func (a *Authz) resolveOrgID(ctx context.Context, query string, arg uuid.UUID, entityName string, entityID uuid.UUID) (uuid.UUID, error) {
	var orgIDStr *string
	if err := a.db.WithContext(ctx).
		Raw(query, arg).
		Scan(&orgIDStr).Error; err != nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeAuthzQueryFailed).
			With(entityName+"_id", entityID).Wrap(err)
	}
	if orgIDStr == nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			Public(strings.ToUpper(entityName[:1])+entityName[1:]+" not found.").
			With(entityName+"_id", entityID).
			Errorf("%s not found for authz scope", entityName)
	}
	id, err := uuid.Parse(*orgIDStr)
	if err != nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeAuthzQueryFailed).Wrap(err)
	}
	return id, nil
}

func (a *Authz) resolveOrgByClinic(ctx context.Context, clinicID uuid.UUID) (uuid.UUID, error) {
	return a.resolveOrgID(ctx,
		`SELECT organization_id::text FROM domain.clinics WHERE id = ?`,
		clinicID, "clinic", clinicID)
}

func (a *Authz) resolveOrgByDepartment(ctx context.Context, deptID uuid.UUID) (uuid.UUID, error) {
	return a.resolveOrgID(ctx,
		`SELECT c.organization_id::text FROM domain.departments d JOIN domain.clinics c ON c.id = d.clinic_id WHERE d.id = ?`,
		deptID, "department", deptID)
}

func (a *Authz) resolveOrgByEmployee(ctx context.Context, empID uuid.UUID) (uuid.UUID, error) {
	return a.resolveOrgID(ctx,
		`SELECT organization_id::text FROM domain.employees WHERE id = ?`,
		empID, "employee", empID)
}

func (a *Authz) resolveOrgByVacation(ctx context.Context, vacID uuid.UUID) (uuid.UUID, error) {
	return a.resolveOrgID(ctx,
		`SELECT e.organization_id::text FROM domain.employee_vacations v JOIN domain.employees e ON e.id = v.employee_id WHERE v.id = ?`,
		vacID, "vacation", vacID)
}

func (a *Authz) resolveOrgByCategory(ctx context.Context, catID uuid.UUID) (uuid.UUID, error) {
	return a.resolveOrgID(ctx,
		`SELECT organization_id::text FROM domain.incident_categories WHERE id = ?`,
		catID, "category", catID)
}

func (a *Authz) resolveOrgByType(ctx context.Context, typeID uuid.UUID) (uuid.UUID, error) {
	return a.resolveOrgID(ctx,
		`SELECT organization_id::text FROM domain.incident_types WHERE id = ?`,
		typeID, "type", typeID)
}
