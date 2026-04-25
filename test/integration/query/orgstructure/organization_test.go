//go:build integration

package orgstructure_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	orgread "github.com/medincident/medincident-backend/internal/service/query/orgstructure"
)

// TestOrganizationReader_Get_ReturnsViewSeededByProjector seeds a
// projection row via the command-side projector and reads it back
// through OrganizationReader. The read surface preserves the full card.
func TestOrganizationReader_Get_ReturnsViewSeededByProjector(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	orgID := uuid.Must(uuid.NewV7())
	now := time.Now().UTC().Truncate(time.Second)
	org := &model.Organization{
		ID:           orgID,
		Name:         "Acme Clinics",
		Description:  null.StringFrom("a description above min"),
		LegalAddress: model.Address{Text: "1 Real Street", Point: null.ValueFrom(model.Point{Longitude: 4.25, Latitude: 50.85})},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrganizationCreated(tx, org)
	}))

	reader := orgread.NewOrganizationReader(testDB, &logger)
	view, err := reader.Get(ctx, orgID)
	require.NoError(t, err)
	require.Equal(t, orgID, view.ID)
	require.Equal(t, "Acme Clinics", view.Name)
	require.NotNil(t, view.Description)
	require.Equal(t, "a description above min", *view.Description)
	require.Equal(t, "1 Real Street", view.LegalAddress.Text)
	require.NotNil(t, view.LegalAddress.Point)
	require.InDelta(t, 4.25, view.LegalAddress.Point.Longitude, 0.0001)
	require.InDelta(t, 50.85, view.LegalAddress.Point.Latitude, 0.0001)
}

// TestOrganizationReader_Get_NotFound returns the expected oops code.
func TestOrganizationReader_Get_NotFound(t *testing.T) {
	resetProjections(t)
	logger := zerolog.Nop()

	reader := orgread.NewOrganizationReader(testDB, &logger)
	_, err := reader.Get(context.Background(), uuid.Must(uuid.NewV7()))
	require.Error(t, err)
}

// TestOrganizationReader_List_And_Count returns seeded rows in
// reverse-chronological order and the right total.
func TestOrganizationReader_List_And_Count(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	ids := []uuid.UUID{
		uuid.Must(uuid.NewV7()),
		uuid.Must(uuid.NewV7()),
		uuid.Must(uuid.NewV7()),
	}
	for i, id := range ids {
		org := &model.Organization{
			ID:           id,
			Name:         "Org " + id.String()[:8],
			LegalAddress: model.Address{Text: "Street"},
			CreatedAt:    now.Add(time.Duration(i) * time.Minute),
			UpdatedAt:    now.Add(time.Duration(i) * time.Minute),
		}
		require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return projector.OrganizationCreated(tx, org)
		}))
	}

	reader := orgread.NewOrganizationReader(testDB, &logger)
	items, err := reader.List(ctx, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 3)
	// Most-recently created first: reverse of insertion order.
	require.Equal(t, ids[2], items[0].ID)
	require.Equal(t, ids[0], items[2].ID)

	total, err := reader.Count(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
}

// TestOrganizationReader_Search_FiltersByNameSubstring seeds three
// organizations and asserts ILIKE %acme% returns both Acme rows.
func TestOrganizationReader_Search_FiltersByNameSubstring(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	names := []string{"Acme Clinics", "Brightside Health", "Acme Medical"}
	ids := make([]uuid.UUID, len(names))
	for i, name := range names {
		ids[i] = uuid.Must(uuid.NewV7())
		org := &model.Organization{
			ID:           ids[i],
			Name:         name,
			LegalAddress: model.Address{Text: "Street"},
			CreatedAt:    now.Add(time.Duration(i) * time.Minute),
			UpdatedAt:    now.Add(time.Duration(i) * time.Minute),
		}
		require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return projector.OrganizationCreated(tx, org)
		}))
	}

	reader := orgread.NewOrganizationReader(testDB, &logger)
	items, err := reader.Search(ctx, "acme", orgread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	got := map[uuid.UUID]string{items[0].ID: items[0].Name, items[1].ID: items[1].Name}
	require.Contains(t, got, ids[0])
	require.Contains(t, got, ids[2])
}

// TestClinicReader_Get_And_ListByOrganization covers the clinic reader.
func TestClinicReader_Get_And_ListByOrganization(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	orgID := uuid.Must(uuid.NewV7())
	now := time.Now().UTC().Truncate(time.Second)
	org := &model.Organization{
		ID: orgID, Name: "O", LegalAddress: model.Address{Text: "S"}, CreatedAt: now, UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrganizationCreated(tx, org)
	}))

	clinicID := uuid.Must(uuid.NewV7())
	clinic := &model.Clinic{
		ID:              clinicID,
		OrganizationID:  orgID,
		Name:            "Downtown Clinic",
		Description:     null.StringFrom("downtown branch, long enough description"),
		PhysicalAddress: model.Address{Text: "1 City Plaza"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.ClinicCreated(tx, clinic)
	}))

	reader := orgread.NewClinicReader(testDB, authzSvc, &logger)
	got, err := reader.Get(ctx, sysadminCaller, clinicID)
	require.NoError(t, err)
	require.Equal(t, clinicID, got.ID)
	require.Equal(t, "Downtown Clinic", got.Name)
	require.Equal(t, "1 City Plaza", got.PhysicalAddress.Text)

	list, err := reader.ListByOrganization(ctx, sysadminCaller, orgID, orgread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, clinicID, list[0].ID)
}

// TestDepartmentReader_Get_And_ListByClinic covers the department reader.
func TestDepartmentReader_Get_And_ListByClinic(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	now := time.Now().UTC().Truncate(time.Second)

	org := &model.Organization{ID: orgID, Name: "O", LegalAddress: model.Address{Text: "S"}, CreatedAt: now, UpdatedAt: now}
	clinic := &model.Clinic{ID: clinicID, OrganizationID: orgID, Name: "C", PhysicalAddress: model.Address{Text: "S"}, CreatedAt: now, UpdatedAt: now}
	dept := &model.Department{ID: deptID, ClinicID: clinicID, Name: "Radiology", Description: null.StringFrom("radiology department, long description"), CreatedAt: now, UpdatedAt: now}

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.OrganizationCreated(tx, org); err != nil {
			return err
		}
		if err := projector.ClinicCreated(tx, clinic); err != nil {
			return err
		}
		return projector.DepartmentCreated(tx, dept)
	}))

	reader := orgread.NewDepartmentReader(testDB, authzSvc, &logger)
	got, err := reader.Get(ctx, sysadminCaller, deptID)
	require.NoError(t, err)
	require.Equal(t, "Radiology", got.Name)

	list, err := reader.ListByClinic(ctx, sysadminCaller, clinicID, orgread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, list, 1)
}
