//go:build integration

package projector_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
)

// TestCategory_LifecycleProjection runs a representative category flow
// through every projector touching incident_categories:
// Created → UpdateDetails → Move → Deactivate → Reactivate → Deleted.
func TestCategory_LifecycleProjection(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())
	catID := uuid.Must(uuid.NewV7())
	parentID := uuid.Must(uuid.NewV7())

	cat := &model.IncidentCategory{
		ID:             catID,
		OrganizationID: orgID,
		Name:           "Cat A",
		Description:    null.StringFrom("first description"),
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.CategoryCreated(tx, cat)
	}))

	cat.Name = "Cat B"
	cat.Description = null.StringFrom("updated description")
	cat.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.CategoryUpdateDetails(tx, cat)
	}))
	var name string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT name FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&name))
	require.Equal(t, "Cat B", name)

	// Move under a new parent.
	newParent := uuid.NullUUID{UUID: parentID, Valid: true}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.CategoryMove(tx, catID, newParent, now.Add(2*time.Hour))
	}))
	var gotParent *uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT parent_category_id FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&gotParent))
	require.NotNil(t, gotParent)
	require.Equal(t, parentID, *gotParent)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.CategoryDeactivate(tx, catID, now.Add(3*time.Hour))
	}))
	var active bool
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT is_active FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&active))
	require.False(t, active)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.CategoryReactivate(tx, catID, now.Add(4*time.Hour))
	}))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT is_active FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&active))
	require.True(t, active)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.CategoryDeleted(tx, catID)
	}))
	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}

// TestType_LifecycleProjection runs an analogous flow for incident
// types: Created → UpdateDetails → Move → Deactivate → Reactivate → Deleted.
func TestType_LifecycleProjection(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())
	catID := uuid.Must(uuid.NewV7())
	newCatID := uuid.Must(uuid.NewV7())
	typeID := uuid.Must(uuid.NewV7())

	typ := &model.IncidentType{
		ID:             typeID,
		OrganizationID: orgID,
		CategoryID:     catID,
		Name:           "Type A",
		Description:    null.StringFrom("first description"),
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.TypeCreated(tx, typ)
	}))

	typ.Name = "Type B"
	typ.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.TypeUpdateDetails(tx, typ)
	}))

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.TypeMove(tx, typeID, newCatID, now.Add(2*time.Hour))
	}))
	var gotCatID uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT category_id FROM projections.incident_types WHERE id = ?`, typeID,
	).Row().Scan(&gotCatID))
	require.Equal(t, newCatID, gotCatID)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.TypeDeactivate(tx, typeID, now.Add(3*time.Hour))
	}))
	var active bool
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT is_active FROM projections.incident_types WHERE id = ?`, typeID,
	).Row().Scan(&active))
	require.False(t, active)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.TypeReactivate(tx, typeID, now.Add(4*time.Hour))
	}))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT is_active FROM projections.incident_types WHERE id = ?`, typeID,
	).Row().Scan(&active))
	require.True(t, active)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.TypeDeleted(tx, typeID)
	}))
	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.incident_types WHERE id = ?`, typeID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}
