//go:build integration

package servicerequest_integration_test

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
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
)

const (
	sysadminZitadelID = "sysadmin"
	executorZitadelID = "executor"
)

var (
	sysadminCaller = authz.Caller{ZitadelUserID: sysadminZitadelID}
	testDB         *gorm.DB
	testLogger     = zerolog.Nop()
	authzSvc       *authz.Authz
	requestSvc     *requestsvc.ServiceRequestService
	orgSvc         *orgsvc.OrganizationService
	clinicSvc      *orgsvc.ClinicService
	deptSvc        *orgsvc.DepartmentService
	requestTypeSvc *classifiersvc.RequestTypeService
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
	requestSvc = requestsvc.NewServiceRequestService(testDB, authzSvc, &testLogger)
	orgSvc = orgsvc.NewOrganizationService(testDB, authzSvc, &testLogger)
	clinicSvc = orgsvc.NewClinicService(testDB, authzSvc, &testLogger)
	deptSvc = orgsvc.NewDepartmentService(testDB, authzSvc, &testLogger)
	requestTypeSvc = classifiersvc.NewRequestTypeService(testDB, authzSvc, &testLogger)
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
		`TRUNCATE TABLE domain.service_request_executors CASCADE`,
		`TRUNCATE TABLE domain.service_requests CASCADE`,
		`TRUNCATE TABLE domain.employees CASCADE`,
		`TRUNCATE TABLE domain.request_types CASCADE`,
		`TRUNCATE TABLE domain.departments CASCADE`,
		`TRUNCATE TABLE domain.clinics CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
		`TRUNCATE TABLE projections.service_requests,
		               projections.service_request_status_history,
		               projections.service_request_executor_history,
		               projections.request_types,
		               projections.users,
		               projections.employees,
		               projections.departments,
		               projections.clinics,
		               projections.organizations CASCADE`,
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
	// ServiceRequestService.resolveActorDisplayName queries projections.users.
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
		t.Fatalf("seed projections.users %s: %v", zitadelID, err)
	}
}

func seedOrg(t *testing.T) uuid.UUID {
	t.Helper()
	res, err := orgSvc.Create(context.Background(), orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:         "Организация",
			LegalAddress: orgsvc.AddressInput{Text: "г. Тест, ул. Тестовая, д. 1"},
		},
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
			Name:            "Клиника",
			PhysicalAddress: orgsvc.AddressInput{Text: "г. Тест, ул. Клиническая, д. 1"},
		},
	})
	require.NoError(t, err)
	return res.ID
}

func seedDept(t *testing.T, clinicID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := deptSvc.Create(context.Background(), orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clinicID.String(), Name: "Отделение"},
	})
	require.NoError(t, err)
	return res.ID
}

func seedRequestType(t *testing.T, orgID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := requestTypeSvc.Create(context.Background(), classifiersvc.CreateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateRequestTypePayload{OrganizationID: orgID.String(), Name: "Тип заявки"},
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
		`INSERT INTO domain.employees (id, zitadel_user_id, organization_id, department_id)
		 VALUES ($1, $2, $3, $4)`,
		empID, zitadelID, orgID, deptID,
	); err != nil {
		t.Fatalf("insert domain.employees: %v", err)
	}
	return empID
}

func loadRequest(t *testing.T, id uuid.UUID) model.ServiceRequest {
	t.Helper()
	var sr model.ServiceRequest
	require.NoError(t, testDB.First(&sr, "id = ?", id).Error)
	return sr
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
