//go:build integration

package authz_integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

var (
	testDB   *gorm.DB
	authzSvc *authz.Authz
)

// Fixture IDs shared across tests.
var (
	orgA     uuid.UUID
	orgB     uuid.UUID
	clinicA1 uuid.UUID
	clinicB1 uuid.UUID
	deptA1a  uuid.UUID
	deptB1a  uuid.UUID

	empAlice uuid.UUID
	empBob   uuid.UUID
	empCarol uuid.UUID
	empDave  uuid.UUID

	bobVacationID uuid.UUID
	aliceVacID    uuid.UUID
	categoryA     uuid.UUID
	typeA         uuid.UUID
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:17-alpine",
		tcpostgres.WithDatabase("medincident"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start postgres: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = container.Terminate(ctx) }()

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get DSN: %v\n", err)
		os.Exit(1)
	}

	if err := runMigrations(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "failed to migrate: %v\n", err)
		os.Exit(1)
	}

	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open gorm: %v\n", err)
		os.Exit(1)
	}

	authzSvc = authz.New(testDB)

	os.Exit(m.Run())
}

func runMigrations(dsn string) error {
	cmd := exec.Command("go", "tool", "dbmate",
		"--url", dsn,
		"--migrations-dir", "../../../db/migrations",
		"--no-dump-schema",
		"up",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func resetDB(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
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
			t.Fatalf("truncate %q: %v", q, err)
		}
	}
}

func seedFixtures(t *testing.T) {
	t.Helper()
	must := func(err error) { require.NoError(t, err) }

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

	// Organizations.
	must(testDB.Exec(`INSERT INTO domain.organizations (id, name, legal_address) VALUES (?, ?, ROW(?, NULL)::domain.address)`,
		orgA, "Org A", "addr-a").Error)
	must(testDB.Exec(`INSERT INTO domain.organizations (id, name, legal_address) VALUES (?, ?, ROW(?, NULL)::domain.address)`,
		orgB, "Org B", "addr-b").Error)

	// Clinics.
	must(testDB.Exec(`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`,
		clinicA1, orgA, "Clinic A1", "caddr-a1").Error)
	must(testDB.Exec(`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`,
		clinicB1, orgB, "Clinic B1", "caddr-b1").Error)

	// Departments.
	must(testDB.Exec(`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`,
		deptA1a, clinicA1, "Dept A1a").Error)
	must(testDB.Exec(`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`,
		deptB1a, clinicB1, "Dept B1a").Error)

	// Employees.
	must(testDB.Exec(`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empAlice, "alice", orgA, deptA1a).Error)
	must(testDB.Exec(`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empBob, "bob", orgA, deptA1a).Error)
	must(testDB.Exec(`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empCarol, "carol", orgA, deptA1a).Error)
	must(testDB.Exec(`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES (?, ?, ?, ?)`,
		empDave, "dave", orgB, deptB1a).Error)

	// System admin.
	must(testDB.Exec(`INSERT INTO domain.system_admins (zitadel_user_id) VALUES (?)`, "sysadmin").Error)

	// Org admin: Bob is admin of OrgA, Carol is his deputy.
	must(testDB.Exec(`INSERT INTO domain.org_admins (organization_id, employee_id, deputy_employee_id) VALUES (?, ?, ?)`,
		orgA, empBob, empCarol).Error)

	// Active vacation for Bob (started 1h ago, no end).
	must(testDB.Exec(`INSERT INTO domain.employee_vacations (id, employee_id, starts_at, ends_at) VALUES (?, ?, now() - interval '1 hour', NULL)`,
		bobVacationID, empBob).Error)

	// Vacation for Alice (for RequireOrgAdminViaVacation tests).
	must(testDB.Exec(`INSERT INTO domain.employee_vacations (id, employee_id, starts_at, ends_at) VALUES (?, ?, now() - interval '2 hours', now() + interval '24 hours')`,
		aliceVacID, empAlice).Error)

	// Incident category in OrgA.
	must(testDB.Exec(`INSERT INTO domain.incident_categories (id, organization_id, name) VALUES (?, ?, ?)`,
		categoryA, orgA, "Category A").Error)

	// Incident type under the category.
	must(testDB.Exec(`INSERT INTO domain.incident_types (id, organization_id, category_id, name) VALUES (?, ?, ?, ?)`,
		typeA, orgA, categoryA, "Type A").Error)
}

func ctxT(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}
