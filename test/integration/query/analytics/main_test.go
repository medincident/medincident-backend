//go:build integration

package analytics_query_integration_test

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
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	requestclassifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
	analyticsread "github.com/medincident/medincident-backend/internal/service/query/analytics"
)

const (
	sysadminZitadelID = "sysadmin"
	userAZitadelID    = "user-a"
	userBZitadelID    = "user-b"
)

var (
	sysadminCaller = authz.Caller{ZitadelUserID: sysadminZitadelID}
	testLogger     = zerolog.Nop()
	testDB         *gorm.DB
	authzSvc       *authz.Authz
	analyticsRdr   *analyticsread.Reader

	orgSvc         *orgsvc.OrganizationService
	clinicSvc      *orgsvc.ClinicService
	deptSvc        *orgsvc.DepartmentService
	categorySvc    *classifiersvc.IncidentCategoryService
	typeSvc        *classifiersvc.IncidentTypeService
	incidentSvc    *incidentsvc.IncidentService
	requestTypeSvc *requestclassifiersvc.RequestTypeService
	requestSvc     *requestsvc.ServiceRequestService
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
	defer func() { _ = container.Terminate(ctx) }()

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get dsn: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "tool", "dbmate",
		"--url", dsn, "--migrations-dir", "../../../../db/migrations",
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
	analyticsRdr = analyticsread.NewReader(testDB, authzSvc, &testLogger)

	orgSvc = orgsvc.NewOrganizationService(testDB, authzSvc, &testLogger)
	clinicSvc = orgsvc.NewClinicService(testDB, authzSvc, &testLogger)
	deptSvc = orgsvc.NewDepartmentService(testDB, authzSvc, &testLogger)
	categorySvc = classifiersvc.NewIncidentCategoryService(testDB, authzSvc, &testLogger)
	typeSvc = classifiersvc.NewIncidentTypeService(testDB, authzSvc, &testLogger)
	incidentSvc = incidentsvc.NewIncidentService(testDB, authzSvc, &testLogger)
	requestTypeSvc = requestclassifiersvc.NewRequestTypeService(testDB, authzSvc, &testLogger)
	requestSvc = requestsvc.NewServiceRequestService(testDB, authzSvc, &testLogger)

	os.Exit(m.Run())
}

