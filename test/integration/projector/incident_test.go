//go:build integration

package projector_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	classifierv1 "github.com/medincident/medincident-backend/pkg/event/incident/classifier/v1"
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
		return qprojector.CategoryCreated(tx, cat.ID.String(), cat.CreatedAt, &classifierv1.IncidentCategoryCreated{
			OrganizationId: cat.OrganizationID.String(),
			Name:           cat.Name,
			Description:    wrapperspb.String(cat.Description.String),
			IsActive:       cat.IsActive,
			CreatedAt:      timestamppb.New(cat.CreatedAt),
		})
	}))

	cat.Name = "Cat B"
	cat.Description = null.StringFrom("updated description")
	cat.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.CategoryDetailsUpdated(tx, cat.ID.String(), cat.UpdatedAt, &classifierv1.IncidentCategoryDetailsUpdated{
			Name:        cat.Name,
			Description: wrapperspb.String(cat.Description.String),
			UpdatedAt:   timestamppb.New(cat.UpdatedAt),
		})
	}))
	var name string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT name FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&name))
	require.Equal(t, "Cat B", name)

	// Move under a new parent.
	newParent := uuid.NullUUID{UUID: parentID, Valid: true}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.CategoryMoved(tx, catID.String(), now.Add(2*time.Hour), &classifierv1.IncidentCategoryMoved{
			NewParentCategoryId: wrapperspb.String(newParent.UUID.String()),
			UpdatedAt:           timestamppb.New(now.Add(2 * time.Hour)),
		})
	}))
	var gotParent *uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT parent_category_id FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&gotParent))
	require.NotNil(t, gotParent)
	require.Equal(t, parentID, *gotParent)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.CategoryDeactivated(tx, catID.String(), now.Add(3*time.Hour), &classifierv1.IncidentCategoryDeactivated{
			UpdatedAt: timestamppb.New(now.Add(3 * time.Hour)),
		})
	}))
	var active bool
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT is_active FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&active))
	require.False(t, active)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.CategoryReactivated(tx, catID.String(), now.Add(4*time.Hour), &classifierv1.IncidentCategoryReactivated{
			UpdatedAt: timestamppb.New(now.Add(4 * time.Hour)),
		})
	}))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT is_active FROM projections.incident_categories WHERE id = ?`, catID,
	).Row().Scan(&active))
	require.True(t, active)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.CategoryDeleted(tx, catID.String(), time.Now(), &classifierv1.IncidentCategoryDeleted{})
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
		return qprojector.TypeCreated(tx, typ.ID.String(), typ.CreatedAt, &classifierv1.IncidentTypeCreated{
			OrganizationId: typ.OrganizationID.String(),
			CategoryId:     typ.CategoryID.String(),
			Name:           typ.Name,
			Description:    wrapperspb.String(typ.Description.String),
			IsActive:       typ.IsActive,
			CreatedAt:      timestamppb.New(typ.CreatedAt),
		})
	}))

	typ.Name = "Type B"
	typ.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.TypeDetailsUpdated(tx, typ.ID.String(), typ.UpdatedAt, &classifierv1.IncidentTypeDetailsUpdated{
			Name:      typ.Name,
			UpdatedAt: timestamppb.New(typ.UpdatedAt),
		})
	}))

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.TypeMoved(tx, typeID.String(), now.Add(2*time.Hour), &classifierv1.IncidentTypeMoved{
			NewCategoryId: newCatID.String(),
			UpdatedAt:     timestamppb.New(now.Add(2 * time.Hour)),
		})
	}))
	var gotCatID uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT category_id FROM projections.incident_types WHERE id = ?`, typeID,
	).Row().Scan(&gotCatID))
	require.Equal(t, newCatID, gotCatID)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.TypeDeactivated(tx, typeID.String(), now.Add(3*time.Hour), &classifierv1.IncidentTypeDeactivated{
			UpdatedAt: timestamppb.New(now.Add(3 * time.Hour)),
		})
	}))
	var active bool
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT is_active FROM projections.incident_types WHERE id = ?`, typeID,
	).Row().Scan(&active))
	require.False(t, active)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.TypeReactivated(tx, typeID.String(), now.Add(4*time.Hour), &classifierv1.IncidentTypeReactivated{
			UpdatedAt: timestamppb.New(now.Add(4 * time.Hour)),
		})
	}))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT is_active FROM projections.incident_types WHERE id = ?`, typeID,
	).Row().Scan(&active))
	require.True(t, active)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.TypeDeleted(tx, typeID.String(), time.Now(), &classifierv1.IncidentTypeDeleted{})
	}))
	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.incident_types WHERE id = ?`, typeID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}
