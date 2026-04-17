//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	systemadminv1 "github.com/medincident/medincident-command-service/pkg/event/system_admin/v1"
)

func TestGrantSystemAdmin_Success(t *testing.T) {
	_ = takeFixture(t)
	truncateOutbox(t)

	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		ZitadelUserID: testUserAliceID,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.system_admins WHERE zitadel_user_id = ?`, testUserAliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count)

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectSystemAdminGranted, rows[0].Subject)
	ev := &systemadminv1.SystemAdminGranted{}
	env := decodePayload(t, rows[0], ev)
	assert.Equal(t, membership.AggregateTypeSystemAdmin, env.AggregateType)
	assert.Equal(t, testUserAliceID, env.AggregateId)
}

func TestGrantSystemAdmin_ZitadelUserNotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		ZitadelUserID: "nonexistent-user-id",
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeZitadelUserNotFound, oopsCode(t, err))
}

func TestGrantSystemAdmin_AlreadyGranted(t *testing.T) {
	_ = takeFixture(t)
	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		ZitadelUserID: testUserAliceID,
	}))
	err := empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		ZitadelUserID: testUserAliceID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeSystemAdminAlreadyGranted, oopsCode(t, err))
}

func TestGrantSystemAdmin_EmptyInput(t *testing.T) {
	err := empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{ZitadelUserID: "   "})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeSystemAdminZitadelUserIDEmpty, oopsCode(t, err))
}

func TestRevokeSystemAdmin_Success(t *testing.T) {
	_ = takeFixture(t)
	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		ZitadelUserID: testUserAliceID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.RevokeSystemAdmin(ctxT(t), membership.RevokeSystemAdminCommand{
		ZitadelUserID: testUserAliceID,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.system_admins WHERE zitadel_user_id = ?`, testUserAliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectSystemAdminRevoked, rows[0].Subject)
}

func TestRevokeSystemAdmin_NotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.RevokeSystemAdmin(ctxT(t), membership.RevokeSystemAdminCommand{
		ZitadelUserID: "nonexistent-user-id",
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeSystemAdminNotFound, oopsCode(t, err))
}

func TestSystemAdmin_IndependentOfEmployee(t *testing.T) {
	_ = takeFixture(t)
	// Grant SystemAdmin to a Zitadel user who is NOT hired as employee.
	require.NoError(t, empSvc.GrantSystemAdmin(ctxT(t), membership.GrantSystemAdminCommand{
		ZitadelUserID: testUserBobID,
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
		ZitadelUserID: testUserAliceID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{ID: aliceID}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.system_admins WHERE zitadel_user_id = ?`, testUserAliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count, "SystemAdmin row must survive Terminate — they're independent")
}
