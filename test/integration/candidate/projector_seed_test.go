//go:build integration

package candidate_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/query/membership"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
)

// resetForHireProjector truncates projections tables used by the projector
// seeding tests (in addition to the domain/auth tables from resetAll).
func resetForHireProjector(t *testing.T) {
	t.Helper()
	resetAll(t)
	raw, err := testDB.DB()
	if err != nil {
		t.Fatalf("get raw db: %v", err)
	}
	if _, err := raw.Exec(`TRUNCATE TABLE
		projections.organization_counters,
		projections.clinic_counters,
		projections.department_counters,
		projections.employees,
		projections.departments,
		projections.clinics,
		projections.organizations
	CASCADE`); err != nil {
		t.Fatalf("truncate projections: %v", err)
	}
}

// seedOrgProjection seeds a minimal organization projection row.
func seedOrgProjection(t *testing.T, ctx context.Context, orgID uuid.UUID, now time.Time) {
	t.Helper()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrganizationCreated(tx, orgID.String(), now, &orgv1.OrganizationCreated{
			Name:         "Org " + orgID.String()[:8],
			LegalAddress: &orgv1.Address{Text: "addr"},
			CreatedAt:    timestamppb.New(now),
		})
	}))
}

// seedClinicProjection seeds a minimal clinic projection row under the org.
func seedClinicProjection(t *testing.T, ctx context.Context, clinicID, orgID uuid.UUID, now time.Time) {
	t.Helper()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicCreated(tx, clinicID.String(), now, &clinicv1.ClinicCreated{
			OrganizationId:  orgID.String(),
			Name:            "Clinic " + clinicID.String()[:8],
			PhysicalAddress: &clinicv1.Address{Text: "addr"},
			CreatedAt:       timestamppb.New(now),
		})
	}))
}

// seedDeptProjection seeds a minimal department projection row under the clinic.
func seedDeptProjection(t *testing.T, ctx context.Context, deptID, clinicID uuid.UUID, now time.Time) {
	t.Helper()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DepartmentCreated(tx, deptID.String(), now, &deptv1.DepartmentCreated{
			ClinicId:  clinicID.String(),
			Name:      "Dept " + deptID.String()[:8],
			CreatedAt: timestamppb.New(now),
		})
	}))
}

// TestForHire_ProjectorSeededUsers verifies the full hire-to-ForHire
// pipeline without the async identity consumer: the EmployeeHired
// projector must seed projections.users so the user shows up as a
// hire candidate for a different organization.
func TestForHire_ProjectorSeededUsers(t *testing.T) {
	resetForHireProjector(t)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	// Two organisations. User will be employed in orgB, queried as
	// candidate for orgA.
	orgA := uuid.Must(uuid.NewV7())
	orgB := uuid.Must(uuid.NewV7())
	clinicA := uuid.Must(uuid.NewV7())
	clinicB := uuid.Must(uuid.NewV7())
	deptA := uuid.Must(uuid.NewV7())
	deptB := uuid.Must(uuid.NewV7())

	// Seed projections for both orgs.
	for _, id := range []uuid.UUID{orgA, orgB} {
		seedOrgProjection(t, ctx, id, now)
	}
	seedClinicProjection(t, ctx, clinicA, orgA, now)
	seedClinicProjection(t, ctx, clinicB, orgB, now)
	seedDeptProjection(t, ctx, deptA, clinicA, now)
	seedDeptProjection(t, ctx, deptB, clinicB, now)

	// Seed domain.* tables for authz checks.
	domainOrg(t, orgA)
	domainOrg(t, orgB)
	domainClinic(t, clinicA, orgA)
	domainClinic(t, clinicB, orgB)
	domainDepartment(t, deptA, clinicA)
	domainDepartment(t, deptB, clinicB)

	// System admin caller — simplest authz path.
	callerZID := "sysadmin-projector-test"
	domainSystemAdmin(t, callerZID)
	callerA := authz.Caller{ZitadelUserID: callerZID}

	// Org admin caller for orgA.
	callerAdminZID := "org-admin-projector-test"
	callerAdmin := orgAdminCaller(t, callerAdminZID, orgA, clinicA, deptA)

	// Hire userB into orgB via the projector (simulates the query-server
	// processing the EmployeeHired event after the command-server writes it).
	userBZID := "user-b-projector"
	empBID := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.EmployeeHired(tx, empBID.String(), now, &empv1.EmployeeHired{
			ZitadelUserId:     userBZID,
			OrganizationId:    orgB.String(),
			DepartmentId:      deptB.String(),
			HiredAt:           timestamppb.New(now),
			UserName:          "bjohnson",
			FirstName:         "Bob",
			LastName:          "Johnson",
			DisplayName:       "Bob Johnson",
			Email:             "bob@example.com",
			PreferredLanguage: "en",
		})
	}))

	log := zerolog.Nop()
	reader := membership.NewCandidateReader(testDB, authzSvc, &log)

	// System admin sees userB as candidate for orgA.
	results, _, err := reader.ForHire(ctx, callerA, orgA, "", "", 50)
	require.NoError(t, err)
	assert.True(t, containsZitadelUserID(results, userBZID),
		"userB (employed in orgB) must appear as candidate for orgA")

	// Org admin of orgA also sees userB.
	results2, _, err := reader.ForHire(ctx, callerAdmin, orgA, "", "", 50)
	require.NoError(t, err)
	assert.True(t, containsZitadelUserID(results2, userBZID),
		"userB must appear for orgA admin too")
}
