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

func isDeptActive(t *testing.T, id uuid.UUID) bool {
	t.Helper()
	var active bool
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.departments WHERE id = ?`, id).Row().Scan(&active))
	return active
}

// seedDepartment creates org → clinic → department and returns the department id.
func seedDepartment(t *testing.T) uuid.UUID {
	t.Helper()
	clinicID := seedClinic(t)
	res, err := deptSvc.Create(context.Background(), orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clinicID.String(), Name: "Родительский отдел"},
	})
	require.NoError(t, err)
	return res.ID
}

func TestDepartment_Deactivate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	deptID := seedDepartment(t)

	require.NoError(t, deptSvc.Deactivate(ctx, orgsvc.DeactivateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateDepartmentPayload{ID: deptID.String()},
	}))

	assert.False(t, isDeptActive(t, deptID))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.department.v1.deactivated"))
}

func TestDepartment_Deactivate_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := deptSvc.Deactivate(ctx, orgsvc.DeactivateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateDepartmentPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeDepartmentNotFound, codeOf(t, err))
}

func TestDepartment_Activate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	deptID := seedDepartment(t)

	require.NoError(t, deptSvc.Deactivate(ctx, orgsvc.DeactivateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateDepartmentPayload{ID: deptID.String()},
	}))

	require.NoError(t, deptSvc.Activate(ctx, orgsvc.ActivateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateDepartmentPayload{ID: deptID.String()},
	}))

	assert.True(t, isDeptActive(t, deptID))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.department.v1.activated"))
}

func TestDepartment_Activate_ParentClinicInactive(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	deptID := seedDepartment(t)

	// Deactivate the parent clinic (cascades to department).
	var clinicID uuid.UUID
	require.NoError(t, testDB.Raw(`SELECT clinic_id FROM domain.departments WHERE id = ?`, deptID).Row().Scan(&clinicID))
	require.NoError(t, clinSvc.Deactivate(ctx, orgsvc.DeactivateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeactivateClinicPayload{ID: clinicID.String()},
	}))

	err := deptSvc.Activate(ctx, orgsvc.ActivateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateDepartmentPayload{ID: deptID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeDepartmentActivateParentInactive, codeOf(t, err))
}

func TestDepartment_Activate_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := deptSvc.Activate(ctx, orgsvc.ActivateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.ActivateDepartmentPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeDepartmentNotFound, codeOf(t, err))
}

func TestDepartment_Delete_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	deptID := seedDepartment(t)

	require.NoError(t, deptSvc.Delete(ctx, orgsvc.DeleteDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeleteDepartmentPayload{ID: deptID.String()},
	}))

	assert.Equal(t, 0, countDepartments(t))
	assert.Equal(t, 1, countOutboxEvents(t, "medincident.event.department.v1.deleted"))
}

func TestDepartment_Delete_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := deptSvc.Delete(ctx, orgsvc.DeleteDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.DeleteDepartmentPayload{ID: uuid.New().String()},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeDepartmentNotFound, codeOf(t, err))
}
