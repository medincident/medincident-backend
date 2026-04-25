//go:build integration

package membership_test

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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

const sysadminZitadelID = "sysadmin"

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()

	authzSvc *authz.Authz
	// empSvc is wired with a nil Zitadel verifier; ScheduleVacation does
	// not call the verifier so it is safe to omit for these tests.
	empSvc *membership.EmployeeService

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
	empSvc = membership.NewEmployeeService(testDB, authzSvc, nil, &testLogger)

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

// resetDB truncates membership-relevant tables and re-seeds the
// system-admin caller used by every test.
func resetDB(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	truncate := []string{
		`TRUNCATE TABLE domain.system_admins CASCADE`,
		`TRUNCATE TABLE domain.employee_vacations CASCADE`,
		`TRUNCATE TABLE domain.employees CASCADE`,
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

// seedEmployee inserts a minimal employee row directly into the DB,
// bypassing Zitadel verification. The authz check for the vacation
// commands short-circuits on SystemAdmin, so no projection rows are
// needed here. Returns the employee ID.
func seedEmployee(t *testing.T) uuid.UUID {
	t.Helper()

	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	empID := uuid.Must(uuid.NewV7())

	require.NoError(t, testDB.Exec(
		`INSERT INTO domain.organizations (id, name, legal_address) VALUES (?, ?, ROW(?, NULL)::domain.address)`,
		orgID, "Org", "addr",
	).Error)
	require.NoError(t, testDB.Exec(
		`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`,
		clinicID, orgID, "Clinic", "addr",
	).Error)
	require.NoError(t, testDB.Exec(
		`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`,
		deptID, clinicID, "Dept",
	).Error)

	emp := model.Employee{
		ID:             empID,
		ZitadelUserID:  fmt.Sprintf("zitadel-%s", empID),
		OrganizationID: orgID,
		DepartmentID:   deptID,
	}
	require.NoError(t, testDB.Create(&emp).Error)

	return empID
}

// oopsCode extracts the oops Code as a string from any error.
func oopsCode(t *testing.T, err error) string {
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

// TestScheduleVacation_FirstVacation verifies that the first vacation
// for a new employee is accepted without error.
func TestScheduleVacation_FirstVacation(t *testing.T) {
	resetDB(t)
	empID := seedEmployee(t)
	start := time.Now().Add(72 * time.Hour)

	res, err := empSvc.ScheduleVacation(context.Background(), membership.ScheduleVacationCommand{
		Caller: sysadminCaller,
		Payload: membership.ScheduleVacationPayload{
			EmployeeID: empID.String(),
			StartsAt:   start,
		},
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, res.ID)

	var count int
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.employee_vacations WHERE employee_id = ?`, empID,
	).Scan(&count).Error)
	assert.Equal(t, 1, count)
}

// TestScheduleVacation_NonOverlapping verifies that two non-overlapping
// vacations for the same employee are both accepted.
func TestScheduleVacation_NonOverlapping(t *testing.T) {
	resetDB(t)
	empID := seedEmployee(t)

	start1 := time.Now().Add(72 * time.Hour)
	end1 := start1.Add(48 * time.Hour)
	_, err := empSvc.ScheduleVacation(context.Background(), membership.ScheduleVacationCommand{
		Caller: sysadminCaller,
		Payload: membership.ScheduleVacationPayload{
			EmployeeID: empID.String(),
			StartsAt:   start1,
			EndsAt:     &end1,
		},
	})
	require.NoError(t, err)

	// Second vacation starts after the first one ends.
	start2 := end1.Add(24 * time.Hour)
	end2 := start2.Add(48 * time.Hour)
	_, err = empSvc.ScheduleVacation(context.Background(), membership.ScheduleVacationCommand{
		Caller: sysadminCaller,
		Payload: membership.ScheduleVacationPayload{
			EmployeeID: empID.String(),
			StartsAt:   start2,
			EndsAt:     &end2,
		},
	})
	require.NoError(t, err)

	var count int
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.employee_vacations WHERE employee_id = ?`, empID,
	).Scan(&count).Error)
	assert.Equal(t, 2, count)
}

// TestScheduleVacation_Overlap verifies that scheduling a vacation that
// overlaps an existing one for the same employee is rejected with
// ErrCodeVacationOverlap.
func TestScheduleVacation_Overlap(t *testing.T) {
	resetDB(t)
	empID := seedEmployee(t)

	start1 := time.Now().Add(72 * time.Hour)
	end1 := start1.Add(96 * time.Hour)
	_, err := empSvc.ScheduleVacation(context.Background(), membership.ScheduleVacationCommand{
		Caller: sysadminCaller,
		Payload: membership.ScheduleVacationPayload{
			EmployeeID: empID.String(),
			StartsAt:   start1,
			EndsAt:     &end1,
		},
	})
	require.NoError(t, err)

	// Overlapping: starts in the middle of the first vacation.
	start2 := start1.Add(24 * time.Hour)
	end2 := end1.Add(48 * time.Hour)
	_, err = empSvc.ScheduleVacation(context.Background(), membership.ScheduleVacationCommand{
		Caller: sysadminCaller,
		Payload: membership.ScheduleVacationPayload{
			EmployeeID: empID.String(),
			StartsAt:   start2,
			EndsAt:     &end2,
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationOverlap, oopsCode(t, err))
}
