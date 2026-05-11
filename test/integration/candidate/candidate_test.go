//go:build integration

package candidate_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/query/membership"
)

// ---------------------------------------------------------------------------
// Seed helpers
// ---------------------------------------------------------------------------

// seedUser inserts a row into projections.users.
func seedUser(t *testing.T, id, firstName, lastName, email string, updatedAt time.Time) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO projections.users (id, user_name, first_name, last_name, display_name, email, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, email, firstName, lastName, firstName+" "+lastName, email, updatedAt,
	).Error
	require.NoError(t, err)
}

// seedEmployee inserts into projections.employee_cards.
// clinicID is *string (nil = no clinic assigned).
func seedEmployee(
	t *testing.T,
	employeeID, zitadelUserID, orgID, deptID string,
	clinicID *string,
	updatedAt time.Time,
) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO projections.employee_cards
		 (employee_id, zitadel_user_id, organization_id, department_id, clinic_id, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		employeeID, zitadelUserID, orgID, deptID, clinicID, updatedAt,
	).Error
	require.NoError(t, err)
}

// seedOrgAdmin inserts into projections.org_admins.
func seedOrgAdmin(t *testing.T, employeeID, orgID string) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO projections.org_admins (organization_id, employee_id) VALUES (?, ?)`,
		orgID, employeeID,
	).Error
	require.NoError(t, err)
}

// seedOrgHead inserts into projections.org_heads.
func seedOrgHead(t *testing.T, employeeID, orgID string) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO projections.org_heads (organization_id, employee_id) VALUES (?, ?)`,
		orgID, employeeID,
	).Error
	require.NoError(t, err)
}

// seedOrgDispatcher inserts into projections.org_dispatchers.
func seedOrgDispatcher(t *testing.T, employeeID, orgID string) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO projections.org_dispatchers (organization_id, employee_id) VALUES (?, ?)`,
		orgID, employeeID,
	).Error
	require.NoError(t, err)
}

// seedClinicHead inserts into projections.clinic_heads.
func seedClinicHead(t *testing.T, employeeID, clinicID string) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO projections.clinic_heads (clinic_id, employee_id) VALUES (?, ?)`,
		clinicID, employeeID,
	).Error
	require.NoError(t, err)
}

// seedDeptResponsible inserts into projections.department_responsibles.
func seedDeptResponsible(t *testing.T, employeeID, deptID string) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO projections.department_responsibles (department_id, employee_id) VALUES (?, ?)`,
		deptID, employeeID,
	).Error
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// Domain fixtures for authz
// ---------------------------------------------------------------------------

// domainOrg inserts a minimal org into domain.organizations.
func domainOrg(t *testing.T, orgID uuid.UUID) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO domain.organizations (id, name, legal_address) VALUES (?, ?, ROW(?, NULL)::domain.address)`,
		orgID, "Test Org "+orgID.String()[:8], "addr",
	).Error
	require.NoError(t, err)
}

// domainClinic inserts a minimal clinic into domain.clinics.
func domainClinic(t *testing.T, clinicID, orgID uuid.UUID) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`,
		clinicID, orgID, "Clinic "+clinicID.String()[:8], "addr",
	).Error
	require.NoError(t, err)
}

// domainDepartment inserts a minimal department into domain.departments.
func domainDepartment(t *testing.T, deptID, clinicID uuid.UUID) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`,
		deptID, clinicID, "Dept "+deptID.String()[:8],
	).Error
	require.NoError(t, err)
}

// domainEmployee inserts an employee into domain.employees.
func domainEmployee(t *testing.T, empID uuid.UUID, zitadelUserID string, orgID, deptID uuid.UUID) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empID, zitadelUserID, orgID, deptID,
	).Error
	require.NoError(t, err)
}

