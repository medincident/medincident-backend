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

// makeIncidentTypeForAllow creates an active incident type owned by a fresh
// organization and returns its IDs. Used by every test in this file to keep
// setup short and focused on the AllowForPatients behaviour under test.
func makeIncidentTypeForAllow(t *testing.T) (orgID, catID, typeID uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	orgID = insertOrganization(t, "Org")
	catRes, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Cat",
		},
	})
	require.NoError(t, err)
	typeRes, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{
			CategoryID: catRes.ID.String(),
			Name:       "Ty",
		},
	})
	require.NoError(t, err)
	return orgID, catRes.ID, typeRes.ID
}

func TestType_AllowForPatients_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	_, _, tpID := makeIncidentTypeForAllow(t)

	// Pre-condition: the freshly-created type is not allowed for patients.
	row := loadType(t, tpID)
	assert.False(t, row.IsAllowedForPatients)
	createdUpdatedAt := row.UpdatedAt

	_, err := typeSvc.AllowForPatients(ctx, classifiersvc.AllowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)

	row = loadType(t, tpID)
	assert.True(t, row.IsAllowedForPatients)
	assert.True(t, row.UpdatedAt.After(createdUpdatedAt),
		"updated_at must advance after Allow")

	// Projection mirrors the new state.
	var projAllowed bool
	require.NoError(t, testDB.Raw(
		`SELECT is_allowed_for_patients FROM projections.incident_types WHERE id = ?`,
		tpID,
	).Row().Scan(&projAllowed))
	assert.True(t, projAllowed)
}

func TestType_AllowForPatients_Idempotent(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	_, _, tpID := makeIncidentTypeForAllow(t)

	_, err := typeSvc.AllowForPatients(ctx, classifiersvc.AllowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)
	first := loadType(t, tpID).UpdatedAt

	// Second Allow on already-allowed type: must succeed and NOT advance updated_at.
	_, err = typeSvc.AllowForPatients(ctx, classifiersvc.AllowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)

	row := loadType(t, tpID)
	assert.True(t, row.IsAllowedForPatients)
	assert.Equal(t, first, row.UpdatedAt,
		"idempotent Allow must not advance updated_at")
}

func TestType_AllowForPatients_InactiveTypeRejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	_, _, tpID := makeIncidentTypeForAllow(t)

	_, err := typeSvc.Deactivate(ctx, classifiersvc.DeactivateIncidentTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateIncidentTypePayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)

	_, err = typeSvc.AllowForPatients(ctx, classifiersvc.AllowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentTypeInactive, codeOf(t, err))

	assert.False(t, loadType(t, tpID).IsAllowedForPatients)
}

func TestType_AllowForPatients_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	missing := uuid.New()
	_, err := typeSvc.AllowForPatients(ctx, classifiersvc.AllowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: missing.String()},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentTypeNotFound, codeOf(t, err))
}

// Allow → Deactivate → Reactivate must leave is_allowed_for_patients = TRUE.
// Deactivation does not silently clear administrator intent.
func TestType_AllowForPatients_PersistsAcrossDeactivateReactivate(t *testing.T) {
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
	assert.True(t, loadType(t, tpID).IsAllowedForPatients,
		"Deactivate must not clear is_allowed_for_patients")

	_, err = typeSvc.Reactivate(ctx, classifiersvc.ReactivateIncidentTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.ReactivateIncidentTypePayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)
	assert.True(t, loadType(t, tpID).IsAllowedForPatients,
		"Reactivate must leave is_allowed_for_patients untouched")
}

// Cascading deactivation through the parent category MUST leave
// is_allowed_for_patients untouched on the affected types.
func TestType_AllowForPatients_PersistsAcrossCategoryCascade(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	_, catID, tpID := makeIncidentTypeForAllow(t)

	_, err := typeSvc.AllowForPatients(ctx, classifiersvc.AllowIncidentTypeForPatientsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: tpID.String()},
	})
	require.NoError(t, err)

	_, err = categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateIncidentCategoryPayload{CategoryID: catID.String()},
	})
	require.NoError(t, err)

	row := loadType(t, tpID)
	assert.False(t, row.IsActive,
		"category cascade must deactivate the type")
	assert.True(t, row.IsAllowedForPatients,
		"category cascade must NOT clear is_allowed_for_patients on the type")
}
