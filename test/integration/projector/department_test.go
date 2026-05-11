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
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
)

// seedClinic is a test helper that creates an organization + clinic via
// the projector so department tests can attach to a real parent.
func seedClinic(t *testing.T, ctx context.Context, now time.Time) (*model.Organization, *model.Clinic) {
	t.Helper()
	org := seedOrganization(t, ctx, now)
	clinic := &model.Clinic{
		ID:              uuid.Must(uuid.NewV7()),
		OrganizationID:  org.ID,
		Name:            "Parent Clinic",
		PhysicalAddress: model.Address{Text: "1 Clinic Street"},
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
	return org, clinic
}

// TestDepartmentCreated_WritesRowCountersAndBumpsParents confirms the
// projector writes the department, its counters row, and bumps both
// clinic and organization counters.
func TestDepartmentCreated_WritesRowCountersAndBumpsParents(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org, clinic := seedClinic(t, ctx, now)

	dept := &model.Department{
		ID:          uuid.Must(uuid.NewV7()),
		ClinicID:    clinic.ID,
		Name:        "Dept A",
		Description: null.StringFrom("department description long enough"),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DepartmentCreated(tx, dept.ID.String(), dept.CreatedAt, &deptv1.DepartmentCreated{
			ClinicId:    dept.ClinicID.String(),
			Name:        dept.Name,
			Description: dept.Description.String,
			CreatedAt:   timestamppb.New(dept.CreatedAt),
		})
	}))

	var name string
	var desc *string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT name, description FROM projections.departments WHERE id = ?`, dept.ID).
		Row().Scan(&name, &desc))
	require.Equal(t, "Dept A", name)
	require.NotNil(t, desc)
	require.Equal(t, "department description long enough", *desc)

	var employees int64
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT employees_total FROM projections.department_counters WHERE department_id = ?`, dept.ID).
		Row().Scan(&employees))
	require.Equal(t, int64(0), employees)

	var clinicDeptsTotal int64
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT departments_total FROM projections.clinic_counters WHERE clinic_id = ?`, clinic.ID).
		Row().Scan(&clinicDeptsTotal))
	require.Equal(t, int64(1), clinicDeptsTotal)

	var orgDeptsTotal int64
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT departments_total FROM projections.organization_counters WHERE organization_id = ?`, org.ID).
		Row().Scan(&orgDeptsTotal))
	require.Equal(t, int64(1), orgDeptsTotal)
}

// TestDepartmentDetailsChanged_UpdatesRowAndEmployeeCards confirms the
// projector updates projections.departments and mirrors the
// department_name into projections.employee_cards.
func TestDepartmentDetailsChanged_UpdatesRowAndEmployeeCards(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org, clinic := seedClinic(t, ctx, now)

	dept := &model.Department{
		ID:          uuid.Must(uuid.NewV7()),
		ClinicID:    clinic.ID,
		Name:        "Old Dept",
		Description: null.StringFrom("old description long enough"),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DepartmentCreated(tx, dept.ID.String(), dept.CreatedAt, &deptv1.DepartmentCreated{
			ClinicId:    dept.ClinicID.String(),
			Name:        dept.Name,
			Description: dept.Description.String,
			CreatedAt:   timestamppb.New(dept.CreatedAt),
		})
	}))

	// Seed an employee_card pointing at this dept.
	require.NoError(t, testDB.WithContext(ctx).Exec(`
		INSERT INTO projections.employee_cards
		    (employee_id, zitadel_user_id, organization_id, clinic_id,
		     department_id, department_name, updated_at)
		VALUES (?, 'zit-1', ?, ?, ?, 'Old Dept', ?)`,
		uuid.Must(uuid.NewV7()), org.ID, clinic.ID, dept.ID, now,
	).Error)

	dept.Name = "New Dept"
	dept.Description = null.StringFrom("new description long enough")
	dept.UpdatedAt = now.Add(time.Hour)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DepartmentDetailsChanged(tx, dept.ID.String(), dept.UpdatedAt, &deptv1.DepartmentDetailsChanged{
			Name:        dept.Name,
			Description: dept.Description.String,
			UpdatedAt:   timestamppb.New(dept.UpdatedAt),
		})
	}))

	var name string
	var desc *string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT name, description FROM projections.departments WHERE id = ?`, dept.ID).
		Row().Scan(&name, &desc))
	require.Equal(t, "New Dept", name)
	require.Equal(t, "new description long enough", *desc)

	var cardDeptName string
	require.NoError(t, testDB.WithContext(ctx).
		Raw(`SELECT department_name FROM projections.employee_cards WHERE department_id = ?`, dept.ID).
		Row().Scan(&cardDeptName))
	require.Equal(t, "New Dept", cardDeptName)
}
