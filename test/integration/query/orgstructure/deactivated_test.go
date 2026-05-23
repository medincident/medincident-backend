//go:build integration

package orgstructure_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	orgread "github.com/medincident/medincident-backend/internal/service/query/orgstructure"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
)

// noRolesCaller is authenticated but has no row in any domain role table.
var noRolesCaller = authz.Caller{ZitadelUserID: "deactivated-test-noroles"}

// TestOrganizationReader_Get_DeactivatedVisibility seeds an inactive
// organization and verifies that the system admin sees it while a caller
// with no roles receives not_found.
func TestOrganizationReader_Get_DeactivatedVisibility(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qprojector.OrganizationCreated(tx, orgID.String(), now, &orgv1.OrganizationCreated{
			Name:         "Inactive Org",
			LegalAddress: &orgv1.Address{Text: "Street"},
			CreatedAt:    timestamppb.New(now),
		}); err != nil {
			return err
		}
		return tx.Exec(`UPDATE projections.organizations SET is_active = FALSE WHERE id = ?`, orgID).Error
	}))

	reader := orgread.NewOrganizationReader(testDB, authzSvc, &logger)

	// Sysadmin sees the deactivated record.
	view, err := reader.Get(ctx, sysadminCaller, orgID)
	require.NoError(t, err)
	require.Equal(t, orgID, view.ID)

	// Non-admin caller receives not_found (same as for a missing record).
	_, err = reader.Get(ctx, noRolesCaller, orgID)
	require.Error(t, err)
}

