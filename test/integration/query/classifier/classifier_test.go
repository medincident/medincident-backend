//go:build integration

package classifier_query_integration_test

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
	classifierread "github.com/medincident/medincident-backend/internal/service/query/incident/classifier"
)

// TestReader_Category_Get_And_Subtree seeds a three-level tree and
// verifies Get/ListActiveRootCategories/ListCategorySubtree.
func TestReader_Category_Get_And_Subtree(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())
	rootID := uuid.Must(uuid.NewV7())
	childID := uuid.Must(uuid.NewV7())
	grandID := uuid.Must(uuid.NewV7())

	root := &model.IncidentCategory{
		ID: rootID, OrganizationID: orgID,
		Name: "Safety", Description: null.StringFrom("root category"), IsActive: true,
		CreatedAt: now, UpdatedAt: now,
	}
	child := &model.IncidentCategory{
		ID: childID, OrganizationID: orgID,
		ParentCategoryID: uuid.NullUUID{UUID: rootID, Valid: true},
		Name:             "Falls",
		IsActive:         true,
		CreatedAt:        now.Add(time.Second),
		UpdatedAt:        now.Add(time.Second),
	}
	grand := &model.IncidentCategory{
		ID: grandID, OrganizationID: orgID,
		ParentCategoryID: uuid.NullUUID{UUID: childID, Valid: true},
		Name:             "In Ward",
		IsActive:         true,
		CreatedAt:        now.Add(2 * time.Second),
		UpdatedAt:        now.Add(2 * time.Second),
	}

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.CategoryCreated(tx, root); err != nil {
			return err
		}
		if err := projector.CategoryCreated(tx, child); err != nil {
			return err
		}
		return projector.CategoryCreated(tx, grand)
	}))

	reader := classifierread.NewReader(testDB, &logger)
	got, err := reader.GetCategory(ctx, childID)
	require.NoError(t, err)
	require.Equal(t, "Falls", got.Name)
	require.NotNil(t, got.ParentCategoryID)
	require.Equal(t, rootID, *got.ParentCategoryID)

	roots, err := reader.ListActiveRootCategories(ctx, orgID)
	require.NoError(t, err)
	require.Len(t, roots, 1)
	require.Equal(t, rootID, roots[0].ID)

	subtree, err := reader.ListCategorySubtree(ctx, rootID)
	require.NoError(t, err)
	require.Len(t, subtree, 3)
	// Ordered by created_at ASC.
	require.Equal(t, rootID, subtree[0].ID)
	require.Equal(t, childID, subtree[1].ID)
	require.Equal(t, grandID, subtree[2].ID)
}

// TestReader_Type_Get_And_ListActiveTypesByOrganization covers the
// type-reader paths.
func TestReader_Type_Get_And_ListActiveTypesByOrganization(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())
	catID := uuid.Must(uuid.NewV7())
	typeID := uuid.Must(uuid.NewV7())

	cat := &model.IncidentCategory{
		ID: catID, OrganizationID: orgID, Name: "Cat", IsActive: true,
		CreatedAt: now, UpdatedAt: now,
	}
	typ := &model.IncidentType{
		ID: typeID, OrganizationID: orgID, CategoryID: catID,
		Name: "Fall from bed", IsActive: true,
		CreatedAt: now, UpdatedAt: now,
	}

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.CategoryCreated(tx, cat); err != nil {
			return err
		}
		return projector.TypeCreated(tx, typ)
	}))

	reader := classifierread.NewReader(testDB, &logger)
	got, err := reader.GetType(ctx, typeID)
	require.NoError(t, err)
	require.Equal(t, "Fall from bed", got.Name)

	active, err := reader.ListActiveTypesByOrganization(ctx, orgID)
	require.NoError(t, err)
	require.Len(t, active, 1)
	require.Equal(t, typeID, active[0].ID)

	byCat, err := reader.ListTypesByCategory(ctx, catID)
	require.NoError(t, err)
	require.Len(t, byCat, 1)
}
