//go:build integration

package projector_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
)

// seedOrg inserts a minimal projections.organizations row and returns the ID.
func seedOrg(t *testing.T, id uuid.UUID) {
	t.Helper()
	now := time.Now().UTC()
	require.NoError(t, testDB.Exec(`
		INSERT INTO projections.organizations
		    (id, name, legal_address_text, is_active, created_at, updated_at)
		VALUES (?, 'Test Org', 'addr', TRUE, ?, ?)`,
		id, now, now,
	).Error)
	require.NoError(t, testDB.Exec(`
		INSERT INTO projections.organization_counters
		    (organization_id, employees_total, clinics_total, departments_total, updated_at)
		VALUES (?, 0, 0, 0, ?)`,
		id, now,
	).Error)
}

// seedClinicRow inserts a minimal projections.clinics row.
func seedClinicRow(t *testing.T, id, orgID uuid.UUID) {
	t.Helper()
	now := time.Now().UTC()
	require.NoError(t, testDB.Exec(`
		INSERT INTO projections.clinics
		    (id, organization_id, name, physical_address_text, is_active, created_at, updated_at)
		VALUES (?, ?, 'Test Clinic', 'addr', TRUE, ?, ?)`,
		id, orgID, now, now,
	).Error)
	require.NoError(t, testDB.Exec(`
		INSERT INTO projections.clinic_counters
		    (clinic_id, organization_id, employees_total, departments_total, updated_at)
		VALUES (?, ?, 0, 0, ?)`,
		id, orgID, now,
	).Error)
}

// seedDeptRow inserts a minimal projections.departments row.
func seedDeptRow(t *testing.T, id, clinicID uuid.UUID) {
	t.Helper()
	// Look up the org_id from the clinic row already inserted.
	var orgID uuid.UUID
	require.NoError(t, testDB.Raw(
		`SELECT organization_id FROM projections.clinics WHERE id = ?`, clinicID,
	).Row().Scan(&orgID))
	now := time.Now().UTC()
	require.NoError(t, testDB.Exec(`
		INSERT INTO projections.departments
		    (id, clinic_id, name, is_active, created_at, updated_at)
		VALUES (?, ?, 'Test Dept', TRUE, ?, ?)`,
		id, clinicID, now, now,
	).Error)
	require.NoError(t, testDB.Exec(`
		INSERT INTO projections.department_counters
		    (department_id, clinic_id, organization_id, employees_total, updated_at)
		VALUES (?, ?, ?, 0, ?)`,
		id, clinicID, orgID, now,
	).Error)
}

// seedEmployeeRow inserts a minimal projections.employees row.
func seedEmployeeRow(t *testing.T, empID, deptID, orgID uuid.UUID, zitadelID string) {
	t.Helper()
	now := time.Now().UTC()
	require.NoError(t, testDB.Exec(`
		INSERT INTO projections.employees
		    (id, zitadel_user_id, department_id, organization_id, is_active, hired_at, updated_at)
		VALUES (?, ?, ?, ?, TRUE, ?, ?)`,
		empID, zitadelID, deptID, orgID, now, now,
	).Error)
}

func isActiveInProjection(t *testing.T, table, id string) bool {
	t.Helper()
	var active bool
	require.NoError(t, testDB.Raw(
		"SELECT is_active FROM "+table+" WHERE id = ?", id,
	).Row().Scan(&active))
	return active
}

func rowExistsInProjection(t *testing.T, table, id string) bool {
	t.Helper()
	var count int
	require.NoError(t, testDB.Raw(
		"SELECT count(*) FROM "+table+" WHERE id = ?", id,
	).Scan(&count).Error)
	return count > 0
}

func TestOrganizationDeactivated_SetsIsActiveFalse(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	id := uuid.Must(uuid.NewV7())
	seedOrg(t, id)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrganizationDeactivated(tx, id.String(), now, &orgv1.OrganizationDeactivated{
			OrganizationId: id.String(),
			UpdatedAt:      timestamppb.New(now),
		})
	}))

	assert.False(t, isActiveInProjection(t, "projections.organizations", id.String()))
}

func TestOrganizationActivated_SetsIsActiveTrue(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	id := uuid.Must(uuid.NewV7())
	seedOrg(t, id)

	// First deactivate.
	now := time.Now().UTC()
	require.NoError(t, testDB.Exec(
		`UPDATE projections.organizations SET is_active = FALSE WHERE id = ?`, id,
	).Error)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrganizationActivated(tx, id.String(), now, &orgv1.OrganizationActivated{
			OrganizationId: id.String(),
			UpdatedAt:      timestamppb.New(now),
		})
	}))

	assert.True(t, isActiveInProjection(t, "projections.organizations", id.String()))
}

func TestOrganizationDeleted_RemovesRow(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	id := uuid.Must(uuid.NewV7())
	seedOrg(t, id)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrganizationDeleted(tx, id.String(), now, &orgv1.OrganizationDeleted{
			OrganizationId: id.String(),
		})
	}))

	assert.False(t, rowExistsInProjection(t, "projections.organizations", id.String()))
}

