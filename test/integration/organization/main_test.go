//go:build integration

package organization_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// testPool is the shared pgxpool all test functions in this package use.
// Set up once in TestMain: a single Postgres 16 testcontainer is booted,
// migrations are applied via dbmate, and the pool is pointed at the
// container for the rest of the run.
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	setupCtx, setupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer setupCancel()

	container, err := tcpg.Run(setupCtx,
		"postgres:16",
		tcpg.WithDatabase("medincident"),
		tcpg.WithUsername("medincident"),
		tcpg.WithPassword("medincident"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(90*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "testcontainers postgres run: %v\n", err)
		return 1
	}
	defer func() {
		// Teardown runs on a fresh context so it is not affected by
		// setupCtx having fired (2m) during a long integration run.
		teardownCtx, teardownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer teardownCancel()
		_ = container.Terminate(teardownCtx)
	}()

	dsn, err := container.ConnectionString(setupCtx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "container ConnectionString: %v\n", err)
		return 1
	}

	if err := applyMigrations(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "applyMigrations: %v\n", err)
		return 1
	}

	testPool, err = pgxpool.New(setupCtx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pgxpool.New: %v\n", err)
		return 1
	}
	defer testPool.Close()

	return m.Run()
}

// applyMigrations runs dbmate up against the given DSN. dbmate is
// available via 'go tool dbmate' (pinned in go.mod's tool directive),
// so we invoke the Go toolchain directly instead of expecting the binary
// on PATH.
func applyMigrations(dsn string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}
	cmd := exec.Command(
		"go", "tool", "dbmate",
		"--migrations-dir", filepath.Join(repoRoot, "db", "migrations"),
		"up",
	)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "DATABASE_URL="+dsn)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// findRepoRoot walks up from the current working directory until it
// finds a go.mod, so the test does not need to know its own depth.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", dir)
		}
		dir = parent
	}
}
