//go:build integration

// Package async_integration_test also covers НС (несчастный случай /
// incidents): the IncidentCreated projector requires registrar_organization_id,
// registrar_clinic_id, and registrar_department_id which must be looked up
// from projections.employees.  This test pins that end-to-end invariant.
package async_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
)

// TestIncidentProjectedAsync verifies the full CQRS pipeline for incidents:
// command → outbox → publisher → NATS → domain consumer → projections.incidents.
//
// This test specifically catches the bug where the IncidentCreated projector
// omitted registrar_organization_id / registrar_clinic_id / registrar_department_id,
// causing a NOT NULL constraint violation and silent drop of every incident event.
func TestIncidentProjectedAsync(t *testing.T) {
	resetDBs(t)
	ctx := context.Background()

	// ── 1. Org structure ────────────────────────────────────────────────
	org, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{Name: "Incident Org", LegalAddress: orgsvc.AddressInput{Text: "addr"}},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.organizations WHERE id = ?`, org.ID)

	clin, err := clinSvc.Create(ctx, orgsvc.CreateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{OrganizationID: org.ID.String(), Name: "Clinic", PhysicalAddress: orgsvc.AddressInput{Text: "addr"}},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.clinics WHERE id = ?`, clin.ID)

	dept, err := deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clin.ID.String(), Name: "Dept"},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.departments WHERE id = ?`, dept.ID)

	// ── 2. Hire registrar employee (direct insert + outbox event) ───────
	empID, err := uuid.NewV7()
	require.NoError(t, err)
	appendEmployeeHiredEvent(t, empID, sysadminZitadelID, org.ID, dept.ID, nil)
	requireProjected(t, `SELECT count(*) FROM projections.employees WHERE id = ?`, empID)

	// ── 3. Incident classifier ───────────────────────────────────────────
	cat, err := classifierCatSvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{OrganizationID: org.ID.String(), Name: "Category"},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.incident_categories WHERE id = ?`, cat.ID)

	typ, err := classifierTypSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{CategoryID: cat.ID.String(), Name: "Type"},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.incident_types WHERE id = ?`, typ.ID)

	// ── 4. Create incident ───────────────────────────────────────────────
	// sysadminCaller is now an employee of org (inserted above), so
	// authz.MemberOf.Organization passes.
	inc, err := incidentSvc.Create(ctx, &incidentsvc.CreateIncidentCommand{
		Caller: sysadminCaller,
		Payload: incidentsvc.CreateIncidentPayload{
			DepartmentID: dept.ID.String(),
			CategoryID:   cat.ID.String(),
			TypeID:       typ.ID.String(),
			OccurredAt:   time.Now().Add(-time.Hour).UTC().Format(time.RFC3339Nano),
		},
	})
	require.NoError(t, err)

	// ── 5. Assert projection ─────────────────────────────────────────────
	requireProjected(t, `SELECT count(*) FROM projections.incidents WHERE id = ?`, inc.ID)

	var orgIDVal, clinIDVal, deptIDVal, regOrgIDVal, regClinIDVal, regDeptIDVal string
	require.NoError(t, queryDB.Raw(`
		SELECT organization_id::text, clinic_id::text, department_id::text,
		       registrar_organization_id::text, registrar_clinic_id::text, registrar_department_id::text
		  FROM projections.incidents WHERE id = ?`, inc.ID,
	).Row().Scan(&orgIDVal, &clinIDVal, &deptIDVal, &regOrgIDVal, &regClinIDVal, &regDeptIDVal))

	assert.Equal(t, org.ID.String(), orgIDVal, "projections.incidents.organization_id")
	assert.Equal(t, clin.ID.String(), clinIDVal, "projections.incidents.clinic_id")
	assert.Equal(t, dept.ID.String(), deptIDVal, "projections.incidents.department_id")
	assert.Equal(t, org.ID.String(), regOrgIDVal, "projections.incidents.registrar_organization_id")
	assert.Equal(t, clin.ID.String(), regClinIDVal, "projections.incidents.registrar_clinic_id")
	assert.Equal(t, dept.ID.String(), regDeptIDVal, "projections.incidents.registrar_department_id")
}
