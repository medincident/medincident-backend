//go:build integration

package stats_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	statsread "github.com/medincident/medincident-backend/internal/service/query/stats"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
)

// TestReader_GetOrganizationStats returns counter values seeded by the
// projector chain.
func TestReader_GetOrganizationStats(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
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

	reader := statsread.NewReader(testDB, authzSvc, &logger)
	s, err := reader.GetOrganizationStats(ctx, sysadminCaller, orgID)
	require.NoError(t, err)
	require.Equal(t, int64(1), s.ClinicsTotal)
	require.Equal(t, int64(1), s.DepartmentsTotal)
	require.Equal(t, int64(0), s.EmployeesTotal)
}

// TestReader_GetClinicStats soft-miss yields zero-filled.
func TestReader_GetClinicStats_SoftMiss(t *testing.T) {
	resetProjections(t)
	logger := zerolog.Nop()
	reader := statsread.NewReader(testDB, authzSvc, &logger)
	s, err := reader.GetClinicStats(context.Background(), sysadminCaller, uuid.Must(uuid.NewV7()))
	require.NoError(t, err)
	require.Equal(t, int64(0), s.EmployeesTotal)
	require.Equal(t, int64(0), s.DepartmentsTotal)
	require.Equal(t, uuid.UUID{}, s.OrganizationID)
}

// TestReader_GetDepartmentStats returns the populated counter row.
func TestReader_GetDepartmentStats(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
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

	reader := statsread.NewReader(testDB, authzSvc, &logger)
	s, err := reader.GetDepartmentStats(ctx, sysadminCaller, deptID)
	require.NoError(t, err)
	require.Equal(t, orgID, s.OrganizationID)
	require.NotNil(t, s.ClinicID)
	require.Equal(t, clinicID, *s.ClinicID)
}
