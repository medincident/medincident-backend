//go:build integration

package requestclassifier_integration_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
)

const sysadminZitadelID = "sysadmin"

var (
	sysadminCaller = authz.Caller{ZitadelUserID: sysadminZitadelID}
	testDB         *gorm.DB
	testLogger     = zerolog.Nop()
	authzSvc       *authz.Authz
	typeSvc        *classifiersvc.RequestTypeService
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
				WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: %v\n", err)
		os.Exit(1)
	}
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get dsn: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "tool", "dbmate",
		"--url", dsn, "--migrations-dir", "../../../db/migrations",
		"--no-dump-schema", "up",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}

	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "gorm open: %v\n", err)
		os.Exit(1)
	}

	authzSvc = authz.New(testDB)
	typeSvc = classifiersvc.NewRequestTypeService(testDB, authzSvc, &testLogger)
	code := m.Run()
	if err := container.Terminate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "terminate postgres: %v\n", err)
	}
	os.Exit(code)
}

func resetDB(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	for _, q := range []string{
		`TRUNCATE TABLE domain.system_admins CASCADE`,
		`TRUNCATE TABLE domain.request_types CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
		`TRUNCATE TABLE projections.request_types CASCADE`,
	} {
		if _, err := raw.Exec(q); err != nil {
			t.Fatalf("truncate %q: %v", q, err)
		}
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1)`, sysadminZitadelID,
	); err != nil {
		t.Fatalf("seed sysadmin: %v", err)
	}
}

func insertOrg(t *testing.T) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	require.NoError(t, err)
	require.NoError(t, testDB.Create(&model.Organization{
		ID:           id,
		Name:         "Org",
		LegalAddress: model.Address{Text: "г. Тест, ул. Тестовая"},
	}).Error)
	return id
}

func loadType(t *testing.T, id uuid.UUID) model.RequestType {
	t.Helper()
	var row model.RequestType
	require.NoError(t, testDB.First(&row, "id = ?", id).Error)
	return row
}

func codeOf(t *testing.T, err error) string {
	t.Helper()
	var oe oops.OopsError
	if !errors.As(err, &oe) {
		t.Fatalf("expected oops error, got %T: %v", err, err)
	}
	code, _ := oe.Code().(string)
	return code
}
