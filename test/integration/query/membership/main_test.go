//go:build integration

package membership_query_integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// sysadminZitadelID is the Zitadel user ID of the system-admin seeded
// by resetProjections. Membership scoped reads run through authz, so
// the integration suite needs a real domain.system_admins row — not
// only the projection — to exercise SystemAdmin branches.
const sysadminZitadelID = "sysadmin"

var sysadminCaller = authz.Caller{ZitadelUserID: sysadminZitadelID}

var (
	testDB   *gorm.DB
	authzSvc *authz.Authz
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
		"--migrations-dir", "../../../../db/migrations",
		"--no-dump-schema",
		"up",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func resetProjections(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`TRUNCATE TABLE
		domain.employee_vacations,
		domain.employees,
		domain.system_admins,
		domain.departments,
		domain.clinics,
		domain.organizations,
		projections.organization_counters,
		projections.clinic_counters,
		projections.department_counters,
		projections.employee_cards,
		projections.employee_vacations,
		projections.employees,
		projections.departments,
		projections.clinics,
		projections.organizations,
		projections.sessions,
		projections.users,
		projections.clinic_heads,
		projections.department_responsibles,
		projections.org_admins,
		projections.org_dispatchers,
		projections.org_heads,
		projections.system_admins,
		projections.incident_categories,
		projections.incident_types CASCADE`); err != nil {
		t.Fatalf("truncate projections: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1)`,
		sysadminZitadelID,
	); err != nil {
		t.Fatalf("seed sysadmin: %v", err)
	}
}
