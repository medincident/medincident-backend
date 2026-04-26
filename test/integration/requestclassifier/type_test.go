//go:build integration

package requestclassifier_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
)

func TestRequestType_Create_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{
			OrganizationID: orgID.String(),
			Name:           "Тип заявки",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	row := loadType(t, res.ID)
	assert.Equal(t, "Тип заявки", row.Name)
	assert.True(t, row.IsActive)
	assert.Equal(t, orgID, row.OrganizationID)

	var projName string
	var projActive bool
	require.NoError(t, testDB.Raw(
		`SELECT name, is_active FROM projections.request_types WHERE id = ?`, res.ID,
	).Row().Scan(&projName, &projActive))
	assert.Equal(t, "Тип заявки", projName)
	assert.True(t, projActive)
}

func TestRequestType_Create_NameConflict(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	_, err := typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Дубль"},
	})
	require.NoError(t, err)

	_, err = typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Дубль"},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeRequestTypeNameConflict, codeOf(t, err))
}

func TestRequestType_UpdateDetails_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Старое имя"},
	})
	require.NoError(t, err)

	desc := "Описание заявки длиннее восьми символов"
	_, err = typeSvc.UpdateDetails(ctx, classifiersvc.UpdateRequestTypeDetailsCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.UpdateRequestTypeDetailsPayload{
			TypeID:      res.ID.String(),
			Name:        "Новое имя",
			Description: &desc,
		},
	})
	require.NoError(t, err)

	row := loadType(t, res.ID)
	assert.Equal(t, "Новое имя", row.Name)
	assert.True(t, row.Description.Valid)

	var projName string
	require.NoError(t, testDB.Raw(
		`SELECT name FROM projections.request_types WHERE id = ?`, res.ID,
	).Scan(&projName).Error)
	assert.Equal(t, "Новое имя", projName)
}

func TestRequestType_UpdateDetails_NotFound(t *testing.T) {
	resetDB(t)
	_, err := typeSvc.UpdateDetails(context.Background(), classifiersvc.UpdateRequestTypeDetailsCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.UpdateRequestTypeDetailsPayload{TypeID: uuid.New().String(), Name: "Тест"},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeRequestTypeNotFound, codeOf(t, err))
}

func TestRequestType_Deactivate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Тип"},
	})
	require.NoError(t, err)

	_, err = typeSvc.Deactivate(ctx, classifiersvc.DeactivateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateRequestTypePayload{TypeID: res.ID.String()},
	})
	require.NoError(t, err)
	assert.False(t, loadType(t, res.ID).IsActive)

	var projActive bool
	require.NoError(t, testDB.Raw(
		`SELECT is_active FROM projections.request_types WHERE id = ?`, res.ID,
	).Scan(&projActive).Error)
	assert.False(t, projActive)
}

func TestRequestType_Deactivate_Idempotent(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Тип"},
	})
	require.NoError(t, err)

	_, err = typeSvc.Deactivate(ctx, classifiersvc.DeactivateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateRequestTypePayload{TypeID: res.ID.String()},
	})
	require.NoError(t, err)

	_, err = typeSvc.Deactivate(ctx, classifiersvc.DeactivateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateRequestTypePayload{TypeID: res.ID.String()},
	})
	require.NoError(t, err)
}

func TestRequestType_Reactivate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Тип"},
	})
	require.NoError(t, err)

	_, _ = typeSvc.Deactivate(ctx, classifiersvc.DeactivateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateRequestTypePayload{TypeID: res.ID.String()},
	})

	_, err = typeSvc.Reactivate(ctx, classifiersvc.ReactivateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.ReactivateRequestTypePayload{TypeID: res.ID.String()},
	})
	require.NoError(t, err)
	assert.True(t, loadType(t, res.ID).IsActive)

	var projActive bool
	require.NoError(t, testDB.Raw(
		`SELECT is_active FROM projections.request_types WHERE id = ?`, res.ID,
	).Scan(&projActive).Error)
	assert.True(t, projActive)
}

func TestRequestType_Reactivate_NameConflict(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	r1, err := typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Дубль"},
	})
	require.NoError(t, err)

	_, _ = typeSvc.Deactivate(ctx, classifiersvc.DeactivateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateRequestTypePayload{TypeID: r1.ID.String()},
	})

	// Another active type with the same name.
	_, err = typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Дубль"},
	})
	require.NoError(t, err)

	_, err = typeSvc.Reactivate(ctx, classifiersvc.ReactivateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.ReactivateRequestTypePayload{TypeID: r1.ID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeRequestTypeNameConflict, codeOf(t, err))
}

func TestRequestType_Delete_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Тип"},
	})
	require.NoError(t, err)

	_, err = typeSvc.Delete(ctx, classifiersvc.DeleteRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeleteRequestTypePayload{TypeID: res.ID.String()},
	})
	require.NoError(t, err)

	var n int64
	testDB.Raw(`SELECT COUNT(*) FROM domain.request_types WHERE id = ?`, res.ID).Scan(&n)
	assert.Equal(t, int64(0), n)

	testDB.Raw(`SELECT COUNT(*) FROM projections.request_types WHERE id = ?`, res.ID).Scan(&n)
	assert.Equal(t, int64(0), n)
}