// TestOrganizationReader_List_IncludeDeactivated seeds one inactive and one
// active organization. Verifies:
//   - includeDeactivated=false: only the active org is returned.
//   - includeDeactivated=true (sysadmin): both orgs are returned.
//   - includeDeactivated=true (no-roles caller): permission denied.
func TestOrganizationReader_List_IncludeDeactivated(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	activeID := uuid.Must(uuid.NewV7())
	inactiveID := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qprojector.OrganizationCreated(tx, activeID.String(), now, &orgv1.OrganizationCreated{
			Name:         "Active Org",
			LegalAddress: &orgv1.Address{Text: "Street"},
			CreatedAt:    timestamppb.New(now),
		}); err != nil {
			return err
		}
		if err := qprojector.OrganizationCreated(tx, inactiveID.String(), now.Add(time.Second), &orgv1.OrganizationCreated{
			Name:         "Inactive Org",
			LegalAddress: &orgv1.Address{Text: "Street"},
			CreatedAt:    timestamppb.New(now.Add(time.Second)),
		}); err != nil {
			return err
		}
		return tx.Exec(`UPDATE projections.organizations SET is_active = FALSE WHERE id = ?`, inactiveID).Error
	}))

	reader := orgread.NewOrganizationReader(testDB, authzSvc, &logger)

	// includeDeactivated=false: only the active org.
	result, err := reader.List(ctx, sysadminCaller, false, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Len(t, result.Items, 1)
	require.Equal(t, activeID, result.Items[0].ID)

	// includeDeactivated=true (admin): both orgs.
	result, err = reader.List(ctx, sysadminCaller, true, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Len(t, result.Items, 2)

	// includeDeactivated=true (no-roles caller): permission denied.
	_, err = reader.List(ctx, noRolesCaller, true, orgread.ListQuery{Limit: 10})
	require.Error(t, err)
}

// TestOrganizationReader_Search_IncludeDeactivated verifies that
// include_deactivated=true in Search requires system admin.
func TestOrganizationReader_Search_IncludeDeactivated(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qprojector.OrganizationCreated(tx, orgID.String(), now, &orgv1.OrganizationCreated{
			Name:         "Dormant Health",
			LegalAddress: &orgv1.Address{Text: "Street"},
			CreatedAt:    timestamppb.New(now),
		}); err != nil {
			return err
		}
		return tx.Exec(`UPDATE projections.organizations SET is_active = FALSE WHERE id = ?`, orgID).Error
	}))

	reader := orgread.NewOrganizationReader(testDB, authzSvc, &logger)

	// includeDeactivated=false: inactive org is excluded from search.
	result, err := reader.Search(ctx, sysadminCaller, "dormant", false, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Empty(t, result.Items)

	// includeDeactivated=true (admin): inactive org appears in search.
	result, err = reader.Search(ctx, sysadminCaller, "dormant", true, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Len(t, result.Items, 1)
	require.Equal(t, orgID, result.Items[0].ID)

	// includeDeactivated=true (no-roles caller): permission denied.
	_, err = reader.Search(ctx, noRolesCaller, "dormant", true, orgread.ListQuery{Limit: 10})
	require.Error(t, err)
}

// TestClinicReader_IncludeDeactivated seeds an inactive clinic and verifies
// it is hidden by default (includeDeactivated=false) and visible to the
// system admin when includeDeactivated=true.
func TestClinicReader_IncludeDeactivated(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())
	activeID := uuid.Must(uuid.NewV7())
	inactiveID := uuid.Must(uuid.NewV7())

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qprojector.OrganizationCreated(tx, orgID.String(), now, &orgv1.OrganizationCreated{
			Name:         "O",
			LegalAddress: &orgv1.Address{Text: "S"},
			CreatedAt:    timestamppb.New(now),
		}); err != nil {
			return err
		}
		if err := qprojector.ClinicCreated(tx, activeID.String(), now, &clinicv1.ClinicCreated{
			OrganizationId:  orgID.String(),
			Name:            "Active Clinic",
			PhysicalAddress: &clinicv1.Address{Text: "S"},
			CreatedAt:       timestamppb.New(now),
		}); err != nil {
			return err
		}
		if err := qprojector.ClinicCreated(tx, inactiveID.String(), now.Add(time.Second), &clinicv1.ClinicCreated{
			OrganizationId:  orgID.String(),
			Name:            "Inactive Clinic",
			PhysicalAddress: &clinicv1.Address{Text: "S"},
			CreatedAt:       timestamppb.New(now.Add(time.Second)),
		}); err != nil {
			return err
		}
		return tx.Exec(`UPDATE projections.clinics SET is_active = FALSE WHERE id = ?`, inactiveID).Error
	}))

	reader := orgread.NewClinicReader(testDB, authzSvc, &logger)

	// includeDeactivated=false: only the active clinic.
	result, err := reader.ListByOrganization(ctx, sysadminCaller, orgID, false, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Len(t, result.Items, 1)
	require.Equal(t, activeID, result.Items[0].ID)

	// includeDeactivated=true (admin): both clinics.
	result, err = reader.ListByOrganization(ctx, sysadminCaller, orgID, true, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Len(t, result.Items, 2)
}

// TestDepartmentReader_IncludeDeactivated seeds an inactive department and
// verifies it is hidden by default and visible to the admin with
// includeDeactivated=true.
func TestDepartmentReader_IncludeDeactivated(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	activeID := uuid.Must(uuid.NewV7())
	inactiveID := uuid.Must(uuid.NewV7())

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qprojector.OrganizationCreated(tx, orgID.String(), now, &orgv1.OrganizationCreated{
			Name:         "O",
			LegalAddress: &orgv1.Address{Text: "S"},
			CreatedAt:    timestamppb.New(now),
		}); err != nil {
			return err
		}
		if err := qprojector.ClinicCreated(tx, clinicID.String(), now, &clinicv1.ClinicCreated{
			OrganizationId:  orgID.String(),
			Name:            "C",
			PhysicalAddress: &clinicv1.Address{Text: "S"},
			CreatedAt:       timestamppb.New(now),
		}); err != nil {
			return err
		}
		if err := qprojector.DepartmentCreated(tx, activeID.String(), now, &deptv1.DepartmentCreated{
			ClinicId:  clinicID.String(),
			Name:      "Active Dept",
			CreatedAt: timestamppb.New(now),
		}); err != nil {
			return err
		}
		if err := qprojector.DepartmentCreated(tx, inactiveID.String(), now.Add(time.Second), &deptv1.DepartmentCreated{
			ClinicId:  clinicID.String(),
			Name:      "Inactive Dept",
			CreatedAt: timestamppb.New(now.Add(time.Second)),
		}); err != nil {
			return err
		}
		return tx.Exec(`UPDATE projections.departments SET is_active = FALSE WHERE id = ?`, inactiveID).Error
	}))

	reader := orgread.NewDepartmentReader(testDB, authzSvc, &logger)

	// includeDeactivated=false: only the active dept.
	result, err := reader.ListByClinic(ctx, sysadminCaller, clinicID, false, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Len(t, result.Items, 1)
	require.Equal(t, activeID, result.Items[0].ID)

	// includeDeactivated=true (admin): both depts.
	result, err = reader.ListByClinic(ctx, sysadminCaller, clinicID, true, orgread.ListQuery{Limit: 10})
	require.NoError(t, err)
	require.Len(t, result.Items, 2)
}
