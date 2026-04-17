//go:build integration

package orgstructure_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
)

// seedOrganization creates one organization and returns its id.
func seedOrganization(t *testing.T) uuid.UUID {
	t.Helper()
	res, err := orgSvc.Create(context.Background(), orgsvc.CreateOrganizationCommand{
		Name:         "Родительская организация",
		LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
	})
	require.NoError(t, err)
	return res.ID
}

func TestClinic_Create_HappyPath(t *testing.T) {
	resetDB(t)
	orgID := seedOrganization(t)

	res, err := clinSvc.Create(context.Background(), orgsvc.CreateClinicCommand{
		OrganizationID: orgID,
		Name:           "Клиника №1",
		PhysicalAddress: orgsvc.AddressInput{
			Text: "г. Москва, ул. Тверская, д. 10",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	assert.Equal(t, 1, countClinics(t))
	// organization + clinic = 2 outbox rows
	assert.Equal(t, 2, countOutboxEvents(t))
}

func TestClinic_Create_OrganizationNotFound(t *testing.T) {
	resetDB(t)

	_, err := clinSvc.Create(context.Background(), orgsvc.CreateClinicCommand{
		OrganizationID:  uuid.New(),
		Name:            "Сирота клиника",
		PhysicalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Тверская, д. 10"},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeClinicOrganizationNotFound, codeOf(t, err))
	assert.Equal(t, 0, countClinics(t))
	assert.Equal(t, 0, countOutboxEvents(t))
}
