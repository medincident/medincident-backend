//go:build integration

package incidentclassifier_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
)

func TestType_CreateUnderRoot(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	cat, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "Surgical",
	})
	require.NoError(t, err)

	desc := "Fall of patient on ward"
	res, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID:  cat.ID,
		Name:        "Patient fall",
		Description: &desc,
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	row := loadType(t, res.ID)
	assert.Equal(t, "Patient fall", row.Name)
	assert.Equal(t, cat.ID, row.CategoryID)
	assert.Equal(t, orgID, row.OrganizationID)
	assert.True(t, row.IsActive)

	// Sync projector must have written the matching projection row.
	var projName string
	var projActive bool
	require.NoError(t, testDB.Raw(
		`SELECT name, is_active FROM projections.incident_types WHERE id = ?`,
		res.ID,
	).Row().Scan(&projName, &projActive))
	assert.Equal(t, "Patient fall", projName)
	assert.True(t, projActive)
}

func TestType_CreateUnderNestedCategory(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	root, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "Root",
	})
	require.NoError(t, err)
	nested, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID:   orgID,
		ParentCategoryID: &root.ID,
		Name:             "Nested",
	})
	require.NoError(t, err)

	_, err = typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: nested.ID,
		Name:       "Leaf type",
	})
	require.NoError(t, err)
}

func TestType_CreateCategoryNotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	missing := uuid.New()
	_, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: missing,
		Name:       "Ghost",
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentTypeCategoryNotFound, codeOf(t, err))
}

func TestType_CreateGlobalNameConflict(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	c1, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "C1",
	})
	require.NoError(t, err)
	c2, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "C2",
	})
	require.NoError(t, err)

	_, err = typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: c1.ID,
		Name:       "SameName",
	})
	require.NoError(t, err)
	_, err = typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: c2.ID,
		Name:       "SameName",
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentTypeNameConflict, codeOf(t, err))
}

func TestType_UpdateDetails(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")
	cat, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "Cat",
	})
	require.NoError(t, err)
	res, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: cat.ID,
		Name:       "Old",
	})
	require.NoError(t, err)

	desc := "A better description"
	_, err = typeSvc.UpdateDetails(ctx, classifiersvc.UpdateIncidentTypeDetailsCommand{
		TypeID:      res.ID,
		Name:        "New name",
		Description: &desc,
	})
	require.NoError(t, err)

	row := loadType(t, res.ID)
	assert.Equal(t, "New name", row.Name)
	assert.True(t, row.Description.Valid)
}

func TestType_MoveSameOrg(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	c1, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "C1",
	})
	require.NoError(t, err)
	c2, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "C2",
	})
	require.NoError(t, err)
	res, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: c1.ID,
		Name:       "Movable",
	})
	require.NoError(t, err)

	_, err = typeSvc.Move(ctx, classifiersvc.MoveIncidentTypeCommand{
		TypeID:        res.ID,
		NewCategoryID: c2.ID,
	})
	require.NoError(t, err)
	row := loadType(t, res.ID)
	assert.Equal(t, c2.ID, row.CategoryID)
}

func TestType_MoveCrossOrg(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgA := insertOrganization(t, "Org A")
	orgB := insertOrganization(t, "Org B")

	catA, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgA,
		Name:           "Cat A",
	})
	require.NoError(t, err)
	catB, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgB,
		Name:           "Cat B",
	})
	require.NoError(t, err)
	res, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: catA.ID,
		Name:       "Type",
	})
	require.NoError(t, err)

	_, err = typeSvc.Move(ctx, classifiersvc.MoveIncidentTypeCommand{
		TypeID:        res.ID,
		NewCategoryID: catB.ID,
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentTypeMoveOrganizationMismatch, codeOf(t, err))
}

func TestType_Deactivate(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")
	cat, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "Cat",
	})
	require.NoError(t, err)
	res, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: cat.ID,
		Name:       "Type1",
	})
	require.NoError(t, err)

	_, err = typeSvc.Deactivate(ctx, classifiersvc.DeactivateIncidentTypeCommand{TypeID: res.ID})
	require.NoError(t, err)
	assert.False(t, loadType(t, res.ID).IsActive)
}

func TestType_ReactivateBlockedByInactiveCategory(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	cat, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "Cat",
	})
	require.NoError(t, err)
	res, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: cat.ID,
		Name:       "Type1",
	})
	require.NoError(t, err)

	// Deactivate category (cascades to the type too).
	_, err = categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{CategoryID: cat.ID})
	require.NoError(t, err)

	_, err = typeSvc.Reactivate(ctx, classifiersvc.ReactivateIncidentTypeCommand{TypeID: res.ID})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentTypeReactivateInactiveAncestor, codeOf(t, err))
}

func TestType_Delete(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")
	cat, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: orgID,
		Name:           "Cat",
	})
	require.NoError(t, err)
	res, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID: cat.ID,
		Name:       "Type1",
	})
	require.NoError(t, err)

	_, err = typeSvc.Delete(ctx, classifiersvc.DeleteIncidentTypeCommand{TypeID: res.ID})
	require.NoError(t, err)
	assert.Equal(t, 0, countIncidentTypes(t))
}
