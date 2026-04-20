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

// TestOrganizationCreated_WritesRowAndCounters confirms the projector
// writes both the projections.organizations row AND the matching
// projections.organization_counters row inside the same transaction.
func TestOrganizationCreated_WritesRowAndCounters(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	id := uuid.Must(uuid.NewV7())
	now := time.Now().UTC().Truncate(time.Second)
	org := &model.Organization{
		ID:           id,
		Name:         "Acme Clinics",
		Description:  null.StringFrom("a longer description above min"),
		LegalAddress: model.Address{Text: "1 Real Street"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrganizationCreated(tx, org)
	}))

	var name string
	var desc *string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT name, description FROM projections.organizations WHERE id = ?`, id).
		Row().Scan(&name, &desc))
	require.Equal(t, "Acme Clinics", name)
	require.NotNil(t, desc)
	require.Equal(t, "a longer description above min", *desc)

	var employees, clinics, departments int64
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT employees_total, clinics_total, departments_total
		FROM projections.organization_counters WHERE organization_id = ?`, id).
		Row().Scan(&employees, &clinics, &departments))
	require.Equal(t, int64(0), employees)
	require.Equal(t, int64(0), clinics)
	require.Equal(t, int64(0), departments)
}

// TestOrganizationDetailsChanged_UpdatesRowAndEmployeeCards confirms
// the projector updates both projections.organizations and the
// organization_name mirror on projections.employee_cards.
func TestOrganizationDetailsChanged_UpdatesRowAndEmployeeCards(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	orgID := uuid.Must(uuid.NewV7())
	now := time.Now().UTC().Truncate(time.Second)
	org := &model.Organization{
		ID:           orgID,
		Name:         "Original Name",
		Description:  null.StringFrom("original description above min"),
		LegalAddress: model.Address{Text: "1 Real Street"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrganizationCreated(tx, org)
	}))

	// Seed an employee_card row that should be re-mirrored on rename.
	require.NoError(t, testDB.WithContext(ctx).Exec(`
		INSERT INTO projections.employee_cards
		    (employee_id, zitadel_user_id, organization_id, organization_name,
		     department_id, updated_at)
		VALUES (?, 'zit-1', ?, 'Original Name', ?, ?)`,
		uuid.Must(uuid.NewV7()), orgID, uuid.Must(uuid.NewV7()), now,
	).Error)

	// Mutate and project the change.
	org.Name = "New Name"
	org.Description = null.StringFrom("new description above min")
	org.UpdatedAt = now.Add(time.Hour)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrganizationDetailsChanged(tx, org)
	}))

	var name string
	var desc *string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT name, description FROM projections.organizations WHERE id = ?`, orgID).
		Row().Scan(&name, &desc))
	require.Equal(t, "New Name", name)
	require.Equal(t, "new description above min", *desc)

	var cardOrgName string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT organization_name FROM projections.employee_cards WHERE organization_id = ?`, orgID).
		Row().Scan(&cardOrgName))
	require.Equal(t, "New Name", cardOrgName)
}

// TestOrganizationLegalAddressChanged_UpdatesRow confirms the projector
// updates legal_address_text/longitude/latitude on projections.organizations.
func TestOrganizationLegalAddressChanged_UpdatesRow(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	orgID := uuid.Must(uuid.NewV7())
	now := time.Now().UTC().Truncate(time.Second)
	org := &model.Organization{
		ID:           orgID,
		Name:         "Acme",
		LegalAddress: model.Address{Text: "1 Real Street"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrganizationCreated(tx, org)
	}))

	org.LegalAddress = model.Address{
		Text:  "2 Different Avenue",
		Point: null.ValueFrom(model.Point{Longitude: 12.5, Latitude: -7.25}),
	}
	org.UpdatedAt = now.Add(time.Hour)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrganizationLegalAddressChanged(tx, org)
	}))

	var addrText string
	var lon, lat *float64
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT legal_address_text, legal_address_longitude, legal_address_latitude
		FROM projections.organizations WHERE id = ?`, orgID).
		Row().Scan(&addrText, &lon, &lat))
	require.Equal(t, "2 Different Avenue", addrText)
	require.NotNil(t, lon)
	require.InDelta(t, 12.5, *lon, 0.0001)
	require.NotNil(t, lat)
	require.InDelta(t, -7.25, *lat, 0.0001)
}
