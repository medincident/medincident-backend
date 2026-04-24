//go:build integration

package orgstructure_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
)

// seedOrganization creates one organization and returns its id.
func seedOrganization(t *testing.T) uuid.UUID {
	t.Helper()
	res, err := orgSvc.Create(context.Background(), orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:         "Родительская организация",
			LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
		},
	})
	require.NoError(t, err)
	return res.ID
}

func TestClinic_Create_HappyPath(t *testing.T) {
	resetDB(t)
	orgID := seedOrganization(t)

	res, err := clinSvc.Create(context.Background(), orgsvc.CreateClinicCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{
			OrganizationID: orgID.String(),
			Name:           "Клиника №1",
			PhysicalAddress: orgsvc.AddressInput{
				Text: "г. Москва, ул. Тверская, д. 10",
			},
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	assert.Equal(t, 1, countClinics(t))
	// Projection row written inline with the domain row.
	assert.Equal(t, 1, countProjectionClinics(t))
}

func TestClinic_Create_OrganizationNotFound(t *testing.T) {
	resetDB(t)

	_, err := clinSvc.Create(context.Background(), orgsvc.CreateClinicCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{
			OrganizationID:  uuid.New().String(),
			Name:            "Сирота клиника",
			PhysicalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Тверская, д. 10"},
		},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeClinicOrganizationNotFound, codeOf(t, err))
	assert.Equal(t, 0, countClinics(t))
	assert.Equal(t, 0, countProjectionClinics(t))
}
