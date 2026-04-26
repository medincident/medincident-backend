//go:build integration

package announcement_query_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	announcementread "github.com/medincident/medincident-backend/internal/service/query/announcement"
)

func TestGetAnnouncement_HappyPath(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	aID := insertAnnouncement(t, orgID, uuid.NullUUID{}, uuid.NullUUID{}, "Главное объявление")

	got, err := reader.GetAnnouncement(ctx, sysadminZitadelID, aID)
	require.NoError(t, err)
	assert.Equal(t, "Главное объявление", got.Title)
	assert.Equal(t, int64(1), got.ViewCount)
}

func TestGetAnnouncement_NotFound(t *testing.T) {
	resetProjections(t)
	_, err := reader.GetAnnouncement(context.Background(), sysadminZitadelID, uuid.New())
	require.Error(t, err)
}

func TestGetAnnouncement_AdminSeesArchived(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	aID := insertAnnouncement(t, orgID, uuid.NullUUID{}, uuid.NullUUID{}, "Архивное")

	require.NoError(t, testDB.Exec(
		`UPDATE domain.announcements SET is_archived=true WHERE id=?`, aID,
	).Error)

	got, err := reader.GetAnnouncement(ctx, sysadminZitadelID, aID)
	require.NoError(t, err)
	assert.True(t, got.IsArchived)
}

func TestListForOrganization_BasicList(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	insertAnnouncement(t, orgID, uuid.NullUUID{}, uuid.NullUUID{}, "Объявление 1")
	insertAnnouncement(t, orgID, uuid.NullUUID{}, uuid.NullUUID{}, "Объявление 2")

	res, err := reader.ListForOrganization(ctx, sysadminZitadelID, orgID, announcementread.ListFilter{Limit: 10})
	require.NoError(t, err)
	assert.Len(t, res.Items, 2)
	assert.Nil(t, res.NextCursor)
}

func TestListForOrganization_EmployeeSeesOwnOrg(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	clinicID := insertClinic(t, orgID)
	deptID := insertDept(t, clinicID)
	insertAnnouncement(t, orgID, uuid.NullUUID{}, uuid.NullUUID{}, "Видимое объявление")
	insertEmployee(t, "emp1", orgID, deptID)

	res, err := reader.ListForOrganization(ctx, "emp1", orgID, announcementread.ListFilter{Limit: 10})
	require.NoError(t, err)
	assert.Len(t, res.Items, 1)
}

func TestListForOrganization_EmployeeWrongOrg(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	org1 := insertOrg(t)
	org2 := insertOrg(t)
	clinic2 := insertClinic(t, org2)
	dept2 := insertDept(t, clinic2)
	insertAnnouncement(t, org1, uuid.NullUUID{}, uuid.NullUUID{}, "Объявление орг 1")
	insertEmployee(t, "emp_other", org2, dept2)

	res, err := reader.ListForOrganization(ctx, "emp_other", org1, announcementread.ListFilter{Limit: 10})
	require.NoError(t, err)
	assert.Empty(t, res.Items)
}

func TestListForClinic_SeesOrgAndClinicLevel(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	clinicID := insertClinic(t, orgID)
	deptID := insertDept(t, clinicID)
	insertAnnouncement(t, orgID, uuid.NullUUID{}, uuid.NullUUID{}, "Орг-уровень")
	insertAnnouncement(t, orgID, uuid.NullUUID{UUID: clinicID, Valid: true}, uuid.NullUUID{}, "Клиника-уровень")
	insertEmployee(t, "clinic_emp", orgID, deptID)

	res, err := reader.ListForClinic(ctx, "clinic_emp", clinicID, announcementread.ListFilter{Limit: 10})
	require.NoError(t, err)
	assert.Len(t, res.Items, 2)
}

func TestListForDepartment_SeesAllLevels(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	clinicID := insertClinic(t, orgID)
	deptID := insertDept(t, clinicID)
	insertAnnouncement(t, orgID, uuid.NullUUID{}, uuid.NullUUID{}, "Орг-уровень")
	insertAnnouncement(t, orgID, uuid.NullUUID{UUID: clinicID, Valid: true}, uuid.NullUUID{}, "Клиника-уровень")
	insertAnnouncement(t, orgID,
		uuid.NullUUID{UUID: clinicID, Valid: true},
		uuid.NullUUID{UUID: deptID, Valid: true},
		"Отделение-уровень",
	)
	insertEmployee(t, "dept_emp", orgID, deptID)

	res, err := reader.ListForDepartment(ctx, "dept_emp", deptID, announcementread.ListFilter{Limit: 10})
	require.NoError(t, err)
	assert.Len(t, res.Items, 3)
}
