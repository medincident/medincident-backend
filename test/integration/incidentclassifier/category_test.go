//go:build integration

package incidentclassifier_integration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
)

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

// uuidStrPtr returns a *string whose value is id.String(). Used when
// constructing payloads whose optional UUID fields are typed as
// *string.
func uuidStrPtr(id uuid.UUID) *string {
	s := id.String()
	return &s
}

func TestCategory_CreateRoot(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Test Org 1")

	desc := "Root category for surgical incidents"
	result, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Surgical",
			Description:    &desc,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, result.ID)

	row := loadCategory(t, result.ID)
	assert.Equal(t, "Surgical", row.Name)
	assert.False(t, row.ParentCategoryID.Valid)
	assert.True(t, row.IsActive)
	assert.True(t, row.Description.Valid)

	// Sync projector must have written the matching projection row
	// inside the same transaction as the domain insert.
	var projName string
	var projActive bool
	require.NoError(t, testDB.Raw(
		`SELECT name, is_active FROM projections.incident_categories WHERE id = ?`,
		result.ID,
	).Row().Scan(&projName, &projActive))
	assert.Equal(t, "Surgical", projName)
	assert.True(t, projActive)
}

func TestCategory_CreateChild(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	root, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Clinical",
		},
	})
	require.NoError(t, err)

	child, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(root.ID),
			Name:             "Medication error",
		},
	})
	require.NoError(t, err)

	row := loadCategory(t, child.ID)
	assert.True(t, row.ParentCategoryID.Valid)
	assert.Equal(t, root.ID, row.ParentCategoryID.UUID)
}

func TestCategory_CreateMaxDepthExceeded(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	var parent *string
	for i := 0; i < 5; i++ {
		res, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
			Caller: sysadminCaller,
			Payload: classifiersvc.CreateIncidentCategoryPayload{
				OrganizationID:   orgID.String(),
				ParentCategoryID: parent,
				Name:             names(i),
			},
		})
		require.NoError(t, err)
		parent = uuidStrPtr(res.ID)
	}

	_, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: parent,
			Name:             "level6",
		},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryMaxDepthExceeded, codeOf(t, err))
}

func TestCategory_CreateGlobalNameConflict(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	root, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Duplicate",
		},
	})
	require.NoError(t, err)

	_, err = categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(root.ID),
			Name:             "Duplicate",
		},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryNameConflict, codeOf(t, err))
}

func TestCategory_CreateCrossOrgParent(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgA := insertOrganization(t, "Org A")
	orgB := insertOrganization(t, "Org B")

	parentInA, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgA.String(),
			Name:           "ParentInA",
		},
	})
	require.NoError(t, err)

	_, err = categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgB.String(),
			ParentCategoryID: uuidStrPtr(parentInA.ID),
			Name:             "Child",
		},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryParentOrganizationMismatch, codeOf(t, err))
}

func TestCategory_UpdateDetails(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	res, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Old name",
		},
	})
	require.NoError(t, err)

	desc := "New description text"
	_, err = categorySvc.UpdateDetails(ctx, classifiersvc.UpdateIncidentCategoryDetailsCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.UpdateIncidentCategoryDetailsPayload{
			CategoryID:  res.ID.String(),
			Name:        "New name",
			Description: &desc,
		},
	})
	require.NoError(t, err)

	row := loadCategory(t, res.ID)
	assert.Equal(t, "New name", row.Name)
	assert.True(t, row.Description.Valid)
}

func TestCategory_MoveToRootAndUnderNewParent(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	rootA, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "RootA",
		},
	})
	require.NoError(t, err)
	rootB, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "RootB",
		},
	})
	require.NoError(t, err)
	child, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(rootA.ID),
			Name:             "Child",
		},
	})
	require.NoError(t, err)

	// Move child under rootB
	_, err = categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID:          child.ID.String(),
			NewParentCategoryID: uuidStrPtr(rootB.ID),
		},
	})
	require.NoError(t, err)
	row := loadCategory(t, child.ID)
	assert.Equal(t, rootB.ID, row.ParentCategoryID.UUID)

	// Move child to root level
	_, err = categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID: child.ID.String(),
		},
	})
	require.NoError(t, err)
	row = loadCategory(t, child.ID)
	assert.False(t, row.ParentCategoryID.Valid)
}

func TestCategory_MoveWouldCreateCycle(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	a, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "AA",
		},
	})
	require.NoError(t, err)
	b, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(a.ID),
			Name:             "BB",
		},
	})
	require.NoError(t, err)
	c, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(b.ID),
			Name:             "CC",
		},
	})
	require.NoError(t, err)

	// Move A under C (C is descendant of A) → cycle
	_, err = categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID:          a.ID.String(),
			NewParentCategoryID: uuidStrPtr(c.ID),
		},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryMoveWouldCreateCycle, codeOf(t, err))
}

