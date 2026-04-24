//go:build integration

package membership_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	memberread "github.com/medincident/medincident-backend/internal/service/query/membership"
)

// TestRoleReader_GetClinicHead returns nil (sentinel) when vacant, a
// view when assigned.
func TestRoleReader_GetClinicHead(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	_, clinicID, _ := seedOrgClinicDept(t, ctx, now)

	reader := memberread.NewRoleReader(testDB, authzSvc, &logger)

	view, err := reader.GetClinicHead(ctx, sysadminCaller, clinicID)
	require.ErrorIs(t, err, memberread.ErrRoleVacant)
	require.Nil(t, view)

	empID := uuid.Must(uuid.NewV7())
	head := &model.ClinicHead{
		ClinicID: clinicID, EmployeeID: empID,
		CreatedAt: now, UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.ClinicHeadAssigned(tx, head)
	}))

	view, err = reader.GetClinicHead(ctx, sysadminCaller, clinicID)
	require.NoError(t, err)
	require.NotNil(t, view)
	require.Equal(t, empID, view.EmployeeID)
}

// TestRoleReader_ListSystemAdmins returns every row, newest first.
func TestRoleReader_ListSystemAdmins(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.SystemAdminGranted(tx, "zit-a", now); err != nil {
			return err
		}
		return projector.SystemAdminGranted(tx, "zit-b", now.Add(time.Second))
	}))

	reader := memberread.NewRoleReader(testDB, authzSvc, &logger)
	items, err := reader.ListSystemAdmins(ctx, sysadminCaller, memberread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Equal(t, "zit-b", items[0].ZitadelUserID)
}
