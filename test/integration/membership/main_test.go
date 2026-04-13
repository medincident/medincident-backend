//go:build integration

package membership_integration_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-command-service/internal/services/membership"
	"github.com/medincident/medincident-command-service/internal/zitadel"
)

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()

	empSvc *membership.EmployeeService

	// Known test user IDs that testVerifier reports as existing.
	// Anything else is reported as not found.
	testUserAlice = "test-user-alice"
	testUserBob   = "test-user-bob"
	testUserCarol = "test-user-carol"
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

	verifier := newTestVerifier(testUserAlice, testUserBob, testUserCarol)
	empSvc = membership.NewEmployeeService(testDB, verifier, &testLogger)

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
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// testVerifier is an in-memory UserVerifier that reports a fixed set
// of user IDs as existing and returns zitadel.ErrUserNotFound for
// every other ID. Safe for concurrent use.
type testVerifier struct {
	known map[string]struct{}
}

func newTestVerifier(knownIDs ...string) *testVerifier {
	m := make(map[string]struct{}, len(knownIDs))
	for _, id := range knownIDs {
		m[id] = struct{}{}
	}
	return &testVerifier{known: m}
}

func (v *testVerifier) Verify(_ context.Context, zitadelUserID string) error {
	if _, ok := v.known[zitadelUserID]; ok {
		return nil
	}
	return zitadel.ErrUserNotFound
}

// Compile-time assertion.
var _ zitadel.UserVerifier = (*testVerifier)(nil)

// Silence unused-import warning when errors package isn't otherwise touched.
var _ = errors.New
