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

func isClinicActive(t *testing.T, id uuid.UUID) bool {
	t.Helper()
	var active bool
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.clinics WHERE id = ?`, id).Row().Scan(&active))
	return active
}

func TestClinic_Deactivate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	require.NoError(t, clinSvc.Deactivate(ctx, orgsvc.DeactivateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateClinicPayload{ID: clinicID.String()},
	}))

	assert.False(t, isClinicActive(t, clinicID))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.clinic.v1.deactivated"))
}

func TestClinic_Deactivate_CascadesToDepartments(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	dept1, err := deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clinicID.String(), Name: "Dept 1"},
	})
	require.NoError(t, err)
	dept2, err := deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clinicID.String(), Name: "Dept 2"},
	})
	require.NoError(t, err)

	require.NoError(t, clinSvc.Deactivate(ctx, orgsvc.DeactivateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateClinicPayload{ID: clinicID.String()},
	}))

	assert.False(t, isClinicActive(t, clinicID))

	var deptActive bool
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.departments WHERE id = ?`, dept1.ID).Row().Scan(&deptActive))
	assert.False(t, deptActive)
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.departments WHERE id = ?`, dept2.ID).Row().Scan(&deptActive))
	assert.False(t, deptActive)

	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.clinic.v1.deactivated"))
	assert.Equal(t, 2, countOutboxEvents(t, "medincident.event.department.v1.deactivated"))
}

func TestClinic_Deactivate_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := clinSvc.Deactivate(ctx, orgsvc.DeactivateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateClinicPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeClinicNotFound, codeOf(t, err))
}

func TestClinic_Activate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	require.NoError(t, clinSvc.Deactivate(ctx, orgsvc.DeactivateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateClinicPayload{ID: clinicID.String()},
	}))

	require.NoError(t, clinSvc.Activate(ctx, orgsvc.ActivateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateClinicPayload{ID: clinicID.String()},
	}))

	assert.True(t, isClinicActive(t, clinicID))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.clinic.v1.activated"))
}

func TestClinic_Activate_ParentOrgInactive(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	// Deactivate the org first (also deactivates the clinic via cascade).
	var orgID uuid.UUID
	require.NoError(t, testDB.Raw(`SELECT organization_id FROM domain.clinics WHERE id = ?`, clinicID).Row().Scan(&orgID))
	require.NoError(t, orgSvc.Deactivate(ctx, orgsvc.DeactivateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateOrganizationPayload{ID: orgID.String()},
	}))

	err := clinSvc.Activate(ctx, orgsvc.ActivateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateClinicPayload{ID: clinicID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeClinicActivateParentInactive, codeOf(t, err))
}

func TestClinic_Activate_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := clinSvc.Activate(ctx, orgsvc.ActivateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateClinicPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeClinicNotFound, codeOf(t, err))
}

func TestClinic_Delete_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	require.NoError(t, clinSvc.Delete(ctx, orgsvc.DeleteClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeleteClinicPayload{ID: clinicID.String()},
	}))

	assert.Equal(t, 0, countClinics(t))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.clinic.v1.deleted"))
}

func TestClinic_Delete_CascadesDepartments(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	_, err := deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clinicID.String(), Name: "Dept X"},
	})
	require.NoError(t, err)

	require.NoError(t, clinSvc.Delete(ctx, orgsvc.DeleteClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeleteClinicPayload{ID: clinicID.String()},
	}))

	assert.Equal(t, 0, countClinics(t))
	assert.Equal(t, 0, countDepartments(t))
}

func TestClinic_Delete_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := clinSvc.Delete(ctx, orgsvc.DeleteClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeleteClinicPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeClinicNotFound, codeOf(t, err))
}
