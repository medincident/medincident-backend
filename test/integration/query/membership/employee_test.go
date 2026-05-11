//go:build integration

package membership_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	memberread "github.com/medincident/medincident-backend/internal/service/query/membership"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	vacv1 "github.com/medincident/medincident-backend/pkg/event/vacation/v1"
)

// seedOrgClinicDept inserts an organization + clinic + department
// through the projector so the employee projection has parents.
func seedOrgClinicDept(t *testing.T, ctx context.Context, now time.Time) (uuid.UUID, uuid.UUID, uuid.UUID) {
	t.Helper()
	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	org := &model.Organization{ID: orgID, Name: "O", LegalAddress: model.Address{Text: "S"}, CreatedAt: now, UpdatedAt: now}
	clinic := &model.Clinic{ID: clinicID, OrganizationID: orgID, Name: "C", PhysicalAddress: model.Address{Text: "S"}, CreatedAt: now, UpdatedAt: now}
	dept := &model.Department{ID: deptID, ClinicID: clinicID, Name: "D", CreatedAt: now, UpdatedAt: now}

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qprojector.OrganizationCreated(tx, org.ID.String(), org.CreatedAt, &orgv1.OrganizationCreated{
			Name:         org.Name,
			LegalAddress: &orgv1.Address{Text: org.LegalAddress.Text},
			CreatedAt:    timestamppb.New(org.CreatedAt),
		}); err != nil {
			return err
		}
		if err := qprojector.ClinicCreated(tx, clinic.ID.String(), clinic.CreatedAt, &clinicv1.ClinicCreated{
			OrganizationId:  clinic.OrganizationID.String(),
			Name:            clinic.Name,
			PhysicalAddress: &clinicv1.Address{Text: clinic.PhysicalAddress.Text},
			CreatedAt:       timestamppb.New(clinic.CreatedAt),
		}); err != nil {
			return err
		}
		return qprojector.DepartmentCreated(tx, dept.ID.String(), dept.CreatedAt, &deptv1.DepartmentCreated{
			ClinicId:  dept.ClinicID.String(),
			Name:      dept.Name,
			CreatedAt: timestamppb.New(dept.CreatedAt),
		})
	}))
	return orgID, clinicID, deptID
}

// TestEmployeeReader_Get_AndListByDepartment seeds an employee through
// the projector and reads it back.
func TestEmployeeReader_Get_AndListByDepartment(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID, _, deptID := seedOrgClinicDept(t, ctx, now)

	empID := uuid.Must(uuid.NewV7())
	emp := &model.Employee{
		ID:             empID,
		ZitadelUserID:  "zit-1",
		OrganizationID: orgID,
		DepartmentID:   deptID,
		Position:       null.StringFrom("Nurse"),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.EmployeeHired(tx, emp.ID.String(), emp.CreatedAt, &empv1.EmployeeHired{
			ZitadelUserId:  emp.ZitadelUserID,
			OrganizationId: emp.OrganizationID.String(),
			DepartmentId:   emp.DepartmentID.String(),
			Position:       emp.Position.String,
			HiredAt:        timestamppb.New(emp.CreatedAt),
		})
	}))

	reader := memberread.NewEmployeeReader(testDB, authzSvc, &logger)
	got, err := reader.Get(ctx, sysadminCaller, empID)
	require.NoError(t, err)
	require.Equal(t, empID, got.EmployeeID)
	require.Equal(t, "zit-1", got.ZitadelUserID)
	require.Equal(t, deptID, got.DepartmentID)
	require.NotNil(t, got.Position)
	require.Equal(t, "Nurse", *got.Position)

	list, err := reader.ListByDepartment(ctx, sysadminCaller, deptID, memberread.ListQuery{}, memberread.EmployeeFilter{})
	require.NoError(t, err)
	require.Len(t, list.Items, 1)
	require.Equal(t, empID, list.Items[0].EmployeeID)
}

// TestEmployeeReader_Get_NotFound surfaces the expected code.
func TestEmployeeReader_Get_NotFound(t *testing.T) {
	resetProjections(t)
	logger := zerolog.Nop()
	reader := memberread.NewEmployeeReader(testDB, authzSvc, &logger)
	_, err := reader.Get(context.Background(), sysadminCaller, uuid.Must(uuid.NewV7()))
	require.Error(t, err)
}