func resetDB(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	for _, q := range []string{
		`TRUNCATE TABLE domain.system_admins CASCADE`,
		`TRUNCATE TABLE domain.service_request_executors CASCADE`,
		`TRUNCATE TABLE domain.service_requests CASCADE`,
		`TRUNCATE TABLE domain.patient_incident_buffer CASCADE`,
		`TRUNCATE TABLE domain.incidents CASCADE`,
		`TRUNCATE TABLE domain.employees CASCADE`,
		`TRUNCATE TABLE domain.incident_types CASCADE`,
		`TRUNCATE TABLE domain.incident_categories CASCADE`,
		`TRUNCATE TABLE domain.request_types CASCADE`,
		`TRUNCATE TABLE domain.departments CASCADE`,
		`TRUNCATE TABLE domain.clinics CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
		`TRUNCATE TABLE projections.incidents,
		               projections.incident_status_history,
		               projections.incident_priority_history,
		               projections.patient_incident_buffer,
		               projections.incident_categories,
		               projections.incident_types,
		               projections.service_requests,
		               projections.service_request_status_history,
		               projections.service_request_executor_history,
		               projections.request_types,
		               projections.users,
		               projections.employees,
		               projections.organization_counters,
		               projections.clinic_counters,
		               projections.department_counters,
		               projections.departments,
		               projections.clinics,
		               projections.organizations,
		               projections.clinic_heads,
		               projections.department_responsibles,
		               projections.org_admins,
		               projections.org_dispatchers,
		               projections.org_heads,
		               projections.system_admins CASCADE`,
	} {
		if _, err := raw.Exec(q); err != nil {
			t.Fatalf("truncate: %v", err)
		}
	}
	if _, err := raw.Exec(
		`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1)`, sysadminZitadelID,
	); err != nil {
		t.Fatalf("seed sysadmin: %v", err)
	}
	seedUser(t, sysadminZitadelID, "System Admin")
}

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
		t.Fatalf("seed user %s: %v", zitadelID, err)
	}
}

func seedOrg(t *testing.T) uuid.UUID {
	t.Helper()
	res, err := orgSvc.Create(context.Background(), orgsvc.CreateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{Name: "Тестовая организация", LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"}},
	})
	require.NoError(t, err)
	return res.ID
}

func seedClinic(t *testing.T, orgID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := clinicSvc.Create(context.Background(), orgsvc.CreateClinicCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{
			OrganizationID:  orgID.String(),
			Name:            "Тестовая клиника",
			PhysicalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Тверская, д. 1"},
		},
	})
	require.NoError(t, err)
	return res.ID
}

func seedDept(t *testing.T, clinicID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := deptSvc.Create(context.Background(), orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clinicID.String(), Name: "Тестовое отделение"},
	})
	require.NoError(t, err)
	return res.ID
}

func seedCategory(t *testing.T, orgID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := categorySvc.Create(context.Background(), classifiersvc.CreateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{OrganizationID: orgID.String(), Name: "Категория"},
	})
	require.NoError(t, err)
	return res.ID
}

func seedIncidentType(t *testing.T, categoryID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := typeSvc.Create(context.Background(), classifiersvc.CreateIncidentTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{CategoryID: categoryID.String(), Name: "Тип инцидента"},
	})
	require.NoError(t, err)
	return res.ID
}

func seedCategoryNamed(t *testing.T, orgID uuid.UUID, name string) uuid.UUID {
	t.Helper()
	res, err := categorySvc.Create(context.Background(), classifiersvc.CreateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{OrganizationID: orgID.String(), Name: name},
	})
	require.NoError(t, err)
	return res.ID
}

func seedIncidentTypeNamed(t *testing.T, categoryID uuid.UUID, name string) uuid.UUID {
	t.Helper()
	res, err := typeSvc.Create(context.Background(), classifiersvc.CreateIncidentTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{CategoryID: categoryID.String(), Name: name},
	})
	require.NoError(t, err)
	return res.ID
}

func seedRequestType(t *testing.T, orgID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := requestTypeSvc.Create(context.Background(), requestclassifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: requestclassifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Тип заявки"},
	})
	require.NoError(t, err)
	return res.ID
}

func seedEmployee(t *testing.T, zitadelID string, orgID, deptID uuid.UUID) uuid.UUID {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	empID, err := uuid.NewV7()
	require.NoError(t, err)
	if _, err := raw.Exec(
		`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id) VALUES ($1, $2, $3, $4)`,
		empID, zitadelID, orgID, deptID,
	); err != nil {
		t.Fatalf("insert employee: %v", err)
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

func seedOrgAdmin(t *testing.T, empID, orgID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO domain.org_admins (organization_id, employee_id) VALUES ($1, $2)`, orgID, empID); err != nil {
		t.Fatalf("insert domain.org_admins: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO projections.org_admins (organization_id, employee_id, created_at) VALUES ($1, $2, now())`, orgID, empID); err != nil {
		t.Fatalf("insert projections.org_admins: %v", err)
	}
}

func seedClinicHead(t *testing.T, empID, clinicID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO domain.clinic_heads (clinic_id, employee_id) VALUES ($1, $2)`, clinicID, empID); err != nil {
		t.Fatalf("insert domain.clinic_heads: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO projections.clinic_heads (clinic_id, employee_id, created_at) VALUES ($1, $2, now())`, clinicID, empID); err != nil {
		t.Fatalf("insert projections.clinic_heads: %v", err)
	}
}

func seedDeptResponsible(t *testing.T, empID, deptID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO domain.department_responsibles (department_id, employee_id) VALUES ($1, $2)`, deptID, empID); err != nil {
		t.Fatalf("insert domain.department_responsibles: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO projections.department_responsibles (department_id, employee_id, created_at) VALUES ($1, $2, now())`, deptID, empID); err != nil {
		t.Fatalf("insert projections.department_responsibles: %v", err)
	}
}

