//go:build integration

package orgstructure_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orgsvc "github.com/medincident/medincident-command-service/internal/service/command/orgstructure"
)

// seedClinic creates one organization + one clinic and returns the
// clinic id.
func seedClinic(t *testing.T) uuid.UUID {
	t.Helper()
	orgID := seedOrganization(t)
	res, err := clinSvc.Create(context.Background(), orgsvc.CreateClinicCommand{
		OrganizationID:  orgID,
		Name:            "Родительская клиника",
		PhysicalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Тверская, д. 10"},
	})
	require.NoError(t, err)
	return res.ID
}

func TestDepartment_Create_HappyPath(t *testing.T) {
	resetDB(t)
	clinicID := seedClinic(t)

	res, err := deptSvc.Create(context.Background(), orgsvc.CreateDepartmentCommand{
		ClinicID: clinicID,
		Name:     "Терапия",
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	assert.Equal(t, 1, countDepartments(t))
	// organization + clinic + department = 3 outbox rows
	assert.Equal(t, 3, countOutboxEvents(t))
}

func TestDepartment_Create_ClinicNotFound(t *testing.T) {
	resetDB(t)

	_, err := deptSvc.Create(context.Background(), orgsvc.CreateDepartmentCommand{
		ClinicID: uuid.New(),
		Name:     "Сирота отделение",
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeDepartmentClinicNotFound, codeOf(t, err))
	assert.Equal(t, 0, countDepartments(t))
	assert.Equal(t, 0, countOutboxEvents(t))
}
