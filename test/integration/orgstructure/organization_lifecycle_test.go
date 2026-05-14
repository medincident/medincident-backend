//go:build integration

package orgstructure_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
)

func isOrgActive(t *testing.T, id uuid.UUID) bool {
	t.Helper()
	var active bool
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.organizations WHERE id = ?`, id).Row().Scan(&active))
	return active
}

func TestOrganization_Deactivate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	require.NoError(t, orgSvc.Deactivate(ctx, orgsvc.DeactivateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateOrganizationPayload{ID: orgID.String()},
	}))

	assert.False(t, isOrgActive(t, orgID))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.organization.v1.deactivated"))
}

func TestOrganization_Deactivate_CascadesToClinicsAndDepartments(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	// Create two clinics, each with one department.
	clin1, err := clinSvc.Create(ctx, orgsvc.CreateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{OrganizationID: orgID.String(), Name: "Clinic 1", PhysicalAddress: orgsvc.AddressInput{Text: "addr1"}},
	})
	require.NoError(t, err)
	clin2, err := clinSvc.Create(ctx, orgsvc.CreateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{OrganizationID: orgID.String(), Name: "Clinic 2", PhysicalAddress: orgsvc.AddressInput{Text: "addr2"}},
	})
	require.NoError(t, err)
	dept1, err := deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clin1.ID.String(), Name: "Dept 1"},
	})
	require.NoError(t, err)

	require.NoError(t, orgSvc.Deactivate(ctx, orgsvc.DeactivateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateOrganizationPayload{ID: orgID.String()},
	}))

	assert.False(t, isOrgActive(t, orgID))

	var clinicActive bool
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.clinics WHERE id = ?`, clin1.ID).Row().Scan(&clinicActive))
	assert.False(t, clinicActive, "clinic1 should be deactivated")
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.clinics WHERE id = ?`, clin2.ID).Row().Scan(&clinicActive))
	assert.False(t, clinicActive, "clinic2 should be deactivated")

	var deptActive bool
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.departments WHERE id = ?`, dept1.ID).Row().Scan(&deptActive))
	assert.False(t, deptActive, "dept1 should be deactivated")

	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.organization.v1.deactivated"))
	assert.Equal(t, 2, countOutboxEvents(t, "medincident.event.clinic.v1.deactivated"))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.department.v1.deactivated"))
}

func TestOrganization_Deactivate_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := orgSvc.Deactivate(ctx, orgsvc.DeactivateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateOrganizationPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeOrganizationNotFound, codeOf(t, err))
}

func TestOrganization_Activate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	require.NoError(t, orgSvc.Deactivate(ctx, orgsvc.DeactivateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateOrganizationPayload{ID: orgID.String()},
	}))
	require.False(t, isOrgActive(t, orgID))

	require.NoError(t, orgSvc.Activate(ctx, orgsvc.ActivateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateOrganizationPayload{ID: orgID.String()},
	}))

	assert.True(t, isOrgActive(t, orgID))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.organization.v1.activated"))
}

func TestOrganization_Activate_Idempotent(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	// Already active — should succeed without writing an event.
	require.NoError(t, orgSvc.Activate(ctx, orgsvc.ActivateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateOrganizationPayload{ID: orgID.String()},
	}))

	assert.Equal(t, 0, countOutboxEvents(t, "medincident.event.organization.v1.activated"))
}

func TestOrganization_Activate_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := orgSvc.Activate(ctx, orgsvc.ActivateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateOrganizationPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeOrganizationNotFound, codeOf(t, err))
}

func TestOrganization_Delete_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	require.NoError(t, orgSvc.Delete(ctx, orgsvc.DeleteOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeleteOrganizationPayload{ID: orgID.String()},
	}))

	assert.Equal(t, 0, countOrganizations(t))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.organization.v1.deleted"))
}

func TestOrganization_Delete_CascadesChildren(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	clinRes, err := clinSvc.Create(ctx, orgsvc.CreateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{OrganizationID: orgID.String(), Name: "Clinic X", PhysicalAddress: orgsvc.AddressInput{Text: "addr"}},
	})
	require.NoError(t, err)
	_, err = deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clinRes.ID.String(), Name: "Dept X"},
	})
	require.NoError(t, err)

	require.NoError(t, orgSvc.Delete(ctx, orgsvc.DeleteOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeleteOrganizationPayload{ID: orgID.String()},
	}))

	assert.Equal(t, 0, countOrganizations(t))
	assert.Equal(t, 0, countClinics(t))
	assert.Equal(t, 0, countDepartments(t))
}

func TestOrganization_Delete_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := orgSvc.Delete(ctx, orgsvc.DeleteOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeleteOrganizationPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeOrganizationNotFound, codeOf(t, err))
}
