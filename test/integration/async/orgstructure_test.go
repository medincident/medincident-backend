//go:build integration

package async_integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
)

func TestOrganizationProjectedAsync(t *testing.T) {
	resetDBs(t)
	ctx := context.Background()

	result, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:         "Async Test Org",
			LegalAddress: orgsvc.AddressInput{Text: "ул. Тестовая, 1"},
		},
	})
	require.NoError(t, err)

	requireProjected(t, `SELECT count(*) FROM projections.organizations WHERE id = ?`, result.ID)

	var name string
	require.NoError(t, queryDB.Raw(`SELECT name FROM projections.organizations WHERE id = ?`, result.ID).Scan(&name).Error)
	assert.Equal(t, "Async Test Org", name)
}

func TestClinicProjectedAsync(t *testing.T) {
	resetDBs(t)
	ctx := context.Background()

	org, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{Name: "Org For Clinic", LegalAddress: orgsvc.AddressInput{Text: "addr"}},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.organizations WHERE id = ?`, org.ID)

	clin, err := clinSvc.Create(ctx, orgsvc.CreateClinicCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{
			OrganizationID:  org.ID.String(),
			Name:            "Async Clinic",
			PhysicalAddress: orgsvc.AddressInput{Text: "clinic addr"},
		},
	})
	require.NoError(t, err)

	requireProjected(t, `SELECT count(*) FROM projections.clinics WHERE id = ?`, clin.ID)

	var count int
	require.NoError(t, queryDB.Raw(`SELECT clinics_total FROM projections.organization_counters WHERE organization_id = ?`, org.ID).Scan(&count).Error)
	assert.Equal(t, 1, count)
}

func TestDepartmentProjectedAsync(t *testing.T) {
	resetDBs(t)
	ctx := context.Background()

	org, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{Name: "Org For Dept", LegalAddress: orgsvc.AddressInput{Text: "addr"}},
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
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clin.ID.String(), Name: "Async Dept"},
	})
	require.NoError(t, err)

	requireProjected(t, `SELECT count(*) FROM projections.departments WHERE id = ?`, dept.ID)

	var deptCount int
	require.NoError(t, queryDB.Raw(`SELECT departments_total FROM projections.clinic_counters WHERE clinic_id = ?`, clin.ID).Scan(&deptCount).Error)
	assert.Equal(t, 1, deptCount)
}