// TestEmployeeReader_EmployeeFilters exercises the three filter paths
// end-to-end: one terminated employee, one active employee on an
// active vacation with position "Nurse".
func TestEmployeeReader_EmployeeFilters(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)
	orgID, _, deptID := seedOrgClinicDept(t, ctx, now)

	activeID := uuid.Must(uuid.NewV7())
	termID := uuid.Must(uuid.NewV7())
	active := &model.Employee{ID: activeID, ZitadelUserID: "zit-a", OrganizationID: orgID, DepartmentID: deptID, Position: null.StringFrom("Nurse"), CreatedAt: now, UpdatedAt: now}
	term := &model.Employee{ID: termID, ZitadelUserID: "zit-t", OrganizationID: orgID, DepartmentID: deptID, Position: null.StringFrom("Doctor"), CreatedAt: now, UpdatedAt: now}
	vac := &model.EmployeeVacation{
		ID: uuid.Must(uuid.NewV7()), EmployeeID: activeID,
		StartsAt: now.Add(-time.Hour), EndsAt: null.TimeFrom(now.Add(72 * time.Hour)),
		CreatedAt: now, UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qprojector.EmployeeHired(tx, active.ID.String(), active.CreatedAt, &empv1.EmployeeHired{
			ZitadelUserId:  active.ZitadelUserID,
			OrganizationId: active.OrganizationID.String(),
			DepartmentId:   active.DepartmentID.String(),
			Position:       active.Position.String,
			HiredAt:        timestamppb.New(active.CreatedAt),
		}); err != nil {
			return err
		}
		if err := qprojector.EmployeeHired(tx, term.ID.String(), term.CreatedAt, &empv1.EmployeeHired{
			ZitadelUserId:  term.ZitadelUserID,
			OrganizationId: term.OrganizationID.String(),
			DepartmentId:   term.DepartmentID.String(),
			Position:       term.Position.String,
			HiredAt:        timestamppb.New(term.CreatedAt),
		}); err != nil {
			return err
		}
		if err := qprojector.EmployeeTerminated(tx, term.ID.String(), now, &empv1.EmployeeTerminated{
			OrganizationId: term.OrganizationID.String(),
			DepartmentId:   term.DepartmentID.String(),
			TerminatedAt:   timestamppb.New(now),
		}); err != nil {
			return err
		}
		return qprojector.VacationStarted(tx, vac.EmployeeID.String(), vac.CreatedAt, &vacv1.VacationStarted{
			VacationId: vac.ID.String(),
			StartsAt:   timestamppb.New(vac.StartsAt),
			EndsAt:     timestamppb.New(vac.EndsAt.Time),
			CreatedAt:  timestamppb.New(vac.CreatedAt),
		})
	}))

	reader := memberread.NewEmployeeReader(testDB, authzSvc, &logger)

	// Default: terminated hidden → only the active employee.
	list, err := reader.ListByDepartment(ctx, sysadminCaller, deptID, memberread.ListQuery{}, memberread.EmployeeFilter{})
	require.NoError(t, err)
	require.Len(t, list.Items, 1)
	require.Equal(t, activeID, list.Items[0].EmployeeID)
	total, err := reader.CountByDepartment(ctx, sysadminCaller, deptID, memberread.EmployeeFilter{})
	require.NoError(t, err)
	require.Equal(t, int64(1), total)

	// IncludeTerminated=true → both.
	total, err = reader.CountByDepartment(ctx, sysadminCaller, deptID, memberread.EmployeeFilter{IncludeTerminated: true})
	require.NoError(t, err)
	require.Equal(t, int64(2), total)

	// OnVacation=true → only the active employee (term has no active vacation).
	list, err = reader.ListByDepartment(ctx, sysadminCaller, deptID, memberread.ListQuery{}, memberread.EmployeeFilter{OnVacation: true})
	require.NoError(t, err)
	require.Len(t, list.Items, 1)
	require.Equal(t, activeID, list.Items[0].EmployeeID)

	// Position filter: exact match.
	total, err = reader.CountByDepartment(ctx, sysadminCaller, deptID, memberread.EmployeeFilter{Position: "Nurse"})
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	total, err = reader.CountByDepartment(ctx, sysadminCaller, deptID, memberread.EmployeeFilter{Position: "Surgeon"})
	require.NoError(t, err)
	require.Equal(t, int64(0), total)
}

