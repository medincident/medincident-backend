//go:build integration

package authz_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// resetDBBench mirrors resetDB but accepts *testing.B.
func resetDBBench(b *testing.B) {
	b.Helper()
	raw, err := testDB.DB()
	if err != nil {
		b.Fatalf("get raw db: %v", err)
	}
	truncate := []string{
		`TRUNCATE TABLE domain.incident_types CASCADE`,
		`TRUNCATE TABLE domain.incident_categories CASCADE`,
		`TRUNCATE TABLE domain.department_responsibles CASCADE`,
		`TRUNCATE TABLE domain.clinic_heads CASCADE`,
		`TRUNCATE TABLE domain.org_admins CASCADE`,
		`TRUNCATE TABLE domain.org_heads CASCADE`,
		`TRUNCATE TABLE domain.org_dispatchers CASCADE`,
		`TRUNCATE TABLE domain.system_admins CASCADE`,
		`TRUNCATE TABLE domain.employee_vacations CASCADE`,
		`TRUNCATE TABLE domain.employees CASCADE`,
		`TRUNCATE TABLE domain.departments CASCADE`,
		`TRUNCATE TABLE domain.clinics CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
	}
	for _, q := range truncate {
		if _, err := raw.Exec(q); err != nil {
			b.Fatalf("truncate %q: %v", q, err)
		}
	}
}

// seedFixturesBench mirrors seedFixtures but accepts *testing.B.
func seedFixturesBench(b *testing.B) {
	b.Helper()
	must := func(err error) {
		if err != nil {
			b.Fatalf("seed: %v", err)
		}
	}

	orgA = uuid.Must(uuid.NewV7())
	orgB = uuid.Must(uuid.NewV7())
	clinicA1 = uuid.Must(uuid.NewV7())
	clinicB1 = uuid.Must(uuid.NewV7())
	deptA1a = uuid.Must(uuid.NewV7())
	deptB1a = uuid.Must(uuid.NewV7())
	empAlice = uuid.Must(uuid.NewV7())
	empBob = uuid.Must(uuid.NewV7())
	empCarol = uuid.Must(uuid.NewV7())
	empDave = uuid.Must(uuid.NewV7())
	bobVacationID = uuid.Must(uuid.NewV7())
	aliceVacID = uuid.Must(uuid.NewV7())
	categoryA = uuid.Must(uuid.NewV7())
	typeA = uuid.Must(uuid.NewV7())

	must(testDB.Exec(`INSERT INTO domain.organizations (id, name, legal_address) VALUES (?, ?, ROW(?, NULL)::domain.address)`,
		orgA, "Org A", "addr-a").Error)
	must(testDB.Exec(`INSERT INTO domain.organizations (id, name, legal_address) VALUES (?, ?, ROW(?, NULL)::domain.address)`,
		orgB, "Org B", "addr-b").Error)
	must(testDB.Exec(`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`,
		clinicA1, orgA, "Clinic A1", "caddr-a1").Error)
	must(testDB.Exec(`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`,
		clinicB1, orgB, "Clinic B1", "caddr-b1").Error)
	must(testDB.Exec(`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`,
		deptA1a, clinicA1, "Dept A1a").Error)
	must(testDB.Exec(`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`,
		deptB1a, clinicB1, "Dept B1a").Error)
	must(testDB.Exec(`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empAlice, "alice", orgA, deptA1a).Error)
	must(testDB.Exec(`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empBob, "bob", orgA, deptA1a).Error)
	must(testDB.Exec(`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empCarol, "carol", orgA, deptA1a).Error)
	must(testDB.Exec(`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empDave, "dave", orgB, deptB1a).Error)
	must(testDB.Exec(`INSERT INTO domain.system_admins (zitadel_user_id) VALUES (?)`, "sysadmin").Error)
	must(testDB.Exec(`INSERT INTO domain.org_admins (organization_id, employee_id, deputy_employee_id) VALUES (?, ?, ?)`,
		orgA, empBob, empCarol).Error)
	must(testDB.Exec(`INSERT INTO domain.employee_vacations (id, employee_id, starts_at, ends_at) VALUES (?, ?, now() - interval '1 hour', NULL)`,
		bobVacationID, empBob).Error)
	must(testDB.Exec(`INSERT INTO domain.employee_vacations (id, employee_id, starts_at, ends_at) VALUES (?, ?, now() - interval '2 hours', now() + interval '24 hours')`,
		aliceVacID, empAlice).Error)
	must(testDB.Exec(`INSERT INTO domain.incident_categories (id, organization_id, name) VALUES (?, ?, ?)`,
		categoryA, orgA, "Category A").Error)
	must(testDB.Exec(`INSERT INTO domain.incident_types (id, organization_id, category_id, name) VALUES (?, ?, ?, ?)`,
		typeA, orgA, categoryA, "Type A").Error)
}

func BenchmarkRequire_SystemAdmin(b *testing.B) {
	resetDBBench(b)
	seedFixturesBench(b)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := authzSvc.Require(ctx, "sysadmin", authz.SystemAdmin); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkRequire_AdminOf_Clinic_OrgAdmin(b *testing.B) {
	resetDBBench(b)
	seedFixturesBench(b)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := authzSvc.Require(ctx, "bob", authz.AdminOf.Clinic(clinicA1)); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkRequire_AdminOf_Clinic_Denied measures the deny path —
// alice is neither sysadmin nor org-admin for clinicA1. A nil error
// would mean the seed is broken; we fail fast on that to keep the
// benchmark honest. Any other error (e.g. DB timeout) surfaces
// through b.Fatalf via the benchmark harness's own runtime errors.
func BenchmarkRequire_AdminOf_Clinic_Denied(b *testing.B) {
	resetDBBench(b)
	seedFixturesBench(b)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := authzSvc.Require(ctx, "alice", authz.AdminOf.Clinic(clinicA1)); err == nil {
			b.Fatal("expected permission denied, got nil")
		}
	}
}
