//go:build integration

package incident_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	buffercmd "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
)

// bufferWorld bundles the setup needed by buffer flow tests.
type bufferWorld struct {
	orgID      uuid.UUID
	deptID     uuid.UUID
	categoryID uuid.UUID
	typeID     uuid.UUID
	// dispatcher is an employee with OrgDispatcher role.
	dispatcher   authz.Caller
	dispatcherID uuid.UUID
}

// setupBufferWorld seeds org/clinic/dept/category/type, a patient-allowed type,
// and a dispatcher employee.
func setupBufferWorld(t *testing.T) bufferWorld {
	t.Helper()
	oid := seedOrg(t)
	cid := seedClinic(t, oid)
	did := seedDept(t, cid)
	catID := seedCategory(t, oid)
	// Patient-allowed type for buffer submissions.
	tid := seedType(t, catID, true)

	const dispatcherZitadelID = "dispatcher-zitadel"
	seedUser(t, dispatcherZitadelID, "Диспетчер Организации")
	dispEmpID := seedEmployee(t, dispatcherZitadelID, oid, did)
	seedOrgDispatcher(t, dispEmpID, oid)

	return bufferWorld{
		orgID:        oid,
		deptID:       did,
		categoryID:   catID,
		typeID:       tid,
		dispatcher:   authz.Caller{ZitadelUserID: dispatcherZitadelID},
		dispatcherID: dispEmpID,
	}
}

