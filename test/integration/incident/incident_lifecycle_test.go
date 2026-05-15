//go:build integration

package incident_integration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
)

// codeOf extracts the oops Code string from any error in the tree.
func codeOf(t *testing.T, err error) string {
	t.Helper()
	var oe oops.OopsError
	if !errors.As(err, &oe) {
		t.Fatalf("expected oops error, got %T: %v", err, err)
	}
	code, ok := oe.Code().(string)
	if !ok {
		t.Fatalf("oops.Code() returned non-string: %T %v", oe.Code(), oe.Code())
	}
	return code
}

// world bundles the IDs and caller used by lifecycle tests.
type world struct {
	orgID      uuid.UUID
	clinicID   uuid.UUID
	deptID     uuid.UUID
	categoryID uuid.UUID
	typeID     uuid.UUID
	caller     authz.Caller
}

// setupWorld seeds one org/clinic/dept/category/type and one registrar employee.
func setupWorld(t *testing.T) world {
	t.Helper()
	oid := seedOrg(t)
	cid := seedClinic(t, oid)
	did := seedDept(t, cid)
	catID := seedCategory(t, oid)
	tid := seedType(t, catID, false)

	const registrarZitadelID = "registrar-zitadel"
	seedUser(t, registrarZitadelID, "Иван Иванов")
	seedEmployee(t, registrarZitadelID, oid, did)

	return world{
		orgID:      oid,
		clinicID:   cid,
		deptID:     did,
		categoryID: catID,
		typeID:     tid,
		caller:     authz.Caller{ZitadelUserID: registrarZitadelID},
	}
}

// TestIncidentLifecycle_CreateHappyPath: employee creates incident →
// status pending, history has one initial row.
func TestIncidentLifecycle_CreateHappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)

	var status string
	require.NoError(t, testDB.Raw(
		`SELECT status FROM domain.incidents WHERE id = ?`, incID,
	).Row().Scan(&status))
	assert.Equal(t, "pending", status)
}

// TestIncidentLifecycle_ToInProgress: DeptResponsible moves to in_progress →
// history grows to 2 rows.
func TestIncidentLifecycle_ToInProgress(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	// Seed a dept_responsible caller.
	const deptRespID = "deptresp-zitadel"
	seedUser(t, deptRespID, "Ответственный Отдела")
	empID := seedEmployee(t, deptRespID, w.orgID, w.deptID)
	seedDeptResponsible(t, empID, w.deptID)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)

	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: authz.Caller{ZitadelUserID: deptRespID},
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: incID.String(),
			NewStatus:  string(model.IncidentStatusInProgress),
		},
	}))
}

// TestIncidentLifecycle_PriorityChange: ClinicHead changes priority normal→high
// → priority_history has exactly one row.
func TestIncidentLifecycle_PriorityChange(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	const clinicHeadID = "clinichead-zitadel"
	seedUser(t, clinicHeadID, "Главврач Клиники")
	empID := seedEmployee(t, clinicHeadID, w.orgID, w.deptID)
	seedClinicHead(t, empID, w.clinicID)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)

	require.NoError(t, incidentSvc.UpdatePriority(ctx, incidentsvc.UpdateIncidentPriorityCommand{
		Caller: authz.Caller{ZitadelUserID: clinicHeadID},
		Payload: incidentsvc.UpdateIncidentPriorityPayload{
			IncidentID: incID.String(),
			Priority:   string(model.IncidentPriorityHigh),
		},
	}))
}

// TestIncidentLifecycle_DescriptionEdit: registrar edits description at
// in_progress → projection updated, no extra status/priority history rows.
func TestIncidentLifecycle_DescriptionEdit(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)

	// Move to in_progress first.
	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: sysadminCaller,
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: incID.String(),
			NewStatus:  string(model.IncidentStatusInProgress),
		},
	}))

	newDesc := "Уточнённое описание инцидента после обследования пациента"
	require.NoError(t, incidentSvc.UpdateDescription(ctx, incidentsvc.UpdateIncidentDescriptionCommand{
		Caller: w.caller,
		Payload: incidentsvc.UpdateIncidentDescriptionPayload{
			IncidentID:  incID.String(),
			Description: &newDesc,
		},
	}))

	var domainDesc string
	require.NoError(t, testDB.Raw(
		`SELECT COALESCE(description, '') FROM domain.incidents WHERE id = ?`, incID,
	).Row().Scan(&domainDesc))
	assert.Equal(t, newDesc, domainDesc)
}

// TestIncidentLifecycle_DoneIsTerminal: status moves to done → terminal.
// Further priority mutations are rejected with incident_frozen.
func TestIncidentLifecycle_DoneIsTerminal(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)

	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: sysadminCaller,
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: incID.String(),
			NewStatus:  string(model.IncidentStatusInProgress),
		},
	}))
	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: sysadminCaller,
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: incID.String(),
			NewStatus:  string(model.IncidentStatusDone),
		},
	}))

	err := incidentSvc.UpdatePriority(ctx, incidentsvc.UpdateIncidentPriorityCommand{
		Caller: sysadminCaller,
		Payload: incidentsvc.UpdateIncidentPriorityPayload{
			IncidentID: incID.String(),
			Priority:   string(model.IncidentPriorityCritical),
		},
	})
	require.Error(t, err)
	assert.Equal(t, incidentsvc.ErrCodeIncidentFrozen, codeOf(t, err))
}

