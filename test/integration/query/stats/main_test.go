//go:build integration

package stats_query_integration_test

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
)

var testDB *gorm.DB

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

	cmd := exec.Command("go", "tool", "dbmate",
		"--url", dsn,
		"--migrations-dir", "../../../../db/migrations",
		"--no-dump-schema",
		"up",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
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

	os.Exit(m.Run())
}

func resetProjections(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`TRUNCATE TABLE projections.organization_counters,
		projections.clinic_counters,
		projections.department_counters,
		projections.employee_cards,
		projections.employee_vacations,
		projections.employees,
		projections.departments,
		projections.clinics,
		projections.organizations CASCADE`); err != nil {
		t.Fatalf("truncate projections: %v", err)
	}
}