// TestBufferFlow_SubmitAndRead: patient submits buffer entry → reader returns it.
func TestBufferFlow_SubmitAndRead(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	bw := setupBufferWorld(t)

	const patientID = "patient-zitadel"
	seedUser(t, patientID, "Пациент Тестовый")
	patient := authz.Caller{ZitadelUserID: patientID}

	desc := "Болит голова"
	res, err := bufferSvc.Submit(ctx, buffercmd.SubmitCommand{
		Caller: patient,
		Payload: buffercmd.SubmitPayload{
			OrganizationID: bw.orgID.String(),
			Description:    &desc,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	// Reader returns it.
	view, err := bufferRdr.GetBufferEntry(ctx, patientID, res.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BufferStatusPending, view.Status)
	assert.True(t, view.Description.Valid)
	assert.Equal(t, "Болит голова", view.Description.String)
}

// TestBufferFlow_UpdateDescription: patient updates description → projection updated.
func TestBufferFlow_UpdateDescription(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	bw := setupBufferWorld(t)

	const patientID = "patient-zitadel"
	seedUser(t, patientID, "Пациент Тестовый")
	patient := authz.Caller{ZitadelUserID: patientID}

	desc := "Первичная жалоба пациента"
	res, err := bufferSvc.Submit(ctx, buffercmd.SubmitCommand{
		Caller:  patient,
		Payload: buffercmd.SubmitPayload{OrganizationID: bw.orgID.String(), Description: &desc},
	})
	require.NoError(t, err)

	updatedDesc := "Уточнённая жалоба пациента после консультации"
	require.NoError(t, bufferSvc.Update(ctx, buffercmd.UpdateCommand{
		Caller: patient,
		Payload: buffercmd.UpdatePayload{
			BufferID:    res.ID.String(),
			Description: &updatedDesc,
		},
	}))

	view, err := bufferRdr.GetBufferEntry(ctx, patientID, res.ID)
	require.NoError(t, err)
	assert.Equal(t, updatedDesc, view.Description.String)
}

// TestBufferFlow_PatientCancel: patient cancels → status flips to cancelled.
func TestBufferFlow_PatientCancel(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	bw := setupBufferWorld(t)

	const patientID = "patient-cancel-zitadel"
	seedUser(t, patientID, "Отменяющий Пациент")
	patient := authz.Caller{ZitadelUserID: patientID}

	res, err := bufferSvc.Submit(ctx, buffercmd.SubmitCommand{
		Caller:  patient,
		Payload: buffercmd.SubmitPayload{OrganizationID: bw.orgID.String()},
	})
	require.NoError(t, err)

	require.NoError(t, bufferSvc.Cancel(ctx, buffercmd.CancelCommand{
		Caller:  patient,
		Payload: buffercmd.CancelPayload{BufferID: res.ID.String()},
	}))

	view, err := bufferRdr.GetBufferEntry(ctx, patientID, res.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BufferStatusCancelled, view.Status)
}

// TestBufferFlow_PublishCreatesIncident: dispatcher publishes buffer entry →
// incident created with source_buffer_id and source_patient_zitadel_user_id.
// Buffer row goes to published with published_incident_id set.
// Patient GetBufferEntry shows PatientPerspective.
func TestBufferFlow_PublishCreatesIncident(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	bw := setupBufferWorld(t)

	const patientID = "patient-pub-zitadel"
	seedUser(t, patientID, "Публикуемый Пациент")
	patient := authz.Caller{ZitadelUserID: patientID}

	occAt := time.Now().Add(-2 * time.Hour).Format(time.RFC3339Nano)
	desc := "Жалоба пациента на боли в груди"
	res, err := bufferSvc.Submit(ctx, buffercmd.SubmitCommand{
		Caller: patient,
		Payload: buffercmd.SubmitPayload{
			OrganizationID: bw.orgID.String(),
			Description:    &desc,
			OccurredAt:     &occAt,
		},
	})
	require.NoError(t, err)

	publishRes, err := bufferSvc.Publish(ctx, &buffercmd.PublishCommand{
		Caller: bw.dispatcher,
		Payload: buffercmd.PublishPayload{
			BufferID:     res.ID.String(),
			DepartmentID: bw.deptID.String(),
			CategoryID:   bw.categoryID.String(),
			TypeID:       bw.typeID.String(),
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, publishRes.IncidentID)

	// Buffer row is published with published_incident_id.
	bufView, err := bufferRdr.GetBufferEntry(ctx, patientID, res.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BufferStatusPublished, bufView.Status)
	require.True(t, bufView.PublishedIncidentID.Valid)
	assert.Equal(t, publishRes.IncidentID, bufView.PublishedIncidentID.UUID)

	// Incident has source_buffer_id and source_patient_zitadel_user_id.
	var sourceBufferID uuid.NullUUID
	var sourcePatient string
	require.NoError(t, testDB.Raw(
		`SELECT source_buffer_id, COALESCE(source_patient_zitadel_user_id, '')
		 FROM domain.incidents WHERE id = ?`, publishRes.IncidentID,
	).Row().Scan(&sourceBufferID, &sourcePatient))
	require.True(t, sourceBufferID.Valid)
	assert.Equal(t, res.ID, sourceBufferID.UUID)
	assert.Equal(t, patientID, sourcePatient)
}

// TestBufferFlow_PublishThenCloseIsVisibleToPatient: when dispatcher closes
// the incident as done, the patient's view of the incident should show
// status done (via projection).
func TestBufferFlow_PublishThenCloseIsVisibleToPatient(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	bw := setupBufferWorld(t)

	const patientID = "patient-close-zitadel"
	seedUser(t, patientID, "Наблюдающий Пациент")
	patient := authz.Caller{ZitadelUserID: patientID}

	res, err := bufferSvc.Submit(ctx, buffercmd.SubmitCommand{
		Caller:  patient,
		Payload: buffercmd.SubmitPayload{OrganizationID: bw.orgID.String()},
	})
	require.NoError(t, err)

	publishRes, err := bufferSvc.Publish(ctx, &buffercmd.PublishCommand{
		Caller: bw.dispatcher,
		Payload: buffercmd.PublishPayload{
			BufferID:     res.ID.String(),
			DepartmentID: bw.deptID.String(),
			CategoryID:   bw.categoryID.String(),
			TypeID:       bw.typeID.String(),
		},
	})
	require.NoError(t, err)

	// Dispatcher moves incident to done.
	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: bw.dispatcher,
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: publishRes.IncidentID.String(),
			NewStatus:  "INCIDENT_STATUS_IN_PROGRESS",
		},
	}))
	require.NoError(t, incidentSvc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: bw.dispatcher,
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: publishRes.IncidentID.String(),
			NewStatus:  "INCIDENT_STATUS_DONE",
		},
	}))

	// Patient can read the incident via its source_patient_zitadel_user_id link.
	incView, err := incidentRdr.GetIncident(ctx, patientID, publishRes.IncidentID)
	require.NoError(t, err)
	assert.Equal(t, model.IncidentStatusDone, incView.Status)
	assert.True(t, incView.PatientPerspective)
}

// TestBufferFlow_RejectSetsStatusRejected: dispatcher rejects submission →
// buffer status is rejected; patient can still read it.
func TestBufferFlow_RejectSetsStatusRejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	bw := setupBufferWorld(t)

	const patientID = "patient-reject-zitadel"
	seedUser(t, patientID, "Отклонённый Пациент")
	patient := authz.Caller{ZitadelUserID: patientID}

	res, err := bufferSvc.Submit(ctx, buffercmd.SubmitCommand{
		Caller:  patient,
		Payload: buffercmd.SubmitPayload{OrganizationID: bw.orgID.String()},
	})
	require.NoError(t, err)

	require.NoError(t, bufferSvc.Reject(ctx, buffercmd.RejectCommand{
		Caller:  bw.dispatcher,
		Payload: buffercmd.RejectPayload{BufferID: res.ID.String()},
	}))

	view, err := bufferRdr.GetBufferEntry(ctx, patientID, res.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BufferStatusRejected, view.Status)
	// No incident created.
	assert.False(t, view.PublishedIncidentID.Valid)
}