func createIncident(t *testing.T, caller authz.Caller, deptID, categoryID, typeID uuid.UUID) uuid.UUID {
	t.Helper()
	desc := "Тестовое описание инцидента достаточной длины для валидации"
	res, err := incidentSvc.Create(context.Background(), &incidentsvc.CreateIncidentCommand{
		Caller: caller,
		Payload: incidentsvc.CreateIncidentPayload{
			DepartmentID: deptID.String(),
			CategoryID:   categoryID.String(),
			TypeID:       typeID.String(),
			Description:  &desc,
			OccurredAt:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339Nano),
		},
	})
	require.NoError(t, err)
	return res.ID
}

func uuidNew() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id.String()
}

// updateIncidentStatus transitions incidentID via incidentSvc.UpdateStatus.
func updateIncidentStatus(t *testing.T, caller authz.Caller, incidentID uuid.UUID, newStatus string) {
	t.Helper()
	require.NoError(t, incidentSvc.UpdateStatus(context.Background(), incidentsvc.UpdateIncidentStatusCommand{
		Caller: caller,
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: incidentID.String(),
			NewStatus:  newStatus,
		},
	}))
}

// doneIncident transitions incidentID to done (pending→in_progress→done).
func doneIncident(t *testing.T, caller authz.Caller, incidentID uuid.UUID) {
	t.Helper()
	updateIncidentStatus(t, caller, incidentID, "in_progress")
	updateIncidentStatus(t, caller, incidentID, "done")
}

// createRequest creates a service request and returns its ID.
// executorEmpID is the employee who will be the executor.
func createRequest(t *testing.T, caller authz.Caller, deptID, typeID, executorEmpID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := requestSvc.Create(context.Background(), &requestsvc.CreateServiceRequestCommand{
		Caller: caller,
		Payload: requestsvc.CreateServiceRequestPayload{
			DepartmentID:        deptID.String(),
			TypeID:              typeID.String(),
			Description:         "Тестовое описание заявки достаточной длины для прохождения валидации",
			ExecutorEmployeeIDs: []string{executorEmpID.String()},
		},
	})
	require.NoError(t, err)
	return res.ID
}

// completeRequest transitions requestID to completed (created→in_work→pending_review→completed).
// executorCaller is the executor; managerCaller is a privileged role (OrgAdmin etc) for the final step.
func completeRequest(t *testing.T, executorCaller, managerCaller authz.Caller, requestID uuid.UUID) {
	t.Helper()
	require.NoError(t, requestSvc.UpdateStatus(context.Background(), requestsvc.UpdateServiceRequestStatusCommand{
		Caller:  executorCaller,
		Payload: requestsvc.UpdateServiceRequestStatusPayload{ServiceRequestID: requestID.String(), NewStatus: "in_work"},
	}))
	require.NoError(t, requestSvc.UpdateStatus(context.Background(), requestsvc.UpdateServiceRequestStatusCommand{
		Caller:  executorCaller,
		Payload: requestsvc.UpdateServiceRequestStatusPayload{ServiceRequestID: requestID.String(), NewStatus: "pending_review"},
	}))
	require.NoError(t, requestSvc.UpdateStatus(context.Background(), requestsvc.UpdateServiceRequestStatusCommand{
		Caller:  managerCaller,
		Payload: requestsvc.UpdateServiceRequestStatusPayload{ServiceRequestID: requestID.String(), NewStatus: "completed"},
	}))
}

// seedOrgHead grants orgHead role to an employee in an org.
func seedOrgHead(t *testing.T, empID, orgID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO domain.org_heads (organization_id, employee_id) VALUES ($1, $2)`, orgID, empID); err != nil {
		t.Fatalf("insert domain.org_heads: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO projections.org_heads (organization_id, employee_id, created_at) VALUES ($1, $2, now())`, orgID, empID); err != nil {
		t.Fatalf("insert projections.org_heads: %v", err)
	}
}

// seedOrgDispatcher grants orgDispatcher role to an employee in an org.
func seedOrgDispatcher(t *testing.T, empID, orgID uuid.UUID) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO domain.org_dispatchers (organization_id, employee_id) VALUES ($1, $2)`, orgID, empID); err != nil {
		t.Fatalf("insert domain.org_dispatchers: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO projections.org_dispatchers (organization_id, employee_id, created_at) VALUES ($1, $2, now())`, orgID, empID); err != nil {
		t.Fatalf("insert projections.org_dispatchers: %v", err)
	}
}