// domainOrgAdmin inserts an org admin row into domain.org_admins.
func domainOrgAdmin(t *testing.T, empID, orgID uuid.UUID) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO domain.org_admins (organization_id, employee_id) VALUES (?, ?)`,
		orgID, empID,
	).Error
	require.NoError(t, err)
}

// domainSystemAdmin inserts a system admin row into domain.system_admins.
func domainSystemAdmin(t *testing.T, zitadelUserID string) {
	t.Helper()
	err := testDB.Exec(
		`INSERT INTO domain.system_admins (zitadel_user_id) VALUES (?)`,
		zitadelUserID,
	).Error
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// Caller helpers
// ---------------------------------------------------------------------------

// systemAdminCaller creates a Caller who passes SystemAdmin checks by
// inserting zitadel_user_id into domain.system_admins.
func systemAdminCaller(t *testing.T, userID string) authz.Caller {
	t.Helper()
	domainSystemAdmin(t, userID)
	return authz.Caller{ZitadelUserID: userID}
}

// orgAdminCaller creates a Caller who passes AdminOf.Organization checks.
// It inserts domain fixtures so the authz check succeeds:
//   - an employee record for the caller (in orgID/deptID)
//   - an org_admin row
func orgAdminCaller(t *testing.T, callerZitadelID string, orgID, clinicID, deptID uuid.UUID) authz.Caller {
	t.Helper()
	empID := uuid.Must(uuid.NewV7())
	domainEmployee(t, empID, callerZitadelID, orgID, deptID)
	domainOrgAdmin(t, empID, orgID)
	return authz.Caller{ZitadelUserID: callerZitadelID}
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func newReader() *membership.CandidateReader {
	log := zerolog.Nop()
	return membership.NewCandidateReader(testDB, authzSvc, &log)
}

func containsZitadelUserID(views []membership.ZitadelUserView, id string) bool {
	for _, v := range views {
		if v.ZitadelUserID == id {
			return true
		}
	}
	return false
}

func containsEmployeeID(views []membership.EmployeeCardView, id uuid.UUID) bool {
	for _, v := range views {
		if v.EmployeeID == id {
			return true
		}
	}
	return false
}

// ctxT returns a context tied to the test deadline (10s).
func ctxT(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestForHire_ExcludesActiveEmployeeOfSameOrg ensures that a user who is an
// active employee of orgA is not returned as a candidate for orgA, while a
// user who has no employment is returned.
func TestForHire_ExcludesActiveEmployeeOfSameOrg(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	domainOrg(t, orgID)
	domainClinic(t, clinicID, orgID)
	domainDepartment(t, deptID, clinicID)

	// callerID is "admin-user" — an org admin.
	caller := orgAdminCaller(t, "admin-user", orgID, clinicID, deptID)

	// alice is already hired in orgA.
	aliceZID := "alice-zid"
	aliceEmpID := uuid.Must(uuid.NewV7()).String()
	seedUser(t, aliceZID, "Alice", "Smith", "alice@example.com", now)
	seedEmployee(t, aliceEmpID, aliceZID, orgID.String(), deptID.String(), nil, now)

	// bob has no employment anywhere.
	bobZID := "bob-zid"
	seedUser(t, bobZID, "Bob", "Jones", "bob@example.com", now.Add(-time.Second))

	r := newReader()
	results, _, err := r.ForHire(ctx, caller, orgID, "", "", 50)
	require.NoError(t, err)

	assert.False(t, containsZitadelUserID(results, aliceZID), "alice (active in org) must be excluded")
	assert.True(t, containsZitadelUserID(results, bobZID), "bob (no employment) must be included")
}

// TestForHire_IncludesEmployeeOfDifferentOrg ensures that a user who is an
// active employee of orgB appears as a hire candidate for orgA.
func TestForHire_IncludesEmployeeOfDifferentOrg(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgA := uuid.Must(uuid.NewV7())
	orgB := uuid.Must(uuid.NewV7())
	clinicA := uuid.Must(uuid.NewV7())
	clinicB := uuid.Must(uuid.NewV7())
	deptA := uuid.Must(uuid.NewV7())
	deptB := uuid.Must(uuid.NewV7())

	domainOrg(t, orgA)
	domainOrg(t, orgB)
	domainClinic(t, clinicA, orgA)
	domainClinic(t, clinicB, orgB)
	domainDepartment(t, deptA, clinicA)
	domainDepartment(t, deptB, clinicB)

	caller := orgAdminCaller(t, "admin-a", orgA, clinicA, deptA)

	// carol is employed in orgB only.
	carolZID := "carol-zid"
	carolEmpID := uuid.Must(uuid.NewV7()).String()
	seedUser(t, carolZID, "Carol", "Brown", "carol@example.com", now)
	seedEmployee(t, carolEmpID, carolZID, orgB.String(), deptB.String(), nil, now)

	r := newReader()
	results, _, err := r.ForHire(ctx, caller, orgA, "", "", 50)
	require.NoError(t, err)

	assert.True(t, containsZitadelUserID(results, carolZID), "carol (employed in orgB) must appear as candidate for orgA")
}

// TestForHire_SearchFilters verifies that passing searchQuery "diana" returns
// only the matching user and not others.
func TestForHire_SearchFilters(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	domainOrg(t, orgID)
	domainClinic(t, clinicID, orgID)
	domainDepartment(t, deptID, clinicID)

	caller := orgAdminCaller(t, "admin-search", orgID, clinicID, deptID)

	dianaZID := "diana-zid"
	clarkZID := "clark-zid"
	seedUser(t, dianaZID, "Diana", "Prince", "diana@example.com", now)
	seedUser(t, clarkZID, "Clark", "Kent", "clark@example.com", now.Add(-time.Second))

	r := newReader()
	results, _, err := r.ForHire(ctx, caller, orgID, "diana", "", 50)
	require.NoError(t, err)

	assert.True(t, containsZitadelUserID(results, dianaZID), "diana must be returned by search filter")
	assert.False(t, containsZitadelUserID(results, clarkZID), "clark must not appear in 'diana' search")
}

// TestForHire_CursorPagination seeds 3 users, fetches limit=2, then fetches
// the next page with the returned cursor and verifies correct non-overlapping
// pagination.
func TestForHire_CursorPagination(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	domainOrg(t, orgID)
	domainClinic(t, clinicID, orgID)
	domainDepartment(t, deptID, clinicID)

	caller := orgAdminCaller(t, "admin-page", orgID, clinicID, deptID)

	// Seed 3 users with distinct updated_at so ordering is deterministic.
	seedUser(t, "page-user-1", "User", "One", "u1@example.com", now)
	seedUser(t, "page-user-2", "User", "Two", "u2@example.com", now.Add(-time.Second))
	seedUser(t, "page-user-3", "User", "Three", "u3@example.com", now.Add(-2*time.Second))

	r := newReader()

	// Page 1: fetch 2.
	page1, nextCursor, err := r.ForHire(ctx, caller, orgID, "", "", 2)
	require.NoError(t, err)
	assert.Len(t, page1, 2, "page 1 should have 2 items")
	assert.NotEmpty(t, nextCursor, "next cursor must be set after page 1")

	// Page 2: fetch with cursor.
	page2, nextCursor2, err := r.ForHire(ctx, caller, orgID, "", nextCursor, 2)
	require.NoError(t, err)
	assert.Len(t, page2, 1, "page 2 should have 1 remaining item")
	assert.Empty(t, nextCursor2, "next cursor must be empty on last page")

	// No overlap between pages.
	page1IDs := map[string]bool{}
	for _, v := range page1 {
		page1IDs[v.ZitadelUserID] = true
	}
	for _, v := range page2 {
		assert.False(t, page1IDs[v.ZitadelUserID], "page 2 item %s must not overlap with page 1", v.ZitadelUserID)
	}
}

// TestForSystemAdmin_ExcludesExistingSystemAdmin verifies that a user who is
// already a system admin is excluded from ForSystemAdmin results.
func TestForSystemAdmin_ExcludesExistingSystemAdmin(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	// Caller is a system admin.
	caller := systemAdminCaller(t, "sys-caller")

	// frank is already a system admin in projections.system_admins.
	frankZID := "frank-zid"
	seedUser(t, frankZID, "Frank", "Admin", "frank@example.com", now)
	err := testDB.Exec(
		`INSERT INTO projections.system_admins (zitadel_user_id) VALUES (?)`, frankZID,
	).Error
	require.NoError(t, err)

	// grace is a regular user, not a system admin.
	graceZID := "grace-zid"
	seedUser(t, graceZID, "Grace", "Hopper", "grace@example.com", now.Add(-time.Second))

	r := newReader()
	results, _, err := r.ForSystemAdmin(ctx, caller, "", "", 50)
	require.NoError(t, err)

	assert.False(t, containsZitadelUserID(results, frankZID), "frank (already sysadmin) must be excluded")
	assert.True(t, containsZitadelUserID(results, graceZID), "grace (not sysadmin) must be included")
}

// TestForOrgAdmin_ExcludesAlreadyAssigned verifies that an employee who is
// already an org admin is excluded from ForOrgAdmin results.
func TestForOrgAdmin_ExcludesAlreadyAssigned(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	domainOrg(t, orgID)
	domainClinic(t, clinicID, orgID)
	domainDepartment(t, deptID, clinicID)

	caller := orgAdminCaller(t, "admin-org", orgID, clinicID, deptID)

	// emp1 is already an org admin (in projections).
	emp1ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp1ID.String(), "emp1-zid", orgID.String(), deptID.String(), nil, now)
	seedOrgAdmin(t, emp1ID.String(), orgID.String())

	// emp2 is a regular employee, not yet an org admin.
	emp2ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp2ID.String(), "emp2-zid", orgID.String(), deptID.String(), nil, now.Add(-time.Second))

	r := newReader()
	results, _, err := r.ForOrgAdmin(ctx, caller, orgID, "", "", 50)
	require.NoError(t, err)

	assert.False(t, containsEmployeeID(results, emp1ID), "emp1 (already org admin) must be excluded")
	assert.True(t, containsEmployeeID(results, emp2ID), "emp2 (not org admin) must be included")
}

// TestForClinicHead_OnlyReturnsEmployeesOfThatClinic verifies that only
// employees belonging to the target clinic are returned.
func TestForClinicHead_OnlyReturnsEmployeesOfThatClinic(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgID := uuid.Must(uuid.NewV7())
	clinicA := uuid.Must(uuid.NewV7())
	clinicB := uuid.Must(uuid.NewV7())
	deptA := uuid.Must(uuid.NewV7())
	deptB := uuid.Must(uuid.NewV7())

	domainOrg(t, orgID)
	domainClinic(t, clinicA, orgID)
	domainClinic(t, clinicB, orgID)
	domainDepartment(t, deptA, clinicA)
	domainDepartment(t, deptB, clinicB)

	// Caller is an org admin for orgID (which implies AdminOf.Clinic).
	caller := orgAdminCaller(t, "admin-clinic", orgID, clinicA, deptA)

	clinicAStr := clinicA.String()
	clinicBStr := clinicB.String()

	// emp1 belongs to clinicA.
	emp1ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp1ID.String(), "emp1-clinic-zid", orgID.String(), deptA.String(), &clinicAStr, now)

	// emp2 belongs to clinicB.
	emp2ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp2ID.String(), "emp2-clinic-zid", orgID.String(), deptB.String(), &clinicBStr, now.Add(-time.Second))

	// emp3 belongs to clinicA but is already a clinic head.
	emp3ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp3ID.String(), "emp3-clinic-zid", orgID.String(), deptA.String(), &clinicAStr, now.Add(-2*time.Second))
	seedClinicHead(t, emp3ID.String(), clinicA.String())

	r := newReader()
	results, _, err := r.ForClinicHead(ctx, caller, clinicA, "", "", 50)
	require.NoError(t, err)

	assert.True(t, containsEmployeeID(results, emp1ID), "emp1 (in clinicA) must be included")
	assert.False(t, containsEmployeeID(results, emp2ID), "emp2 (in clinicB) must be excluded")
	assert.False(t, containsEmployeeID(results, emp3ID), "emp3 (already clinic head) must be excluded")
}

// TestForDeptResponsible_OnlyReturnsEmployeesOfThatDept verifies that only
// employees belonging to the target department are returned.
func TestForDeptResponsible_OnlyReturnsEmployeesOfThatDept(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptA := uuid.Must(uuid.NewV7())
	deptB := uuid.Must(uuid.NewV7())

	domainOrg(t, orgID)
	domainClinic(t, clinicID, orgID)
	domainDepartment(t, deptA, clinicID)
	domainDepartment(t, deptB, clinicID)

	// Caller is an org admin for orgID (which implies AdminOf.Department).
	caller := orgAdminCaller(t, "admin-dept", orgID, clinicID, deptA)

	clinicStr := clinicID.String()

	// emp1 belongs to deptA.
	emp1ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp1ID.String(), "emp1-dept-zid", orgID.String(), deptA.String(), &clinicStr, now)

	// emp2 belongs to deptB.
	emp2ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp2ID.String(), "emp2-dept-zid", orgID.String(), deptB.String(), &clinicStr, now.Add(-time.Second))

	// emp3 belongs to deptA but is already a dept responsible.
	emp3ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp3ID.String(), "emp3-dept-zid", orgID.String(), deptA.String(), &clinicStr, now.Add(-2*time.Second))
	seedDeptResponsible(t, emp3ID.String(), deptA.String())

	r := newReader()
	results, _, err := r.ForDeptResponsible(ctx, caller, deptA, "", "", 50)
	require.NoError(t, err)

	assert.True(t, containsEmployeeID(results, emp1ID), "emp1 (in deptA) must be included")
	assert.False(t, containsEmployeeID(results, emp2ID), "emp2 (in deptB) must be excluded")
	assert.False(t, containsEmployeeID(results, emp3ID), "emp3 (already dept responsible) must be excluded")
}

// TestForOrgHead_ExcludesAlreadyAssigned verifies that an employee who is
// already an org head is excluded from ForOrgHead results.
func TestForOrgHead_ExcludesAlreadyAssigned(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	domainOrg(t, orgID)
	domainClinic(t, clinicID, orgID)
	domainDepartment(t, deptID, clinicID)

	caller := orgAdminCaller(t, "admin-orghead", orgID, clinicID, deptID)

	// emp1 is already an org head (in projections).
	emp1ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp1ID.String(), "emp1-orghead-zid", orgID.String(), deptID.String(), nil, now)
	seedOrgHead(t, emp1ID.String(), orgID.String())

	// emp2 is a regular employee, not yet an org head.
	emp2ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp2ID.String(), "emp2-orghead-zid", orgID.String(), deptID.String(), nil, now.Add(-time.Second))

	r := newReader()
	results, _, err := r.ForOrgHead(ctx, caller, orgID, "", "", 50)
	require.NoError(t, err)

	assert.False(t, containsEmployeeID(results, emp1ID), "emp1 (already org head) must be excluded")
	assert.True(t, containsEmployeeID(results, emp2ID), "emp2 (not org head) must be included")
}

// TestForOrgDispatcher_ExcludesAlreadyAssigned verifies that an employee who is
// already an org dispatcher is excluded from ForOrgDispatcher results.
func TestForOrgDispatcher_ExcludesAlreadyAssigned(t *testing.T) {
	resetAll(t)
	ctx := ctxT(t)
	now := time.Now().UTC().Truncate(time.Second)

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	domainOrg(t, orgID)
	domainClinic(t, clinicID, orgID)
	domainDepartment(t, deptID, clinicID)

	caller := orgAdminCaller(t, "admin-orgdisp", orgID, clinicID, deptID)

	// emp1 is already an org dispatcher (in projections).
	emp1ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp1ID.String(), "emp1-orgdisp-zid", orgID.String(), deptID.String(), nil, now)
	seedOrgDispatcher(t, emp1ID.String(), orgID.String())

	// emp2 is a regular employee, not yet an org dispatcher.
	emp2ID := uuid.Must(uuid.NewV7())
	seedEmployee(t, emp2ID.String(), "emp2-orgdisp-zid", orgID.String(), deptID.String(), nil, now.Add(-time.Second))

	r := newReader()
	results, _, err := r.ForOrgDispatcher(ctx, caller, orgID, "", "", 50)
	require.NoError(t, err)

	assert.False(t, containsEmployeeID(results, emp1ID), "emp1 (already org dispatcher) must be excluded")
	assert.True(t, containsEmployeeID(results, emp2ID), "emp2 (not org dispatcher) must be included")
}