func TestClinicDeactivated_SetsIsActiveFalse(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := uuid.Must(uuid.NewV7())
	clinID := uuid.Must(uuid.NewV7())
	seedOrg(t, orgID)
	seedClinicRow(t, clinID, orgID)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicDeactivated(tx, clinID.String(), now, &clinicv1.ClinicDeactivated{
			ClinicId:  clinID.String(),
			UpdatedAt: timestamppb.New(now),
		})
	}))

	assert.False(t, isActiveInProjection(t, "projections.clinics", clinID.String()))
}

func TestClinicActivated_SetsIsActiveTrue(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := uuid.Must(uuid.NewV7())
	clinID := uuid.Must(uuid.NewV7())
	seedOrg(t, orgID)
	seedClinicRow(t, clinID, orgID)
	require.NoError(t, testDB.Exec(
		`UPDATE projections.clinics SET is_active = FALSE WHERE id = ?`, clinID,
	).Error)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicActivated(tx, clinID.String(), now, &clinicv1.ClinicActivated{
			ClinicId:  clinID.String(),
			UpdatedAt: timestamppb.New(now),
		})
	}))

	assert.True(t, isActiveInProjection(t, "projections.clinics", clinID.String()))
}

func TestClinicDeleted_RemovesRow(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := uuid.Must(uuid.NewV7())
	clinID := uuid.Must(uuid.NewV7())
	seedOrg(t, orgID)
	seedClinicRow(t, clinID, orgID)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicDeleted(tx, clinID.String(), now, &clinicv1.ClinicDeleted{
			ClinicId: clinID.String(),
		})
	}))

	assert.False(t, rowExistsInProjection(t, "projections.clinics", clinID.String()))
}

func TestDepartmentDeactivated_SetsIsActiveFalse(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := uuid.Must(uuid.NewV7())
	clinID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	seedOrg(t, orgID)
	seedClinicRow(t, clinID, orgID)
	seedDeptRow(t, deptID, clinID)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DepartmentDeactivated(tx, deptID.String(), now, &deptv1.DepartmentDeactivated{
			DepartmentId: deptID.String(),
			UpdatedAt:    timestamppb.New(now),
		})
	}))

	assert.False(t, isActiveInProjection(t, "projections.departments", deptID.String()))
}

func TestDepartmentActivated_SetsIsActiveTrue(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := uuid.Must(uuid.NewV7())
	clinID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	seedOrg(t, orgID)
	seedClinicRow(t, clinID, orgID)
	seedDeptRow(t, deptID, clinID)
	require.NoError(t, testDB.Exec(
		`UPDATE projections.departments SET is_active = FALSE WHERE id = ?`, deptID,
	).Error)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DepartmentActivated(tx, deptID.String(), now, &deptv1.DepartmentActivated{
			DepartmentId: deptID.String(),
			UpdatedAt:    timestamppb.New(now),
		})
	}))

	assert.True(t, isActiveInProjection(t, "projections.departments", deptID.String()))
}

func TestDepartmentDeleted_RemovesRow(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := uuid.Must(uuid.NewV7())
	clinID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	seedOrg(t, orgID)
	seedClinicRow(t, clinID, orgID)
	seedDeptRow(t, deptID, clinID)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DepartmentDeleted(tx, deptID.String(), now, &deptv1.DepartmentDeleted{
			DepartmentId: deptID.String(),
		})
	}))

	assert.False(t, rowExistsInProjection(t, "projections.departments", deptID.String()))
}

func TestEmployeeDeactivated_SetsIsActiveFalse(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := uuid.Must(uuid.NewV7())
	clinID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	empID := uuid.Must(uuid.NewV7())
	seedOrg(t, orgID)
	seedClinicRow(t, clinID, orgID)
	seedDeptRow(t, deptID, clinID)
	seedEmployeeRow(t, empID, deptID, orgID, "test-user-1")

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.EmployeeDeactivated(tx, empID.String(), now, &empv1.EmployeeDeactivated{
			EmployeeId: empID.String(),
			UpdatedAt:  timestamppb.New(now),
		})
	}))

	assert.False(t, isActiveInProjection(t, "projections.employees", empID.String()))
}

func TestEmployeeActivated_SetsIsActiveTrue(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := uuid.Must(uuid.NewV7())
	clinID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	empID := uuid.Must(uuid.NewV7())
	seedOrg(t, orgID)
	seedClinicRow(t, clinID, orgID)
	seedDeptRow(t, deptID, clinID)
	seedEmployeeRow(t, empID, deptID, orgID, "test-user-2")
	require.NoError(t, testDB.Exec(
		`UPDATE projections.employees SET is_active = FALSE WHERE id = ?`, empID,
	).Error)

	now := time.Now().UTC()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.EmployeeActivated(tx, empID.String(), now, &empv1.EmployeeActivated{
			EmployeeId: empID.String(),
			UpdatedAt:  timestamppb.New(now),
		})
	}))

	assert.True(t, isActiveInProjection(t, "projections.employees", empID.String()))
}
