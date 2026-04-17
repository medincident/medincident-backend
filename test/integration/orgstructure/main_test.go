//go:build integration

package orgstructure_integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	orgsvc "github.com/medincident/medincident-command-service/internal/service/command/orgstructure"
)

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()

	orgSvc  *orgsvc.OrganizationService
	clinSvc *orgsvc.ClinicService
	deptSvc *orgsvc.DepartmentService
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

	orgSvc = orgsvc.NewOrganizationService(testDB, &testLogger)
	clinSvc = orgsvc.NewClinicService(testDB, &testLogger)
	deptSvc = orgsvc.NewDepartmentService(testDB, &testLogger)

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

// resetDB truncates every domain.* and outbox.* table between tests.
func resetDB(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	truncate := []string{
		`TRUNCATE TABLE outbox.events RESTART IDENTITY`,
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

func countOrganizations(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.organizations`).Scan(&n).Error)
	return n
}

func countClinics(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.clinics`).Scan(&n).Error)
	return n
}

func countDepartments(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.departments`).Scan(&n).Error)
	return n
}

func countOutboxEvents(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM outbox.events`).Scan(&n).Error)
	return n
}
