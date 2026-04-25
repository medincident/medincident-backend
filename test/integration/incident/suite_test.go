//go:build integration

package incident_integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	queryincident "github.com/medincident/medincident-backend/internal/service/query/incident"
	querybuffer "github.com/medincident/medincident-backend/internal/service/query/incident/buffer"
)

// sysadminZitadelID is the Zitadel user ID seeded into domain.system_admins
// and used as the privileged caller throughout the suite.
const sysadminZitadelID = "sysadmin"

var sysadminCaller = authz.Caller{ZitadelUserID: sysadminZitadelID}

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()

	authzSvc    *authz.Authz
	orgSvc      *orgsvc.OrganizationService
	clinSvc     *orgsvc.ClinicService
	deptSvc     *orgsvc.DepartmentService
	categorySvc *classifiersvc.IncidentCategoryService
	typeSvc     *classifiersvc.IncidentTypeService
	incidentSvc *incidentsvc.IncidentService
	bufferSvc   *buffersvc.BufferService
	incidentRdr *queryincident.Reader
	bufferRdr   *querybuffer.Reader
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
	categorySvc = classifiersvc.NewIncidentCategoryService(testDB, authzSvc, &testLogger)
	typeSvc = classifiersvc.NewIncidentTypeService(testDB, authzSvc, &testLogger)
	incidentSvc = incidentsvc.NewIncidentService(testDB, authzSvc, &testLogger)
	bufferSvc = buffersvc.NewBufferService(testDB, authzSvc, &testLogger)
	incidentRdr = queryincident.NewReader(testDB, &testLogger)
	bufferRdr = querybuffer.NewReader(testDB, &testLogger, incidentRdr)

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

// resetDB truncates all domain.* and projections.* tables and re-seeds
// the system-admin used by every test.
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
		`TRUNCATE TABLE domain.incident_types CASCADE`,
		`TRUNCATE TABLE domain.incident_categories CASCADE`,
		`TRUNCATE TABLE domain.incidents CASCADE`,
		`TRUNCATE TABLE domain.patient_incident_buffer CASCADE`,
		`TRUNCATE TABLE projections.incidents,
		                 projections.incident_status_history,
		                 projections.incident_priority_history,
		                 projections.patient_incident_buffer,
		                 projections.incident_categories,
		                 projections.incident_types,
		                 projections.organization_counters,
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
	// Re-seed system admin.
	if _, err := raw.Exec(
		`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1)`,
		sysadminZitadelID,
	); err != nil {
		t.Fatalf("seed sysadmin domain: %v", err)
	}
	seedUser(t, sysadminZitadelID, "System Admin")
}

// seedUser inserts a projections.users row. Required for any caller that
// triggers incident creation (display_name snapshot lookup).
func seedUser(t *testing.T, zitadelID, displayName string) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO projections.users (id, user_name, first_name, last_name, display_name, email, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, now(), now())
		 ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name`,
		zitadelID, zitadelID, "Test", "User", displayName, zitadelID+"@test.local",
	); err != nil {
		t.Fatalf("seed projections.users %s: %v", zitadelID, err)
	}
}

// seedOrg creates a test organization and returns its ID.
func seedOrg(t *testing.T) uuid.UUID {
	t.Helper()
	res, err := orgSvc.Create(context.Background(), orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:         "Тестовая организация",
			LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
		},
	})
	require.NoError(t, err)
	return res.ID
}

// seedClinic creates a clinic under orgID.
func seedClinic(t *testing.T, orgID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := clinSvc.Create(context.Background(), orgsvc.CreateClinicCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{
			OrganizationID:  orgID.String(),
			Name:            "Тестовая клиника",
			PhysicalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Тверская, д. 10"},
		},
	})
	require.NoError(t, err)
	return res.ID
}

// seedDept creates a department under clinicID.
func seedDept(t *testing.T, clinicID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := deptSvc.Create(context.Background(), orgsvc.CreateDepartmentCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{
			ClinicID: clinicID.String(),
			Name:     "Тестовое отделение",
		},
	})
	require.NoError(t, err)
	return res.ID
}

// seedCategory creates an active incident category for orgID.
func seedCategory(t *testing.T, orgID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := categorySvc.Create(context.Background(), classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: orgID.String(),
			Name:           "Тестовая категория",
		},
	})
	require.NoError(t, err)
	return res.ID
}

// seedType creates an active incident type under categoryID; patient-allowed
// if allowPatients is true.
func seedType(t *testing.T, categoryID uuid.UUID, allowPatients bool) uuid.UUID {
	t.Helper()
	res, err := typeSvc.Create(context.Background(), classifiersvc.CreateIncidentTypeCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{
			CategoryID: categoryID.String(),
			Name:       "Тестовый тип",
		},
	})
	require.NoError(t, err)
	if allowPatients {
		_, err = typeSvc.AllowForPatients(context.Background(), classifiersvc.AllowIncidentTypeForPatientsCommand{
			Caller:  sysadminCaller,
			Payload: classifiersvc.AllowIncidentTypeForPatientsPayload{TypeID: res.ID.String()},
		})
		require.NoError(t, err)
	}
	return res.ID
}

