//go:build integration

package announcement_query_integration_test

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

	"github.com/medincident/medincident-backend/internal/model"
	announcementread "github.com/medincident/medincident-backend/internal/service/query/announcement"
)

const sysadminZitadelID = "sysadmin"

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()
	reader     *announcementread.Reader
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

	reader = announcementread.NewReader(testDB, &testLogger)
	code := m.Run()
	if err := container.Terminate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "terminate postgres: %v\n", err)
	}
	os.Exit(code)
}

func resetProjections(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	for _, q := range []string{
		`TRUNCATE TABLE domain.system_admins CASCADE`,
		`TRUNCATE TABLE domain.department_responsibles CASCADE`,
		`TRUNCATE TABLE domain.clinic_heads CASCADE`,
		`TRUNCATE TABLE domain.org_admins CASCADE`,
		`TRUNCATE TABLE domain.employees CASCADE`,
		`TRUNCATE TABLE domain.departments CASCADE`,
		`TRUNCATE TABLE domain.clinics CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
		`TRUNCATE TABLE domain.announcements CASCADE`,
		`TRUNCATE TABLE projections.announcement_views CASCADE`,
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

func insertClinic(t *testing.T, orgID uuid.UUID) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	require.NoError(t, err)
	require.NoError(t, testDB.Create(&model.Clinic{
		ID:              id,
		OrganizationID:  orgID,
		Name:            "Clinic",
		PhysicalAddress: model.Address{Text: "г. Тест, ул. Клиническая"},
	}).Error)
	return id
}

func insertDept(t *testing.T, clinicID uuid.UUID) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	require.NoError(t, err)
	require.NoError(t, testDB.Create(&model.Department{
		ID:       id,
		ClinicID: clinicID,
		Name:     "Dept",
	}).Error)
	return id
}

func insertAnnouncement(t *testing.T, orgID uuid.UUID, clinicID, deptID uuid.NullUUID, title string) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	require.NoError(t, err)
	now := time.Now()
	a := model.Announcement{
		ID:             id,
		OrganizationID: orgID,
		ClinicID:       clinicID,
		DepartmentID:   deptID,
		AuthorID:       sysadminZitadelID,
		Title:          title,
		Content:        "Тестовое содержание объявления",
		Priority:       model.AnnouncementPriorityNormal,
		IsArchived:     false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.Create(&a).Error)
	require.NoError(t, testDB.Exec(
		`INSERT INTO projections.announcement_views (announcement_id, view_count) VALUES (?, 0)`, id,
	).Error)
	return id
}

func insertEmployee(t *testing.T, zitadelID string, orgID, deptID uuid.UUID) {
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
}
