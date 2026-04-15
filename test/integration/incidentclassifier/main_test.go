//go:build integration

package incidentclassifier_integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-command-service/internal/model"
	classifiersvc "github.com/medincident/medincident-command-service/internal/services/incident/classifier"
)

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()

	categorySvc *classifiersvc.IncidentCategoryService
	typeSvc     *classifiersvc.IncidentTypeService
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
		Logger:                 gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open gorm: %v\n", err)
		os.Exit(1)
	}

	categorySvc = classifiersvc.NewIncidentCategoryService(testDB, &testLogger)
	typeSvc = classifiersvc.NewIncidentTypeService(testDB, &testLogger)

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

// resetDB truncates every table the classifier suite touches.
func resetDB(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	truncate := []string{
		`TRUNCATE TABLE outbox.events RESTART IDENTITY`,
		`TRUNCATE TABLE domain.incident_types CASCADE`,
		`TRUNCATE TABLE domain.incident_categories CASCADE`,
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

// insertOrganization creates a minimal organisation row directly so that
// classifier rows (which have ON DELETE RESTRICT FK to organizations)
// have a valid parent to reference. Returns the new organisation id.
func insertOrganization(t *testing.T, name string) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	require.NoError(t, err)
	row := model.Organization{
		ID:   id,
		Name: name,
		LegalAddress: model.Address{
			Text: "г. Тест, ул. Тестовая",
		},
	}
	require.NoError(t, testDB.Create(&row).Error)
	return id
}

func countIncidentCategories(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.incident_categories`).Scan(&n).Error)
	return n
}

func countIncidentTypes(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.incident_types`).Scan(&n).Error)
	return n
}

func countOutboxEventsWithSubject(t *testing.T, subject string) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM outbox.events WHERE subject = ?`, subject,
	).Scan(&n).Error)
	return n
}

func loadCategory(t *testing.T, id uuid.UUID) model.IncidentCategory {
	t.Helper()
	var row model.IncidentCategory
	require.NoError(t, testDB.First(&row, "id = ?", id).Error)
	return row
}

func loadType(t *testing.T, id uuid.UUID) model.IncidentType {
	t.Helper()
	var row model.IncidentType
	require.NoError(t, testDB.First(&row, "id = ?", id).Error)
	return row
}

// nullUUIDPtr returns the uuid from a null.UUID or nil sentinel for
// unset values — used to make comparisons readable in tests.
var _ = null.String{} // keep import if future tests use it
