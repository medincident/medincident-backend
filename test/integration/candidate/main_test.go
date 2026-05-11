//go:build integration

package candidate_integration_test

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

// testDB is the package-shared *gorm.DB connected to the testcontainers Postgres.
// Each test resets state via resetAll.
var testDB *gorm.DB

// authzSvc is the shared authz.Authz instance backed by domain.* tables.
var authzSvc *authz.Authz

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

// runMigrations invokes dbmate via go tool against the given DSN.
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

// resetAll truncates every table touched by CandidateReader tests between
// tests so each test sees a clean slate.
func resetAll(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	stmts := []string{
		// domain tables (authz checks)
		`TRUNCATE TABLE domain.system_admins CASCADE`,
		`TRUNCATE TABLE domain.org_admins CASCADE`,
		`TRUNCATE TABLE domain.org_heads CASCADE`,
		`TRUNCATE TABLE domain.org_dispatchers CASCADE`,
		`TRUNCATE TABLE domain.clinic_heads CASCADE`,
		`TRUNCATE TABLE domain.department_responsibles CASCADE`,
		`TRUNCATE TABLE domain.employee_vacations CASCADE`,
		`TRUNCATE TABLE domain.employees CASCADE`,
		`TRUNCATE TABLE domain.departments CASCADE`,
		`TRUNCATE TABLE domain.clinics CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
		// projections tables (query results)
		`TRUNCATE TABLE projections.employee_cards CASCADE`,
		`TRUNCATE TABLE projections.org_admins CASCADE`,
		`TRUNCATE TABLE projections.org_heads CASCADE`,
		`TRUNCATE TABLE projections.org_dispatchers CASCADE`,
		`TRUNCATE TABLE projections.clinic_heads CASCADE`,
		`TRUNCATE TABLE projections.department_responsibles CASCADE`,
		`TRUNCATE TABLE projections.system_admins CASCADE`,
		`TRUNCATE TABLE projections.users CASCADE`,
	}
	for _, q := range stmts {
		if _, err := raw.Exec(q); err != nil {
			t.Fatalf("truncate %q: %v", q, err)
		}
	}
}