func TestCategory_MoveWouldExceedDepth(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	// Build chain A(1) -> B(2) -> C(3)  (subtree depth from A = 3)
	a, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "AA",
		},
	})
	require.NoError(t, err)
	b, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(a.ID),
			Name:             "BB",
		},
	})
	require.NoError(t, err)
	_, err = categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(b.ID),
			Name:             "CC",
		},
	})
	require.NoError(t, err)

	// Build chain X(1) -> Y(2) -> Z(3) (target parent depth = 3)
	x, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "XX",
		},
	})
	require.NoError(t, err)
	y, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(x.ID),
			Name:             "YY",
		},
	})
	require.NoError(t, err)
	z, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(y.ID),
			Name:             "ZZ",
		},
	})
	require.NoError(t, err)

	// Move A under Z: 3 + 3 = 6 > 5 → error
	_, err = categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID:          a.ID.String(),
			NewParentCategoryID: uuidStrPtr(z.ID),
		},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryMoveWouldExceedDepth, codeOf(t, err))
}

func TestCategory_MoveCrossOrg(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgA := insertOrganization(t, "Org A")
	orgB := insertOrganization(t, "Org B")

	a, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgA.String(),
			Name:           "InA",
		},
	})
	require.NoError(t, err)
	b, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgB.String(),
			Name:           "InB",
		},
	})
	require.NoError(t, err)

	_, err = categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID:          a.ID.String(),
			NewParentCategoryID: uuidStrPtr(b.ID),
		},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryMoveOrganizationMismatch, codeOf(t, err))
}

func TestCategory_DeactivateCascades(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	// root → child → grand, types on root and grand
	root, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Root",
		},
	})
	require.NoError(t, err)
	child, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(root.ID),
			Name:             "Child",
		},
	})
	require.NoError(t, err)
	grand, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(child.ID),
			Name:             "Grand",
		},
	})
	require.NoError(t, err)
	typeOnRoot, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{
			CategoryID: root.ID.String(),
			Name:       "TypeOnRoot",
		},
	})
	require.NoError(t, err)
	typeOnGrand, err := typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{
			CategoryID: grand.ID.String(),
			Name:       "TypeOnGrand",
		},
	})
	require.NoError(t, err)

	// Clear creation events so we can count only deactivate events.

	_, err = categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateIncidentCategoryPayload{CategoryID: root.ID.String()},
	})
	require.NoError(t, err)

	assert.False(t, loadCategory(t, root.ID).IsActive)
	assert.False(t, loadCategory(t, child.ID).IsActive)
	assert.False(t, loadCategory(t, grand.ID).IsActive)
	assert.False(t, loadType(t, typeOnRoot.ID).IsActive)
	assert.False(t, loadType(t, typeOnGrand.ID).IsActive)
}

func TestCategory_ReactivateBlockedByInactiveAncestor(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	root, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Root",
		},
	})
	require.NoError(t, err)
	child, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(root.ID),
			Name:             "Child",
		},
	})
	require.NoError(t, err)

	_, err = categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateIncidentCategoryPayload{CategoryID: root.ID.String()},
	})
	require.NoError(t, err)

	_, err = categorySvc.Reactivate(ctx, classifiersvc.ReactivateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.ReactivateIncidentCategoryPayload{CategoryID: child.ID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryReactivateInactiveAncestor, codeOf(t, err))
}

func TestCategory_ReactivateRootOnly(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	root, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Root",
		},
	})
	require.NoError(t, err)
	child, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(root.ID),
			Name:             "Child",
		},
	})
	require.NoError(t, err)

	_, err = categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateIncidentCategoryPayload{CategoryID: root.ID.String()},
	})
	require.NoError(t, err)

	_, err = categorySvc.Reactivate(ctx, classifiersvc.ReactivateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.ReactivateIncidentCategoryPayload{CategoryID: root.ID.String()},
	})
	require.NoError(t, err)

	assert.True(t, loadCategory(t, root.ID).IsActive)
	assert.False(t, loadCategory(t, child.ID).IsActive, "child must stay inactive — reactivation is single-node")
}

func TestCategory_ReactivateNameConflict(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	first, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Duplicate",
		},
	})
	require.NoError(t, err)
	_, err = categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateIncidentCategoryPayload{CategoryID: first.ID.String()},
	})
	require.NoError(t, err)

	_, err = categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Duplicate",
		},
	})
	require.NoError(t, err, "creating a new active Duplicate should succeed because the old one is inactive")

	_, err = categorySvc.Reactivate(ctx, classifiersvc.ReactivateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.ReactivateIncidentCategoryPayload{CategoryID: first.ID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryReactivateNameConflict, codeOf(t, err))
}

func TestCategory_DeleteCascades(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrganization(t, "Org")

	root, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Root",
		},
	})
	require.NoError(t, err)
	child, err := categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: uuidStrPtr(root.ID),
			Name:             "Child",
		},
	})
	require.NoError(t, err)
	_, err = typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{
			CategoryID: child.ID.String(),
			Name:       "TypeInChild",
		},
	})
	require.NoError(t, err)
	_, err = typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{
			CategoryID: root.ID.String(),
			Name:       "TypeInRoot",
		},
	})
	require.NoError(t, err)

	_, err = categorySvc.Delete(ctx, classifiersvc.DeleteIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeleteIncidentCategoryPayload{CategoryID: root.ID.String()},
	})
	require.NoError(t, err)

	assert.Equal(t, 0, countIncidentCategories(t))
	assert.Equal(t, 0, countIncidentTypes(t))
}

// names generates short deterministic names for depth tests.
func names(i int) string {
	return "level" + string(rune('a'+i))
}
