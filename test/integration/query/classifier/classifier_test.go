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
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	classifierread "github.com/medincident/medincident-backend/internal/service/query/incident/classifier"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	classifierv1 "github.com/medincident/medincident-backend/pkg/event/incident/classifier/v1"
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
		if err := qprojector.CategoryCreated(tx, root.ID.String(), root.CreatedAt, &classifierv1.IncidentCategoryCreated{
			OrganizationId: root.OrganizationID.String(),
			Name:           root.Name,
			Description:    wrapperspb.String(root.Description.String),
			IsActive:       root.IsActive,
			CreatedAt:      timestamppb.New(root.CreatedAt),
		}); err != nil {
			return err
		}
		if err := qprojector.CategoryCreated(tx, child.ID.String(), child.CreatedAt, &classifierv1.IncidentCategoryCreated{
			OrganizationId:   child.OrganizationID.String(),
			Name:             child.Name,
			ParentCategoryId: wrapperspb.String(child.ParentCategoryID.UUID.String()),
			IsActive:         child.IsActive,
			CreatedAt:        timestamppb.New(child.CreatedAt),
		}); err != nil {
			return err
		}
		return qprojector.CategoryCreated(tx, grand.ID.String(), grand.CreatedAt, &classifierv1.IncidentCategoryCreated{
			OrganizationId:   grand.OrganizationID.String(),
			Name:             grand.Name,
			ParentCategoryId: wrapperspb.String(grand.ParentCategoryID.UUID.String()),
			IsActive:         grand.IsActive,
			CreatedAt:        timestamppb.New(grand.CreatedAt),
		})
	}))

	reader := classifierread.NewReader(testDB, authzSvc, &logger)
	got, err := reader.GetCategory(ctx, sysadminCaller, childID)
	require.NoError(t, err)
	require.Equal(t, "Falls", got.Name)
	require.NotNil(t, got.ParentCategoryID)
	require.Equal(t, rootID, *got.ParentCategoryID)

	roots, err := reader.ListActiveRootCategories(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, roots.Items, 1)
	require.Equal(t, rootID, roots.Items[0].ID)

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
		if err := qprojector.CategoryCreated(tx, cat.ID.String(), cat.CreatedAt, &classifierv1.IncidentCategoryCreated{
			OrganizationId: cat.OrganizationID.String(),
			Name:           cat.Name,
			IsActive:       cat.IsActive,
			CreatedAt:      timestamppb.New(cat.CreatedAt),
		}); err != nil {
			return err
		}
		return qprojector.TypeCreated(tx, typ.ID.String(), typ.CreatedAt, &classifierv1.IncidentTypeCreated{
			OrganizationId: typ.OrganizationID.String(),
			CategoryId:     typ.CategoryID.String(),
			Name:           typ.Name,
			IsActive:       typ.IsActive,
			CreatedAt:      timestamppb.New(typ.CreatedAt),
		})
	}))

	reader := classifierread.NewReader(testDB, authzSvc, &logger)
	got, err := reader.GetType(ctx, sysadminCaller, typeID)
	require.NoError(t, err)
	require.Equal(t, "Fall from bed", got.Name)
	// Newly-created types start with patient submission disabled.
	require.False(t, got.IsAllowedForPatients)

	active, err := reader.ListActiveTypesByOrganization(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, active.Items, 1)
	require.Equal(t, typeID, active.Items[0].ID)

	byCat, err := reader.ListTypesByCategory(ctx, sysadminCaller, catID, classifierread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, byCat.Items, 1)
}

