//go:build integration

package orgstructure_test

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
)

const sysadminZitadelID = "sysadmin"

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()

	authzSvc *authz.Authz
	deptSvc  *orgsvc.DepartmentService
	clinSvc  *orgsvc.ClinicService
	orgSvc   *orgsvc.OrganizationService

	sysadminCaller = authz.Caller{ZitadelUserID: sysadminZitadelID}
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
	orgSvc = orgsvc.NewOrganizationService(testDB, authzSvc, &testLogger)
	clinSvc = orgsvc.NewClinicService(testDB, authzSvc, &testLogger)
	deptSvc = orgsvc.NewDepartmentService(testDB, authzSvc, &testLogger)

	os.Exit(m.Run())
}

func runMigrations(dsn string) error {
	cmd := exec.Command("go", "tool", "dbmate",
		"--url", dsn,
		"--migrations-dir", "../../../../db/migrations",
		"--no-dump-schema",
		"up",
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// resetDB truncates all relevant domain and projection tables and
// re-seeds the system-admin caller used by every test.
func resetDB(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	truncate := []string{
		`TRUNCATE TABLE domain.system_admins CASCADE`,
		`TRUNCATE TABLE domain.departments CASCADE`,
		`TRUNCATE TABLE domain.clinics CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
		`TRUNCATE TABLE projections.organization_counters,
		                 projections.clinic_counters,
		                 projections.department_counters,
		                 projections.employee_cards,
		                 projections.employee_vacations,
		                 projections.employees,
		                 projections.departments,
		                 projections.clinics,
		                 projections.organizations,
		                 projections.clinic_heads,
		                 projections.department_responsibles,
		                 projections.org_admins,
		                 projections.org_dispatchers,
		                 projections.org_heads,
		                 projections.system_admins
		         CASCADE`,
	}
	for _, q := range truncate {
		if _, err := raw.Exec(q); err != nil {
			t.Fatalf("truncate: %v", err)
		}
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1)`,
		sysadminZitadelID,
	); err != nil {
		t.Fatalf("seed sysadmin: %v", err)
	}
}

// codeOf extracts the oops error Code as a string from any error.
func codeOf(t *testing.T, err error) string {
	t.Helper()
	var oe oops.OopsError
	if !errors.As(err, &oe) {
		t.Fatalf("expected oops error, got %T: %v", err, err)
	}
	code, ok := oe.Code().(string)
	if !ok {
		t.Fatalf("oops.Code() returned non-string: %T %v", oe.Code(), oe.Code())
	}
	return code
}

// seedClinic creates one organization + one clinic and returns the clinic ID.
func seedClinic(t *testing.T) uuid.UUID {
	t.Helper()
	orgRes, err := orgSvc.Create(context.Background(), orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:         "Test Organization",
			LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Тверская, д. 1"},
		},
	})
	require.NoError(t, err)

	clinRes, err := clinSvc.Create(context.Background(), orgsvc.CreateClinicCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{
			OrganizationID:  orgRes.ID.String(),
			Name:            "Test Clinic",
			PhysicalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Арбат, д. 5"},
		},
	})
	require.NoError(t, err)
	return clinRes.ID
}

// TestDepartmentCreate_ValidClinic verifies that creating a department
// under an existing clinic succeeds and persists both domain and
// projection rows.
func TestDepartmentCreate_ValidClinic(t *testing.T) {
	resetDB(t)
	clinicID := seedClinic(t)

	res, err := deptSvc.Create(context.Background(), orgsvc.CreateDepartmentCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{
			ClinicID: clinicID.String(),
			Name:     "Терапия",
		},
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, res.ID)

	var count int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.departments WHERE id = ?`, res.ID).Scan(&count).Error)
	assert.Equal(t, 1, count)

	var projCount int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM projections.departments WHERE id = ?`, res.ID).Scan(&projCount).Error)
	assert.Equal(t, 1, projCount)
}

// TestDepartmentCreate_NonExistentClinic verifies that creating a
// department under a clinic that does not exist returns the
// ErrCodeDepartmentClinicNotFound error code.
func TestDepartmentCreate_NonExistentClinic(t *testing.T) {
	resetDB(t)

	_, err := deptSvc.Create(context.Background(), orgsvc.CreateDepartmentCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{
			ClinicID: uuid.New().String(),
			Name:     "Orphan Department",
		},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeDepartmentClinicNotFound, codeOf(t, err))

	var count int
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.departments`).Scan(&count).Error)
	assert.Equal(t, 0, count)
}
