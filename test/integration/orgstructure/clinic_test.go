//go:build integration

package orgstructure_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/model"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
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
}

func TestClinic_UpdatePhysicalAddress_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	err := clinSvc.UpdatePhysicalAddress(ctx, orgsvc.UpdateClinicPhysicalAddressCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateClinicPhysicalAddressPayload{
			ID:      clinicID.String(),
			Address: orgsvc.AddressInput{Text: "г. Москва, ул. Новая, д. 5"},
		},
	})
	require.NoError(t, err)

	var clinic model.Clinic
	require.NoError(t, testDB.First(&clinic, "id = ?", clinicID).Error)
	assert.Equal(t, "г. Москва, ул. Новая, д. 5", clinic.PhysicalAddress.Text)
	assert.False(t, clinic.PhysicalAddress.Point.Valid, "point must be absent")

	env := outboxEnvelope(t, "medincident.event.clinic.v1.physical_address_changed")
	var msg clinicv1.ClinicPhysicalAddressChanged
	require.NoError(t, env.Payload.UnmarshalTo(&msg))
	assert.Equal(t, clinicID.String(), env.AggregateId)
	assert.Equal(t, "г. Москва, ул. Новая, д. 5", msg.GetPhysicalAddress().GetText())
	assert.Nil(t, msg.GetPhysicalAddress().GetPoint())
}

func TestClinic_UpdatePhysicalAddress_HappyPath_WithPoint(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	err := clinSvc.UpdatePhysicalAddress(ctx, orgsvc.UpdateClinicPhysicalAddressCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateClinicPhysicalAddressPayload{
			ID: clinicID.String(),
			Address: orgsvc.AddressInput{
				Text:  "г. Сочи, ул. Сириуса, д. 1",
				Point: &orgsvc.PointInput{Longitude: 39.72, Latitude: 43.58},
			},
		},
	})
	require.NoError(t, err)

	var clinic model.Clinic
	require.NoError(t, testDB.First(&clinic, "id = ?", clinicID).Error)
	assert.Equal(t, "г. Сочи, ул. Сириуса, д. 1", clinic.PhysicalAddress.Text)
	require.True(t, clinic.PhysicalAddress.Point.Valid)
	assert.InDelta(t, 39.72, clinic.PhysicalAddress.Point.V.Longitude, 0.0001)
	assert.InDelta(t, 43.58, clinic.PhysicalAddress.Point.V.Latitude, 0.0001)

	env := outboxEnvelope(t, "medincident.event.clinic.v1.physical_address_changed")
	var msg clinicv1.ClinicPhysicalAddressChanged
	require.NoError(t, env.Payload.UnmarshalTo(&msg))
	assert.Equal(t, clinicID.String(), env.AggregateId)
	assert.Equal(t, "г. Сочи, ул. Сириуса, д. 1", msg.GetPhysicalAddress().GetText())
	require.NotNil(t, msg.GetPhysicalAddress().GetPoint())
	assert.InDelta(t, 39.72, msg.GetPhysicalAddress().GetPoint().GetLongitude(), 0.0001)
	assert.InDelta(t, 43.58, msg.GetPhysicalAddress().GetPoint().GetLatitude(), 0.0001)
}

func TestClinic_UpdatePhysicalAddress_NoOp(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	clinicID := seedClinic(t)

	err := clinSvc.UpdatePhysicalAddress(ctx, orgsvc.UpdateClinicPhysicalAddressCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateClinicPhysicalAddressPayload{
			ID:      clinicID.String(),
			Address: orgsvc.AddressInput{Text: "г. Москва, ул. Тверская, д. 10"},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, 0, countOutboxEvents(t, "medincident.event.clinic.v1.physical_address_changed"))
}

func TestClinic_UpdatePhysicalAddress_ClinicNotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := clinSvc.UpdatePhysicalAddress(ctx, orgsvc.UpdateClinicPhysicalAddressCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateClinicPhysicalAddressPayload{
			ID:      uuid.New().String(),
			Address: orgsvc.AddressInput{Text: "г. Москва, ул. Новая, д. 5"},
		},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeClinicNotFound, codeOf(t, err))
	assert.Equal(t, 0, countClinics(t))
	assert.Equal(t, 0, countOutboxEvents(t, "medincident.event.clinic.v1.physical_address_changed"))
}