// TestReader_UnifiedPatientFiltering seeds a fixture containing every
// interesting combination of (category active?, type active?, type
// allowed-for-patients?) and verifies that all 5 list methods honour their
// contracts for both the employee caller (sysadminCaller) and the patient
// caller (patientCaller).
//
// Fixture:
//
//	Categories (org-level):
//	  surgicalID   — active, root
//	  wardFallsID  — active, parent=surgicalID
//	  archivedID   — inactive, parent=surgicalID
//	  adminOnlyID  — active, root
//
//	Types:
//	  allowedFallID   in wardFallsID  — active, is_allowed_for_patients=true
//	  internalOnlyID  in wardFallsID  — active, is_allowed_for_patients=false
//	  wouldAllowID    in archivedID   — active, is_allowed_for_patients=true
//	  adminOnlyTypeID in adminOnlyID  — active, is_allowed_for_patients=false
func TestReader_UnifiedPatientFiltering(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())

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
			ev := &classifierv1.IncidentCategoryCreated{
				OrganizationId: c.OrganizationID.String(),
				Name:           c.Name,
				IsActive:       c.IsActive,
				CreatedAt:      timestamppb.New(c.CreatedAt),
			}
			if c.ParentCategoryID.Valid {
				ev.ParentCategoryId = wrapperspb.String(c.ParentCategoryID.UUID.String())
			}
			if err := qprojector.CategoryCreated(tx, c.ID.String(), c.CreatedAt, ev); err != nil {
				return err
			}
		}
		for _, tp := range types {
			if err := qprojector.TypeCreated(tx, tp.ID.String(), tp.CreatedAt, &classifierv1.IncidentTypeCreated{
				OrganizationId:       tp.OrganizationID.String(),
				CategoryId:           tp.CategoryID.String(),
				Name:                 tp.Name,
				IsActive:             tp.IsActive,
				IsAllowedForPatients: tp.IsAllowedForPatients,
				CreatedAt:            timestamppb.New(tp.CreatedAt),
			}); err != nil {
				return err
			}
		}
		return nil
	}))

	reader := classifierread.NewReader(testDB, authzSvc, &logger)

	// -------------------------------------------------------------------------
	// Patient caller: filtered views
	// -------------------------------------------------------------------------

	t.Run("patient/ListActiveTypesByOrganization", func(t *testing.T) {
		result, err := reader.ListActiveTypesByOrganization(ctx, patientCaller, orgID, classifierread.ListQuery{})
		require.NoError(t, err)
		ids := typeIDs(result.Items)
		// Active + patient-allowed types regardless of category activity.
		require.True(t, ids[allowedFallID], "active+allowed type in active category must be included")
		require.True(t, ids[wouldAllowID], "active+allowed type in inactive category must be included (type-level filter only)")
		require.False(t, ids[internalOnlyID], "non-allowed type must be excluded")
		require.False(t, ids[adminOnlyTypeID], "non-allowed type must be excluded")
		for _, tp := range result.Items {
			require.True(t, tp.IsActive)
			require.True(t, tp.IsAllowedForPatients)
		}
	})

	t.Run("patient/ListCategoriesByOrganization", func(t *testing.T) {
		result, err := reader.ListCategoriesByOrganization(ctx, patientCaller, orgID, classifierread.ListQuery{})
		require.NoError(t, err)
		ids := categoryIDs(result.Items)
		// Active categories whose subtree contains at least one active+patient-allowed type.
		require.True(t, ids[surgicalID], "ancestor of an allowed type must be visible")
		require.True(t, ids[wardFallsID], "direct parent of an allowed type must be visible")
		require.False(t, ids[archivedID], "inactive category must be excluded")
		require.False(t, ids[adminOnlyID], "category whose subtree has no allowed type must be excluded")
		for _, c := range result.Items {
			require.True(t, c.IsActive, "all returned categories must be active")
		}
	})

	t.Run("patient/ListActiveRootCategories", func(t *testing.T) {
		result, err := reader.ListActiveRootCategories(ctx, patientCaller, orgID, classifierread.ListQuery{})
		require.NoError(t, err)
		ids := categoryIDs(result.Items)
		// Active roots whose subtree contains at least one active+patient-allowed type.
		require.True(t, ids[surgicalID], "root with allowed types in subtree must be included")
		require.False(t, ids[adminOnlyID], "root with no allowed types in subtree must be excluded")
		require.False(t, ids[wardFallsID], "non-root must not appear")
		require.False(t, ids[archivedID], "inactive category must be excluded")
	})

	t.Run("patient/ListTypesByCategory", func(t *testing.T) {
		result, err := reader.ListTypesByCategory(ctx, patientCaller, wardFallsID, classifierread.ListQuery{})
		require.NoError(t, err)
		ids := typeIDs(result.Items)
		require.True(t, ids[allowedFallID], "active+allowed type must be included")
		require.False(t, ids[internalOnlyID], "non-allowed type must be excluded")
	})

	t.Run("patient/ListCategorySubtree", func(t *testing.T) {
		result, err := reader.ListCategorySubtree(ctx, patientCaller, surgicalID)
		require.NoError(t, err)
		ids := categoryIDsSlice(result)
		// 3-CTE walk: down from surgicalID (active only in recursive part),
		// direct = wardFallsID (has allowedFallID), visible = wardFallsID + surgicalID (ancestor).
		// archivedID is inactive so the recursive descent skips it.
		require.True(t, ids[surgicalID], "root must be included as an ancestor of wardFallsID")
		require.True(t, ids[wardFallsID], "direct parent of an allowed type must be included")
		require.False(t, ids[archivedID], "inactive category must be excluded from patient subtree walk")
		require.False(t, ids[adminOnlyID], "category outside subtree must not appear")
	})

	// -------------------------------------------------------------------------
	// Employee caller (sysadmin): unfiltered views
	// -------------------------------------------------------------------------

	t.Run("employee/ListActiveTypesByOrganization", func(t *testing.T) {
		result, err := reader.ListActiveTypesByOrganization(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
		require.NoError(t, err)
		ids := typeIDs(result.Items)
		// All 4 types are active; employee sees all without patient filter.
		require.True(t, ids[allowedFallID])
		require.True(t, ids[internalOnlyID])
		require.True(t, ids[wouldAllowID])
		require.True(t, ids[adminOnlyTypeID])
	})

	t.Run("employee/ListCategoriesByOrganization", func(t *testing.T) {
		result, err := reader.ListCategoriesByOrganization(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
		require.NoError(t, err)
		ids := categoryIDs(result.Items)
		// Employee sees all 4 categories regardless of active status.
		require.True(t, ids[surgicalID])
		require.True(t, ids[wardFallsID])
		require.True(t, ids[archivedID])
		require.True(t, ids[adminOnlyID])
	})

	t.Run("employee/ListActiveRootCategories", func(t *testing.T) {
		result, err := reader.ListActiveRootCategories(ctx, sysadminCaller, orgID, classifierread.ListQuery{})
		require.NoError(t, err)
		ids := categoryIDs(result.Items)
		// Both active roots visible to employee.
		require.True(t, ids[surgicalID])
		require.True(t, ids[adminOnlyID])
		require.False(t, ids[wardFallsID], "non-root must not appear")
		require.False(t, ids[archivedID], "inactive must not appear in active-root list")
	})

	t.Run("employee/ListTypesByCategory", func(t *testing.T) {
		result, err := reader.ListTypesByCategory(ctx, sysadminCaller, wardFallsID, classifierread.ListQuery{})
		require.NoError(t, err)
		ids := typeIDs(result.Items)
		require.True(t, ids[allowedFallID])
		require.True(t, ids[internalOnlyID])
	})

	t.Run("employee/ListCategorySubtree", func(t *testing.T) {
		result, err := reader.ListCategorySubtree(ctx, sysadminCaller, surgicalID)
		require.NoError(t, err)
		ids := categoryIDsSlice(result)
		// Employee sees full subtree including inactive archivedID.
		require.True(t, ids[surgicalID])
		require.True(t, ids[wardFallsID])
		require.True(t, ids[archivedID])
		require.False(t, ids[adminOnlyID], "category outside subtree must not appear")
	})
}

// categoryIDs returns a set of IDs from a CategoryListResult items slice.
func categoryIDs(items []classifierread.CategoryView) map[uuid.UUID]bool {
	m := make(map[uuid.UUID]bool, len(items))
	for _, c := range items {
		m[c.ID] = true
	}
	return m
}

// categoryIDsSlice returns a set of IDs from a []CategoryView slice.
func categoryIDsSlice(items []classifierread.CategoryView) map[uuid.UUID]bool {
	return categoryIDs(items)
}

// typeIDs returns a set of IDs from a TypeListResult items slice.
func typeIDs(items []classifierread.TypeView) map[uuid.UUID]bool {
	m := make(map[uuid.UUID]bool, len(items))
	for _, tp := range items {
		m[tp.ID] = true
	}
	return m
}