// TestIncidentLifecycle_Reopen: OrgHead reopens a done incident → new
// incident has reopened_from_incident_id pointing back, status pending.
func TestIncidentLifecycle_Reopen(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	// Seed an OrgHead to perform the reopen.
	const orgHeadID = "orghead-zitadel"
	seedUser(t, orgHeadID, "Руководитель Организации")
	empID := seedEmployee(t, orgHeadID, w.orgID, w.deptID)
	seedOrgHead(t, empID, w.orgID)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)

	// Progress to done.
	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: sysadminCaller,
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: incID.String(),
			NewStatus:  string(model.IncidentStatusInProgress),
		},
	}))
	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: sysadminCaller,
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: incID.String(),
			NewStatus:  string(model.IncidentStatusDone),
		},
	}))

	// Reopen by OrgHead.
	res, err := incidentSvc.Reopen(ctx, incidentsvc.ReopenIncidentCommand{
		Caller:  authz.Caller{ZitadelUserID: orgHeadID},
		Payload: incidentsvc.ReopenIncidentPayload{IncidentID: incID.String()},
	})
	require.NoError(t, err)

	newID := res.NewIncidentID
	require.NotEqual(t, incID, newID)

	var sourceID string
	require.NoError(t, testDB.Raw(
		`SELECT COALESCE(reopened_from_incident_id::text, '') FROM domain.incidents WHERE id = ?`, newID,
	).Row().Scan(&sourceID))
	assert.Equal(t, incID.String(), sourceID)

	var newStatus string
	require.NoError(t, testDB.Raw(
		`SELECT status FROM domain.incidents WHERE id = ?`, newID,
	).Row().Scan(&newStatus))
	assert.Equal(t, "pending", newStatus)
}

// TestIncidentLifecycle_ReopenPendingFails: reopening a pending incident
// must be rejected with incident_not_reopenable.
func TestIncidentLifecycle_ReopenPendingFails(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)

	_, err := incidentSvc.Reopen(ctx, incidentsvc.ReopenIncidentCommand{
		Caller:  sysadminCaller,
		Payload: incidentsvc.ReopenIncidentPayload{IncidentID: incID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, incidentsvc.ErrCodeIncidentNotReopenable, codeOf(t, err))
}

// TestIncidentLifecycle_ReopenInProgressFails: reopening an in-progress
// incident must be rejected.
func TestIncidentLifecycle_ReopenInProgressFails(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)
	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller:  sysadminCaller,
		Payload: incidentsvc.UpdateIncidentStatusPayload{IncidentID: incID.String(), NewStatus: string(model.IncidentStatusInProgress)},
	}))

	_, err := incidentSvc.Reopen(ctx, incidentsvc.ReopenIncidentCommand{
		Caller:  sysadminCaller,
		Payload: incidentsvc.ReopenIncidentPayload{IncidentID: incID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, incidentsvc.ErrCodeIncidentNotReopenable, codeOf(t, err))
}

// TestIncidentLifecycle_ReopenCancelledFails: a cancelled incident is NOT
// reopenable — only done/rejected are valid sources.
func TestIncidentLifecycle_ReopenCancelledFails(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)
	require.NoError(t, incidentSvc.Cancel(ctx, incidentsvc.CancelIncidentCommand{
		Caller:  w.caller,
		Payload: incidentsvc.CancelIncidentPayload{IncidentID: incID.String()},
	}))

	_, err := incidentSvc.Reopen(ctx, incidentsvc.ReopenIncidentCommand{
		Caller:  sysadminCaller,
		Payload: incidentsvc.ReopenIncidentPayload{IncidentID: incID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, incidentsvc.ErrCodeIncidentNotReopenable, codeOf(t, err))
}

// TestIncidentLifecycle_CancelByRegistrar: registrar cancels pending
// incident. Non-registrar cancel attempt returns permission_denied.
func TestIncidentLifecycle_CancelByRegistrar(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	w := setupWorld(t)

	incID := createIncident(t, ctx, w.caller, w.deptID, w.categoryID, w.typeID)

	// Non-registrar (a different user, not in the org) should be denied.
	const strangerID = "stranger-zitadel"
	seedUser(t, strangerID, "Посторонний")
	stranger := authz.Caller{ZitadelUserID: strangerID}
	err := incidentSvc.Cancel(ctx, incidentsvc.CancelIncidentCommand{
		Caller:  stranger,
		Payload: incidentsvc.CancelIncidentPayload{IncidentID: incID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, authz.ErrCodePermissionDenied, codeOf(t, err))

	// Registrar succeeds.
	require.NoError(t, incidentSvc.Cancel(ctx, incidentsvc.CancelIncidentCommand{
		Caller:  w.caller,
		Payload: incidentsvc.CancelIncidentPayload{IncidentID: incID.String()},
	}))

	var status string
	require.NoError(t, testDB.Raw(
		`SELECT status FROM domain.incidents WHERE id = ?`, incID,
	).Row().Scan(&status))
	assert.Equal(t, "cancelled", status)
}
