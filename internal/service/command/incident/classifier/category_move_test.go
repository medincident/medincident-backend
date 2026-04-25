//go:build integration

package classifier_test

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
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
)

const sysadminZitadelID = "sysadmin"

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()

	authzSvc    *authz.Authz
	categorySvc *classifiersvc.IncidentCategoryService

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
	categorySvc = classifiersvc.NewIncidentCategoryService(testDB, authzSvc, &testLogger)

	os.Exit(m.Run())
}

func runMigrations(dsn string) error {
	cmd := exec.Command("go", "tool", "dbmate",
		"--url", dsn,
		"--migrations-dir", "../../../../../db/migrations",
		"--no-dump-schema",
		"up",
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// resetDB truncates all classifier-related tables and re-seeds the
// system-admin used as the command caller by every test.
func resetDB(t *testing.T) {
	t.Helper()
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	truncate := []string{
		`TRUNCATE TABLE domain.system_admins CASCADE`,
		`TRUNCATE TABLE domain.incident_types CASCADE`,
		`TRUNCATE TABLE domain.incident_categories CASCADE`,
		`TRUNCATE TABLE domain.departments CASCADE`,
		`TRUNCATE TABLE domain.clinics CASCADE`,
		`TRUNCATE TABLE domain.organizations CASCADE`,
		`TRUNCATE TABLE projections.incident_types,
		                 projections.incident_categories,
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
	if _, err := raw.Exec(
		`INSERT INTO domain.system_admins (zitadel_user_id) VALUES ($1)`,
		sysadminZitadelID,
	); err != nil {
		t.Fatalf("seed sysadmin: %v", err)
	}
}

// insertOrg creates a minimal organization row and returns its ID.
func insertOrg(t *testing.T, name string) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	row := model.Organization{
		ID:   id,
		Name: name,
		LegalAddress: model.Address{
			Text: "г. Тест, ул. Тестовая",
		},
	}
	require.NoError(t, testDB.Create(&row).Error)
	return id
}

// createCategory is a thin wrapper around categorySvc.Create that fails
// the test on error and returns the new category ID.
func createCategory(t *testing.T, orgID uuid.UUID, name string, parentID *uuid.UUID) uuid.UUID {
	t.Helper()
	var parentStr *string
	if parentID != nil {
		s := parentID.String()
		parentStr = &s
	}
	res, err := categorySvc.Create(context.Background(), classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   orgID.String(),
			ParentCategoryID: parentStr,
			Name:             name,
		},
	})
	require.NoError(t, err)
	return res.ID
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

// uuidStrPtr returns a pointer to the string representation of id.
func uuidStrPtr(id uuid.UUID) *string {
	s := id.String()
	return &s
}

// TestCategoryMove_ValidMove verifies that moving a category to a new
// parent within the same organisation succeeds.
func TestCategoryMove_ValidMove(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t, "Org")

	rootA := createCategory(t, orgID, "Root A", nil)
	rootB := createCategory(t, orgID, "Root B", nil)
	child := createCategory(t, orgID, "Child", &rootA)

	_, err := categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID:          child.String(),
			NewParentCategoryID: uuidStrPtr(rootB),
		},
	})
	require.NoError(t, err)

	var parentID uuid.NullUUID
	require.NoError(t, testDB.Raw(
		`SELECT parent_category_id FROM domain.incident_categories WHERE id = ?`, child,
	).Row().Scan(&parentID))
	assert.True(t, parentID.Valid)
	assert.Equal(t, rootB, parentID.UUID)
}

// TestCategoryMove_ExceedsMaxDepth verifies that moving a subtree whose
// combined depth would exceed incidentClassifierMaxDepth (5) is rejected
// with ErrCodeIncidentCategoryMoveWouldExceedDepth.
//
// Scenario: A(1)->B(2)->C(3) moved under X(1)->Y(2)->Z(3) gives depth 6.
func TestCategoryMove_ExceedsMaxDepth(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t, "Org")

	// Subtree to be moved: A -> B -> C (internal depth 3).
	a := createCategory(t, orgID, "AA", nil)
	b := createCategory(t, orgID, "BB", &a)
	_ = createCategory(t, orgID, "CC", &b)

	// Target parent chain: X -> Y -> Z (depth of Z = 3).
	x := createCategory(t, orgID, "XX", nil)
	y := createCategory(t, orgID, "YY", &x)
	z := createCategory(t, orgID, "ZZ", &y)

	// Moving A under Z: parentDepth(Z)=3 + subtreeDepth(A)=3 = 6 > 5.
	_, err := categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID:          a.String(),
			NewParentCategoryID: uuidStrPtr(z),
		},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryMoveWouldExceedDepth, codeOf(t, err))
}

// TestCategoryMove_CreatesCycle verifies that moving a category under
// one of its own descendants is rejected with
// ErrCodeIncidentCategoryMoveWouldCreateCycle.
func TestCategoryMove_CreatesCycle(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t, "Org")

	// Build chain: A -> B -> C.
	a := createCategory(t, orgID, "AA", nil)
	b := createCategory(t, orgID, "BB", &a)
	c := createCategory(t, orgID, "CC", &b)

	// Moving A under C would form A -> B -> C -> A.
	_, err := categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID:          a.String(),
			NewParentCategoryID: uuidStrPtr(c),
		},
	})
	require.Error(t, err)
	assert.Equal(t, classifiersvc.ErrCodeIncidentCategoryMoveWouldCreateCycle, codeOf(t, err))
}
