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
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
)

// seedOrganization is a test helper that inserts a fresh organization
// via the projector so clinic tests can attach to a real parent.
func seedOrganization(t *testing.T, ctx context.Context, now time.Time) *model.Organization {
	t.Helper()
	org := &model.Organization{
		ID:           uuid.Must(uuid.NewV7()),
		Name:         "Parent Org",
		Description:  null.StringFrom("parent org description long enough"),
		LegalAddress: model.Address{Text: "1 Legal Street"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrganizationCreated(tx, org.ID.String(), org.CreatedAt, &orgv1.OrganizationCreated{
			Name:         org.Name,
			Description:  org.Description.String,
			LegalAddress: &orgv1.Address{Text: org.LegalAddress.Text},
			CreatedAt:    timestamppb.New(org.CreatedAt),
		})
	}))
	return org
}

// TestClinicCreated_WritesRowCountersAndBumpsOrg confirms the projector
// writes the clinic row, its counter row, and bumps the parent
// organization's clinics_total.
func TestClinicCreated_WritesRowCountersAndBumpsOrg(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org := seedOrganization(t, ctx, now)

	clinic := &model.Clinic{
		ID:              uuid.Must(uuid.NewV7()),
		OrganizationID:  org.ID,
		Name:            "Main Clinic",
		Description:     null.StringFrom("clinic description long enough"),
		PhysicalAddress: model.Address{Text: "2 Physical Street"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicCreated(tx, clinic.ID.String(), clinic.CreatedAt, &clinicv1.ClinicCreated{
			OrganizationId:  clinic.OrganizationID.String(),
			Name:            clinic.Name,
			Description:     clinic.Description.String,
			PhysicalAddress: &clinicv1.Address{Text: clinic.PhysicalAddress.Text},
			CreatedAt:       timestamppb.New(clinic.CreatedAt),
		})
	}))

	var name string
	var desc *string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT name, description FROM projections.clinics WHERE id = ?`, clinic.ID).
		Row().Scan(&name, &desc))
	require.Equal(t, "Main Clinic", name)
	require.NotNil(t, desc)
	require.Equal(t, "clinic description long enough", *desc)

	var employees, departments int64
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT employees_total, departments_total
		FROM projections.clinic_counters WHERE clinic_id = ?`, clinic.ID).
		Row().Scan(&employees, &departments))
	require.Equal(t, int64(0), employees)
	require.Equal(t, int64(0), departments)

	var clinicsTotal int64
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT clinics_total FROM projections.organization_counters WHERE organization_id = ?`, org.ID).
		Row().Scan(&clinicsTotal))
	require.Equal(t, int64(1), clinicsTotal)
}

// TestClinicDetailsChanged_UpdatesRowAndEmployeeCards confirms the
// projector updates projections.clinics and mirrors the clinic_name
// into projections.employee_cards.
func TestClinicDetailsChanged_UpdatesRowAndEmployeeCards(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org := seedOrganization(t, ctx, now)

	clinic := &model.Clinic{
		ID:              uuid.Must(uuid.NewV7()),
		OrganizationID:  org.ID,
		Name:            "Old Clinic",
		Description:     null.StringFrom("old description long enough"),
		PhysicalAddress: model.Address{Text: "2 Physical Street"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicCreated(tx, clinic.ID.String(), clinic.CreatedAt, &clinicv1.ClinicCreated{
			OrganizationId:  clinic.OrganizationID.String(),
			Name:            clinic.Name,
			Description:     clinic.Description.String,
			PhysicalAddress: &clinicv1.Address{Text: clinic.PhysicalAddress.Text},
			CreatedAt:       timestamppb.New(clinic.CreatedAt),
		})
	}))

	// Seed an employee_card row that should be re-mirrored on rename.
	require.NoError(t, testDB.WithContext(ctx).Exec(`
		INSERT INTO projections.employee_cards
		    (employee_id, zitadel_user_id, organization_id, organization_name,
		     clinic_id, clinic_name, department_id, updated_at)
		VALUES (?, 'zit-1', ?, 'Parent Org', ?, 'Old Clinic', ?, ?)`,
		uuid.Must(uuid.NewV7()), org.ID, clinic.ID, uuid.Must(uuid.NewV7()), now,
	).Error)

	clinic.Name = "New Clinic"
	clinic.Description = null.StringFrom("new description long enough")
	clinic.UpdatedAt = now.Add(time.Hour)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicDetailsChanged(tx, clinic.ID.String(), clinic.UpdatedAt, &clinicv1.ClinicDetailsChanged{
			Name:        clinic.Name,
			Description: clinic.Description.String,
			UpdatedAt:   timestamppb.New(clinic.UpdatedAt),
		})
	}))

	var name string
	var desc *string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT name, description FROM projections.clinics WHERE id = ?`, clinic.ID).
		Row().Scan(&name, &desc))
	require.Equal(t, "New Clinic", name)
	require.Equal(t, "new description long enough", *desc)

	var cardClinicName string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT clinic_name FROM projections.employee_cards WHERE clinic_id = ?`, clinic.ID).
		Row().Scan(&cardClinicName))
	require.Equal(t, "New Clinic", cardClinicName)
}

// TestClinicPhysicalAddressChanged_UpdatesRow confirms the projector
// updates the physical_address_* columns on projections.clinics.
func TestClinicPhysicalAddressChanged_UpdatesRow(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org := seedOrganization(t, ctx, now)

	clinic := &model.Clinic{
		ID:              uuid.Must(uuid.NewV7()),
		OrganizationID:  org.ID,
		Name:            "Clinic A",
		PhysicalAddress: model.Address{Text: "2 Physical Street"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicCreated(tx, clinic.ID.String(), clinic.CreatedAt, &clinicv1.ClinicCreated{
			OrganizationId:  clinic.OrganizationID.String(),
			Name:            clinic.Name,
			PhysicalAddress: &clinicv1.Address{Text: clinic.PhysicalAddress.Text},
			CreatedAt:       timestamppb.New(clinic.CreatedAt),
		})
	}))

	clinic.PhysicalAddress = model.Address{
		Text:  "3 New Avenue",
		Point: null.ValueFrom(model.Point{Longitude: 10.5, Latitude: -3.25}),
	}
	clinic.UpdatedAt = now.Add(time.Hour)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicPhysicalAddressChanged(tx, clinic.ID.String(), clinic.UpdatedAt, &clinicv1.ClinicPhysicalAddressChanged{
			PhysicalAddress: &clinicv1.Address{
				Text:  clinic.PhysicalAddress.Text,
				Point: &clinicv1.Point{Longitude: 10.5, Latitude: -3.25},
			},
			UpdatedAt: timestamppb.New(clinic.UpdatedAt),
		})
	}))

	var addrText string
	var lon, lat *float64
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT physical_address_text, physical_address_longitude, physical_address_latitude
		FROM projections.clinics WHERE id = ?`, clinic.ID).
		Row().Scan(&addrText, &lon, &lat))
	require.Equal(t, "3 New Avenue", addrText)
	require.NotNil(t, lon)
	require.InDelta(t, 10.5, *lon, 0.0001)
	require.NotNil(t, lat)
	require.InDelta(t, -3.25, *lat, 0.0001)
}
