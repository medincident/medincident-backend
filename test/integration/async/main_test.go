//go:build integration

// Package async_integration_test validates the full async CQRS pipeline:
// command service → outbox.events → publisher → NATS → domain consumer →
// projections on a separate query DB. Uses three testcontainers: two
// Postgres instances (command DB, query DB) and one NATS with JetStream.
package async_integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	tcnats "github.com/testcontainers/testcontainers-go/modules/nats"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-backend/internal/config"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	domainread "github.com/medincident/medincident-backend/internal/service/query/domain"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
)

const sysadminZitadelID = "sysadmin"

var (
	sysadminCaller = authz.Caller{ZitadelUserID: sysadminZitadelID}
	testLogger     = zerolog.Nop()

	commandDB *gorm.DB
	queryDB   *gorm.DB

	authzSvc *authz.Authz
	orgSvc   *orgsvc.OrganizationService
	clinSvc  *orgsvc.ClinicService
	deptSvc  *orgsvc.DepartmentService

	domainConsumer *domainread.Consumer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// ── Command DB ─────────────────────────────────────────────────────
	cmdContainer, err := tcpostgres.Run(ctx, "postgres:17-alpine",
		tcpostgres.WithDatabase("medincident_cmd"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	must(err, "start command postgres")
	defer func() { _ = cmdContainer.Terminate(ctx) }()

	cmdDSN, err := cmdContainer.ConnectionString(ctx, "sslmode=disable")
	must(err, "get command DSN")
	must(runMigrations(cmdDSN), "migrate command DB")

	commandDB, err = gorm.Open(postgres.Open(cmdDSN), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Silent),
	})
	must(err, "open command gorm")

	// ── Query DB ────────────────────────────────────────────────────────
	qryContainer, err := tcpostgres.Run(ctx, "postgres:17-alpine",
		tcpostgres.WithDatabase("medincident_qry"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	must(err, "start query postgres")
	defer func() { _ = qryContainer.Terminate(ctx) }()

	qryDSN, err := qryContainer.ConnectionString(ctx, "sslmode=disable")
	must(err, "get query DSN")
	must(runMigrations(qryDSN), "migrate query DB")

	queryDB, err = gorm.Open(postgres.Open(qryDSN), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Silent),
	})
	must(err, "open query gorm")

	// ── NATS ────────────────────────────────────────────────────────────
	natsContainer, err := tcnats.Run(ctx, "nats:2-alpine",
		testcontainers.WithWaitStrategy(wait.ForLog("Server is ready").WithStartupTimeout(15*time.Second)),
	)
	must(err, "start nats")
	defer func() { _ = natsContainer.Terminate(ctx) }()

	natsURL, err := natsContainer.ConnectionString(ctx)
	must(err, "get nats URL")

	nc, err := natsgo.Connect(natsURL, natsgo.MaxReconnects(10))
	must(err, "connect nats")
	defer func() { _ = nc.Drain() }()

	js, err := jetstream.New(nc)
	must(err, "init jetstream")

	// Create stream — only in tests; prod stream is admin-provisioned.
	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "medincident_events",
		Subjects: []string{"medincident.event.>"},
		Storage:  jetstream.MemoryStorage,
	})
	must(err, "create test stream")

	// ── Seed sysadmin on command DB ─────────────────────────────────────
	rawCmd, err := commandDB.DB()
	must(err, "get raw command db")
	_, err = rawCmd.Exec(`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1)`, sysadminZitadelID)
	must(err, "seed sysadmin")

	// ── Wire services ───────────────────────────────────────────────────
	authzSvc = authz.New(commandDB)
	orgSvc = orgsvc.NewOrganizationService(commandDB, authzSvc, &testLogger)
	clinSvc = orgsvc.NewClinicService(commandDB, authzSvc, &testLogger)
	deptSvc = orgsvc.NewDepartmentService(commandDB, authzSvc, &testLogger)

	proj := qprojector.NewProjectors()
	natsCfg := &config.NATSConfig{
		URL:         natsURL,
		Stream:      "medincident_events",
		Subjects:    []string{"medincident.event.>"},
		DurableName: "test-domain-consumer",
	}
	domainConsumer = domainread.NewConsumer(js, natsCfg, queryDB, proj, &testLogger)
	must(domainConsumer.Start(ctx), "start domain consumer")

	// ── Publisher goroutine ─────────────────────────────────────────────
	// Inline minimal publisher: poll outbox.events on command DB and
	// publish to NATS. No advisory lock needed in single-process tests.
	go runTestPublisher(ctx, cmdDSN, js)

	os.Exit(m.Run())
}

// runMigrations applies dbmate migrations to the given DSN.
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

func must(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s: %v\n", msg, err)
		os.Exit(1)
	}
}

// resetDBs truncates domain + outbox tables on command DB and projection
// tables on query DB between tests.
func resetDBs(t *testing.T) {
	t.Helper()
	rawCmd, _ := commandDB.DB()
	for _, q := range []string{
		`TRUNCATE TABLE domain.departments CASCADE`,
		`TRUNCATE TABLE domain.clinics CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
		`TRUNCATE TABLE outbox.events CASCADE`,
	} {
		if _, err := rawCmd.Exec(q); err != nil {
			t.Fatalf("truncate command: %v", err)
		}
	}
	// Re-seed sysadmin (idempotent — TestMain already inserts it once).
	if _, err := rawCmd.Exec(`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1) ON CONFLICT DO NOTHING`, sysadminZitadelID); err != nil {
		t.Fatalf("seed sysadmin: %v", err)
	}

	rawQry, _ := queryDB.DB()
	if _, err := rawQry.Exec(`
		DO $$ DECLARE r RECORD; BEGIN
		  FOR r IN SELECT tablename FROM pg_tables WHERE schemaname = 'projections' LOOP
		    EXECUTE 'TRUNCATE TABLE projections.' || quote_ident(r.tablename) || ' CASCADE';
		  END LOOP;
		END $$`); err != nil {
		t.Fatalf("truncate query: %v", err)
	}
}

// requireProjected polls queryDB until the given SELECT count(*) query
// returns > 0, or fails the test after timeout.
func requireProjected(t *testing.T, q string, args ...any) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var n int
		if err := queryDB.Raw(q, args...).Scan(&n).Error; err == nil && n > 0 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("projection not found within timeout for query: %s args: %v", q, args)
}