// TestEmployeeReader_SearchByOrganization seeds three employees in the
// same org with distinct display/email fields and asserts each of a
// full-name, partial-substring, and email-fragment query returns the
// expected row subset.
func TestEmployeeReader_SearchByOrganization(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)
	orgID, _, deptID := seedOrgClinicDept(t, ctx, now)

	seedUser := func(zit, first, last, display, email string) {
		require.NoError(t, testDB.WithContext(ctx).Exec(
			`INSERT INTO projections.users (id, user_name, first_name, last_name, display_name, email, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			zit, zit, first, last, display, email, now, now,
		).Error)
	}
	seedUser("zit-alice", "Alice", "Anders", "Alice Anders", "alice@example.org")
	seedUser("zit-bob", "Bob", "Brown", "Bob Brown", "bob@hospital.test")
	seedUser("zit-carol", "Carol", "Cooper", "Carol Cooper", "carol.cooper@clinic.test")

	mkEmp := func(zit string) uuid.UUID {
		id := uuid.Must(uuid.NewV7())
		emp := &model.Employee{ID: id, ZitadelUserID: zit, OrganizationID: orgID, DepartmentID: deptID, CreatedAt: now, UpdatedAt: now}
		require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return qprojector.EmployeeHired(tx, emp.ID.String(), emp.CreatedAt, &empv1.EmployeeHired{
				ZitadelUserId:  emp.ZitadelUserID,
				OrganizationId: emp.OrganizationID.String(),
				DepartmentId:   emp.DepartmentID.String(),
				HiredAt:        timestamppb.New(emp.CreatedAt),
			})
		}))
		return id
	}
	aliceID := mkEmp("zit-alice")
	bobID := mkEmp("zit-bob")
	carolID := mkEmp("zit-carol")

	reader := memberread.NewEmployeeReader(testDB, authzSvc, &logger)
	ids := func(r memberread.EmployeeListResult) map[uuid.UUID]bool {
		m := make(map[uuid.UUID]bool, len(r.Items))
		for _, v := range r.Items {
			m[v.EmployeeID] = true
		}
		return m
	}

	// Full-name match: "Alice Anders" → one result.
	r, err := reader.SearchByOrganization(ctx, sysadminCaller, orgID, "Alice Anders", memberread.ListQuery{}, memberread.EmployeeFilter{})
	require.NoError(t, err)
	require.Equal(t, map[uuid.UUID]bool{aliceID: true}, ids(r))

	// Partial substring "er" matches Anders and Cooper (last names).
	r, err = reader.SearchByOrganization(ctx, sysadminCaller, orgID, "er", memberread.ListQuery{}, memberread.EmployeeFilter{})
	require.NoError(t, err)
	require.Equal(t, map[uuid.UUID]bool{aliceID: true, carolID: true}, ids(r))

	// Email fragment "hospital" matches only Bob.
	r, err = reader.SearchByOrganization(ctx, sysadminCaller, orgID, "hospital", memberread.ListQuery{}, memberread.EmployeeFilter{})
	require.NoError(t, err)
	require.Equal(t, map[uuid.UUID]bool{bobID: true}, ids(r))
}

// TestEmployeeReader_ListVacationsByEmployee verifies scheduled
// vacations are visible + the state filter works.
func TestEmployeeReader_ListVacationsByEmployee(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID, _, deptID := seedOrgClinicDept(t, ctx, now)

	empID := uuid.Must(uuid.NewV7())
	emp := &model.Employee{
		ID: empID, ZitadelUserID: "zit-x", OrganizationID: orgID, DepartmentID: deptID,
		CreatedAt: now, UpdatedAt: now,
	}
	vacID := uuid.Must(uuid.NewV7())
	vac := &model.EmployeeVacation{
		ID: vacID, EmployeeID: empID,
		StartsAt: now.Add(24 * time.Hour), EndsAt: null.TimeFrom(now.Add(72 * time.Hour)),
		CreatedAt: now, UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qprojector.EmployeeHired(tx, emp.ID.String(), emp.CreatedAt, &empv1.EmployeeHired{
			ZitadelUserId:  emp.ZitadelUserID,
			OrganizationId: emp.OrganizationID.String(),
			DepartmentId:   emp.DepartmentID.String(),
			HiredAt:        timestamppb.New(emp.CreatedAt),
		}); err != nil {
			return err
		}
		return qprojector.VacationScheduled(tx, vac.EmployeeID.String(), vac.CreatedAt, &vacv1.VacationScheduled{
			VacationId: vac.ID.String(),
			StartsAt:   timestamppb.New(vac.StartsAt),
			EndsAt:     timestamppb.New(vac.EndsAt.Time),
			CreatedAt:  timestamppb.New(vac.CreatedAt),
		})
	}))

	reader := memberread.NewEmployeeReader(testDB, authzSvc, &logger)
	all, err := reader.ListVacationsByEmployee(ctx, sysadminCaller, empID, "", memberread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, all.Items, 1)
	require.Equal(t, "scheduled", all.Items[0].State)

	scheduled, err := reader.ListVacationsByEmployee(ctx, sysadminCaller, empID, "scheduled", memberread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, scheduled.Items, 1)

	none, err := reader.ListVacationsByEmployee(ctx, sysadminCaller, empID, "active", memberread.ListQuery{})
	require.NoError(t, err)
	require.Empty(t, none.Items)
}
