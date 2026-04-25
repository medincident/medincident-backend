//go:build integration

package identity_query_integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/testcontainers/testcontainers-go"
	tcnats "github.com/testcontainers/testcontainers-go/modules/nats"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const sysadminZitadelID = "sysadmin"

var (
	testDB  *gorm.DB
	natsURL string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pg, err := tcpostgres.Run(ctx,
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
	defer func() { _ = pg.Terminate(ctx) }()

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
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

	// The nats testcontainers module enables JetStream by default (-js
	// flag is baked into the command line).
	n, err := tcnats.Run(ctx, "nats:2.10-alpine")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start nats: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = n.Terminate(ctx) }()
	natsURL, err = n.ConnectionString(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get nats url: %v\n", err)
		os.Exit(1)
	}

	// Create the zitadel stream up-front so the consumer can subscribe.
	nc, err := nats.Connect(natsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect nats: %v\n", err)
		os.Exit(1)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init jetstream: %v\n", err)
		os.Exit(1)
	}
	if _, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "zitadel",
		Subjects:  []string{"zitadel.>"},
		Retention: jetstream.WorkQueuePolicy,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create stream: %v\n", err)
		os.Exit(1)
	}
	nc.Close()

	os.Exit(m.Run())
}

func resetProjections(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`TRUNCATE TABLE
		domain.system_admins,
		projections.users`); err != nil {
		t.Fatalf("truncate projections: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1)`,
		sysadminZitadelID,
	); err != nil {
		t.Fatalf("seed sysadmin: %v", err)
	}
}
