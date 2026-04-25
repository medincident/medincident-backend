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

	reader := classifierread.NewReader(testDB, authzSvc, &logger)
	got, err := reader.GetCategory(ctx, sysadminCaller, childID)
	require.NoError(t, err)
	require.Equal(t, "Falls", got.Name)
	require.NotNil(t, got.ParentCategoryID)
	require.Equal(t, rootID, *got.ParentCategoryID)

	roots, err := reader.ListActiveRootCategories(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, roots, 1)
	require.Equal(t, rootID, roots[0].ID)

	subtree, err := reader.ListCategorySubtree(ctx, sysadminCaller, rootID)
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

	reader := classifierread.NewReader(testDB, authzSvc, &logger)
	got, err := reader.GetType(ctx, sysadminCaller, typeID)
	require.NoError(t, err)
	require.Equal(t, "Fall from bed", got.Name)
	// Newly-created types start with patient submission disabled.
	require.False(t, got.IsAllowedForPatients)

	active, err := reader.ListActiveTypesByOrganization(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, active, 1)
	require.Equal(t, typeID, active[0].ID)

	byCat, err := reader.ListTypesByCategory(ctx, sysadminCaller, catID, classifierread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, byCat, 1)
}

// TestReader_PatientAllowed_Types_And_VisibleCategories seeds a fixture
// containing every interesting combination of (category active?, type
// active?, type allowed-for-patients?) and verifies that each patient-mode
// reader method honours its own contract:
//
//   - ListPatientAllowedTypesByOrganization filters on the type's own
//     is_active AND is_allowed_for_patients only; ancestor activity is the
//     visibility-tree's concern, not this list's.
//   - ListPatientVisibleCategoriesByOrganization returns categories that
//     are themselves active AND have at least one active+allowed type
//     somewhere in their subtree. It includes intermediate ancestors so
//     the patient sees the full path, and excludes categories whose every
//     descendant type is unavailable (or whose ancestor chain breaks at an
//     inactive category).
func TestReader_PatientAllowed_Types_And_VisibleCategories(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())

	// Fixture (same org):
	//
	//   surgical (active)              ← visible: has allowed type below
	//     wardFalls (active)           ← visible: direct allowed type
	//       allowedFall   (active, allowed)   ← appears in patient list
	//       internalOnly  (active, NOT allowed)
	//     archived (inactive)          ← invisible: ancestor inactive
	//       wouldAllow   (active, allowed)    ← excluded: parent inactive
	//   adminOnly (active)             ← invisible: no allowed type below
	//     adminOnlyType (active, NOT allowed)
	surgicalID := uuid.Must(uuid.NewV7())
	wardFallsID := uuid.Must(uuid.NewV7())
	archivedID := uuid.Must(uuid.NewV7())
	adminOnlyID := uuid.Must(uuid.NewV7())

	allowedFallID := uuid.Must(uuid.NewV7())
	internalOnlyID := uuid.Must(uuid.NewV7())
	wouldAllowID := uuid.Must(uuid.NewV7())
	adminOnlyTypeID := uuid.Must(uuid.NewV7())

	categories := []*model.IncidentCategory{
		{ID: surgicalID, OrganizationID: orgID, Name: "Surgical", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: wardFallsID, OrganizationID: orgID, ParentCategoryID: uuid.NullUUID{UUID: surgicalID, Valid: true}, Name: "Ward falls", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: archivedID, OrganizationID: orgID, ParentCategoryID: uuid.NullUUID{UUID: surgicalID, Valid: true}, Name: "Archived", IsActive: false, CreatedAt: now, UpdatedAt: now},
		{ID: adminOnlyID, OrganizationID: orgID, Name: "Admin-only", IsActive: true, CreatedAt: now, UpdatedAt: now},
	}
	types := []*model.IncidentType{
		{ID: allowedFallID, OrganizationID: orgID, CategoryID: wardFallsID, Name: "Patient fall", IsActive: true, IsAllowedForPatients: true, CreatedAt: now, UpdatedAt: now},
		{ID: internalOnlyID, OrganizationID: orgID, CategoryID: wardFallsID, Name: "Surgical retained item", IsActive: true, IsAllowedForPatients: false, CreatedAt: now, UpdatedAt: now},
		{ID: wouldAllowID, OrganizationID: orgID, CategoryID: archivedID, Name: "Allowed but orphaned", IsActive: true, IsAllowedForPatients: true, CreatedAt: now, UpdatedAt: now},
		{ID: adminOnlyTypeID, OrganizationID: orgID, CategoryID: adminOnlyID, Name: "Internal HR", IsActive: true, IsAllowedForPatients: false, CreatedAt: now, UpdatedAt: now},
	}

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, c := range categories {
			if err := projector.CategoryCreated(tx, c); err != nil {
				return err
			}
		}
		for _, t := range types {
			if err := projector.TypeCreated(tx, t); err != nil {
				return err
			}
		}
		return nil
	}))

	reader := classifierread.NewReader(testDB, authzSvc, &logger)

	// Patient-allowed types: only the active+allowed leaf under an active
	// category chain. wouldAllowID is allowed but its parent category is
	// inactive — included? It depends on intent. The current method filters
	// on the type's own is_active AND is_allowed_for_patients, NOT on
	// ancestor activity. The category-level filter (subtree visibility) is
	// the consumer's choice. We assert this contract explicitly.
	allowedTypes, err := reader.ListPatientAllowedTypesByOrganization(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
	require.NoError(t, err)
	allowedIDs := make(map[uuid.UUID]bool, len(allowedTypes))
	for _, tp := range allowedTypes {
		allowedIDs[tp.ID] = true
		require.True(t, tp.IsActive)
		require.True(t, tp.IsAllowedForPatients)
	}
	require.True(t, allowedIDs[allowedFallID], "patient list must include the active+allowed type")
	require.True(t, allowedIDs[wouldAllowID], "list filters by type-level flags only; ancestor activity is the visibility-tree's concern")
	require.False(t, allowedIDs[internalOnlyID], "non-allowed type must be excluded")
	require.False(t, allowedIDs[adminOnlyTypeID], "non-allowed type must be excluded")

	// Patient-visible categories: surgical (root, has allowed descendant) and
	// wardFalls (direct allowed type). Archived is inactive → excluded.
	// adminOnly has only non-allowed types → excluded. wouldAllow's parent
	// (archived) is inactive, so wouldAllow does not contribute visibility.
	visibleCats, err := reader.ListPatientVisibleCategoriesByOrganization(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
	require.NoError(t, err)
	visibleIDs := make(map[uuid.UUID]bool, len(visibleCats))
	for _, c := range visibleCats {
		visibleIDs[c.ID] = true
		require.True(t, c.IsActive, "patient-visible categories must themselves be active")
	}
	require.True(t, visibleIDs[surgicalID], "ancestor of an allowed type must be visible")
	require.True(t, visibleIDs[wardFallsID], "direct parent of an allowed type must be visible")
	require.False(t, visibleIDs[archivedID], "inactive ancestor must be excluded")
	require.False(t, visibleIDs[adminOnlyID], "category whose subtree has no allowed type must be excluded")
}