// seedEmployee directly inserts a domain.employees row (bypassing the
// Zitadel verifier used by EmployeeService.Hire) and its projection row.
// Returns the new employee UUID.
func seedEmployee(t *testing.T, zitadelID string, orgID, deptID uuid.UUID) uuid.UUID {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	empID, err := uuid.NewV7()
	require.NoError(t, err)
	if _, err := raw.Exec(
		`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id)
		 VALUES ($1, $2, $3, $4)`,
		empID, zitadelID, orgID, deptID,
	); err != nil {
		t.Fatalf("insert domain.employees: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO projections.employees (id, zitadel_user_id, organization_id, department_id, hired_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, now(), now(), now())`,
		empID, zitadelID, orgID, deptID,
	); err != nil {
		t.Fatalf("insert projections.employees: %v", err)
	}
	return empID
}

// seedOrgAdmin directly inserts the org_admins role row + projection.
func seedOrgAdmin(t *testing.T, empID, orgID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.org_admins (organization_id, employee_id) VALUES ($1, $2)`,
		orgID, empID,
	); err != nil {
		t.Fatalf("insert domain.org_admins: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO projections.org_admins (organization_id, employee_id, created_at)
		 VALUES ($1, $2, now())`,
		orgID, empID,
	); err != nil {
		t.Fatalf("insert projections.org_admins: %v", err)
	}
}

// seedOrgHead directly inserts the org_heads role row + projection.
func seedOrgHead(t *testing.T, empID, orgID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.org_heads (organization_id, employee_id) VALUES ($1, $2)`,
		orgID, empID,
	); err != nil {
		t.Fatalf("insert domain.org_heads: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO projections.org_heads (organization_id, employee_id, created_at)
		 VALUES ($1, $2, now())`,
		orgID, empID,
	); err != nil {
		t.Fatalf("insert projections.org_heads: %v", err)
	}
}

// seedClinicHead directly inserts the clinic_heads role row + projection.
func seedClinicHead(t *testing.T, empID, clinicID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.clinic_heads (clinic_id, employee_id) VALUES ($1, $2)`,
		clinicID, empID,
	); err != nil {
		t.Fatalf("insert domain.clinic_heads: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO projections.clinic_heads (clinic_id, employee_id, created_at)
		 VALUES ($1, $2, now())`,
		clinicID, empID,
	); err != nil {
		t.Fatalf("insert projections.clinic_heads: %v", err)
	}
}

// seedDeptResponsible directly inserts the department_responsibles role + projection.
func seedDeptResponsible(t *testing.T, empID, deptID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.department_responsibles (department_id, employee_id) VALUES ($1, $2)`,
		deptID, empID,
	); err != nil {
		t.Fatalf("insert domain.department_responsibles: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO projections.department_responsibles (department_id, employee_id, created_at)
		 VALUES ($1, $2, now())`,
		deptID, empID,
	); err != nil {
		t.Fatalf("insert projections.department_responsibles: %v", err)
	}
}

// seedOrgDispatcher directly inserts the org_dispatchers role row + projection.
func seedOrgDispatcher(t *testing.T, empID, orgID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.org_dispatchers (organization_id, employee_id) VALUES ($1, $2)`,
		orgID, empID,
	); err != nil {
		t.Fatalf("insert domain.org_dispatchers: %v", err)
	}
	if _, err := raw.Exec(
		`INSERT INTO projections.org_dispatchers (organization_id, employee_id, created_at)
		 VALUES ($1, $2, now())`,
		orgID, empID,
	); err != nil {
		t.Fatalf("insert projections.org_dispatchers: %v", err)
	}
}

// createIncident is a convenience wrapper around incidentSvc.Create with
// sensible defaults.
func createIncident(
	t *testing.T,
	ctx context.Context,
	caller authz.Caller,
	deptID, categoryID, typeID uuid.UUID,
) uuid.UUID {
	t.Helper()
	desc := "Тестовое описание инцидента достаточной длины"
	res, err := incidentSvc.Create(ctx, &incidentsvc.CreateIncidentCommand{
		Caller: caller,
		Payload: incidentsvc.CreateIncidentPayload{
			DepartmentID: deptID.String(),
			CategoryID:   categoryID.String(),
			TypeID:       typeID.String(),
			Description:  &desc,
			OccurredAt:   time.Now().Add(-1 * time.Hour),
		},
	})
	require.NoError(t, err)
	return res.ID
}

// countStatusHistory returns the number of status history rows for an incident.
func countStatusHistory(t *testing.T, incidentID uuid.UUID) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM projections.incident_status_history WHERE incident_id = ?`,
		incidentID,
	).Scan(&n).Error)
	return n
}

// countPriorityHistory returns the number of priority history rows for an incident.
func countPriorityHistory(t *testing.T, incidentID uuid.UUID) int {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM projections.incident_priority_history WHERE incident_id = ?`,
		incidentID,
	).Scan(&n).Error)
	return n
}
