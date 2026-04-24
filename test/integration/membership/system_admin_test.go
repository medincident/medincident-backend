//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

func TestGrantSystemAdmin_Success(t *testing.T) {
	_ = takeFixture(t)

	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.GrantSystemAdminPayload{
			ZitadelUserID: testUserAliceID,
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.system_admins WHERE zitadel_user_id = ?`, testUserAliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count)

	// Sync projector must have written the matching projection row.
	var projCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM projections.system_admins WHERE zitadel_user_id = ?`, testUserAliceID,
	).Scan(&projCount).Error)
	assert.Equal(t, int64(1), projCount)
}

func TestGrantSystemAdmin_ZitadelUserNotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.GrantSystemAdminPayload{
			ZitadelUserID: "nonexistent-user-id",
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeZitadelUserNotFound, oopsCode(t, err))
}

func TestGrantSystemAdmin_AlreadyGranted(t *testing.T) {
	_ = takeFixture(t)
	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.GrantSystemAdminPayload{
			ZitadelUserID: testUserAliceID,
		},
	}))
	err := empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.GrantSystemAdminPayload{
			ZitadelUserID: testUserAliceID,
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeSystemAdminAlreadyGranted, oopsCode(t, err))
}

func TestGrantSystemAdmin_WhitespaceOnlyInputRejected(t *testing.T) {
	err := empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		Caller:  sysadminCaller,
		Payload: membership.GrantSystemAdminPayload{ZitadelUserID: "   "},
	})
	require.Error(t, err)
	assert.Equal(t, validation.CodeStringExtraWhitespace, oopsCode(t, err))
}

func TestRevokeSystemAdmin_Success(t *testing.T) {
	_ = takeFixture(t)
	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.GrantSystemAdminPayload{
			ZitadelUserID: testUserAliceID,
		},
	}))

	require.NoError(t, empSvc.RevokeSystemAdmin(ctxT(t), membership.RevokeSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeSystemAdminPayload{
			ZitadelUserID: testUserAliceID,
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.system_admins WHERE zitadel_user_id = ?`, testUserAliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)

	// Sync projector must have deleted the matching projection row.
	var projCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM projections.system_admins WHERE zitadel_user_id = ?`, testUserAliceID,
	).Scan(&projCount).Error)
	assert.Equal(t, int64(0), projCount)
}

func TestRevokeSystemAdmin_NotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.RevokeSystemAdmin(ctxT(t), membership.RevokeSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeSystemAdminPayload{
			ZitadelUserID: "nonexistent-user-id",
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeSystemAdminNotFound, oopsCode(t, err))
}

func TestSystemAdmin_IndependentOfEmployee(t *testing.T) {
	_ = takeFixture(t)
	// Grant SystemAdmin to a Zitadel user who is NOT hired as employee.
	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.GrantSystemAdminPayload{
			ZitadelUserID: testUserBobID,
		},
	}))
	// Verify no employee row was created.
	var empCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.employees WHERE zitadel_user_id = ?`, testUserBobID,
	).Scan(&empCount).Error)
	assert.Equal(t, int64(0), empCount)
}

func TestTerminateEmployee_DoesNotTouchSystemAdmin(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	// Grant Alice system admin. Alice is both a system admin and an employee.
	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.GrantSystemAdminPayload{
			ZitadelUserID: testUserAliceID,
		},
	}))

	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: aliceID.String()},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.system_admins WHERE zitadel_user_id = ?`, testUserAliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count, "SystemAdmin row must survive Terminate — they're independent")
}
