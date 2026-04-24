//go:build integration

package incidentclassifier_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
)

func TestType_DisallowForPatients_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	_, _, tpID := makeIncidentTypeForAllow(t)

	// Pre-allow so we have something to revoke.
	_, err := typeSvc.AllowForPatients(ctx, classifiersvc.AllowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)
	require.True(t, loadType(t, tpID).IsAllowedForPatients)
	allowedUpdatedAt := loadType(t, tpID).UpdatedAt

	_, err = typeSvc.DisallowForPatients(ctx, classifiersvc.DisallowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DisallowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)

	row := loadType(t, tpID)
	assert.False(t, row.IsAllowedForPatients)
	assert.True(t, row.UpdatedAt.After(allowedUpdatedAt),
		"updated_at must advance after Disallow")

	var projAllowed bool
	require.NoError(t, testDB.Raw(
		`SELECT is_allowed_for_patients FROM projections.incident_types WHERE id = ?`,
		tpID,
	).Row().Scan(&projAllowed))
	assert.False(t, projAllowed)
}

func TestType_DisallowForPatients_Idempotent(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	_, _, tpID := makeIncidentTypeForAllow(t)

	// Type starts disallowed by default.
	first := loadType(t, tpID).UpdatedAt

	_, err := typeSvc.DisallowForPatients(ctx, classifiersvc.DisallowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DisallowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)

	row := loadType(t, tpID)
	assert.False(t, row.IsAllowedForPatients)
	assert.Equal(t, first, row.UpdatedAt,
		"idempotent Disallow must not advance updated_at")
}

func TestType_DisallowForPatients_InactiveTypeRejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	_, _, tpID := makeIncidentTypeForAllow(t)

	_, err := typeSvc.AllowForPatients(ctx, classifiersvc.AllowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)
	_, err = typeSvc.Deactivate(ctx, classifiersvc.DeactivateIncidentTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateIncidentTypePayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)

	_, err = typeSvc.DisallowForPatients(ctx, classifiersvc.DisallowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DisallowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentTypeInactive, codeOf(t, err))

	// Failed Disallow is a no-op: the prior Allow=true remains.
	assert.True(t, loadType(t, tpID).IsAllowedForPatients)
}

func TestType_DisallowForPatients_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	missing := uuid.New()
	_, err := typeSvc.DisallowForPatients(ctx, classifiersvc.DisallowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DisallowIncidentTypeForPatientsPayload{TypeID: missing.String()},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentTypeNotFound, codeOf(t, err))
}
