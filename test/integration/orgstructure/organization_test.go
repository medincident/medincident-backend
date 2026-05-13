//go:build integration

package orgstructure_integration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/model"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	"github.com/medincident/medincident-backend/internal/service/validation"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
)

// codeOf extracts the oops Code as a string from any error in a joined
// multi-error tree.
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

func TestOrganization_Create_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	desc := "Крупнейшая частная клиника региона"

	result, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:        "Клиника Пушкина",
			Description: &desc,
			LegalAddress: orgsvc.AddressInput{
				Text:  "г. Москва, ул. Пушкина, д. Колотушкина",
				Point: &orgsvc.PointInput{Longitude: 37.6, Latitude: 55.75},
			},
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, result.ID)

	var row model.Organization
	require.NoError(t, testDB.First(&row, "id = ?", result.ID).Error)
	assert.Equal(t, "Клиника Пушкина", row.Name)
	assert.True(t, row.Description.Valid)
	assert.Equal(t, desc, row.Description.String)
	assert.Equal(t, "г. Москва, ул. Пушкина, д. Колотушкина", row.LegalAddress.Text)
	require.True(t, row.LegalAddress.Point.Valid)
	assert.InDelta(t, 37.6, row.LegalAddress.Point.V.Longitude, 0.0001)
}

func TestOrganization_Create_MultiFieldViolations(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	tooShortDesc := "tiny"

	_, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:        "",
			Description: &tooShortDesc,
			LegalAddress: orgsvc.AddressInput{
				Text:  "abc",
				Point: &orgsvc.PointInput{Longitude: 200, Latitude: -95},
			},
		},
	})
	require.Error(t, err)

	var oe oops.OopsError
	require.True(t, errors.As(err, &oe), "expected oops error, got %T", err)
	require.Equal(t, validation.CodeValidationFailed, oe.Code())
	violations, ok := oe.Context()[validation.ContextKeyViolations].([]validation.Violation)
	require.True(t, ok, "expected []validation.Violation in context, got %T", oe.Context()[validation.ContextKeyViolations])

	rules := map[string]string{}
	for _, v := range violations {
		rules[v.Field] = v.Rule
	}
	assert.Equal(t, "required", rules["name"], "name should fail required, got violations=%+v", violations)
	assert.Equal(t, "min", rules["description"], "description should fail min, got violations=%+v", violations)
	assert.Equal(t, "min", rules["legal_address.text"], "legal_address.text should fail min, got violations=%+v", violations)
	assert.Equal(t, "max", rules["legal_address.point.longitude"], "longitude should fail max, got violations=%+v", violations)
	assert.Equal(t, "min", rules["legal_address.point.latitude"], "latitude should fail min, got violations=%+v", violations)

	assert.Equal(t, 0, countOrganizations(t))
}

func TestOrganization_UpdateDetails_NoOp(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	desc := "Первое описание"

	created, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:         "Тестовая организация",
			Description:  &desc,
			LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
		},
	})
	require.NoError(t, err)

	// Re-send the same values — service must succeed (no-op on domain side).
	require.NoError(t, orgSvc.UpdateDetails(ctx, orgsvc.UpdateOrganizationDetailsCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateOrganizationDetailsPayload{
			ID:          created.ID.String(),
			Name:        "Тестовая организация",
			Description: &desc,
		},
	}))
}

func TestOrganization_UpdateDetails_RealChange(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	initial := "Начальное описание"
	changed := "Обновлённое описание с нормальной длиной"

	created, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{
			Name:         "Орг А тестовая",
			Description:  &initial,
			LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
		},
	})
	require.NoError(t, err)

	require.NoError(t, orgSvc.UpdateDetails(ctx, orgsvc.UpdateOrganizationDetailsCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateOrganizationDetailsPayload{
			ID:          created.ID.String(),
			Name:        "Орг А обновлённая",
			Description: &changed,
		},
	}))

	// Verify the domain row was updated.
	var domainOrg struct{ Name, Description string }
	require.NoError(t, testDB.Raw(
		`SELECT name, description FROM domain.organizations WHERE id = ?`, created.ID,
	).Row().Scan(&domainOrg.Name, &domainOrg.Description))
	assert.Equal(t, "Орг А обновлённая", domainOrg.Name)
	assert.Equal(t, changed, domainOrg.Description)
}

func TestOrganization_UpdateLegalAddress_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	err := orgSvc.UpdateLegalAddress(ctx, orgsvc.UpdateOrganizationLegalAddressCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateOrganizationLegalAddressPayload{
			ID:      orgID.String(),
			Address: orgsvc.AddressInput{Text: "г. Казань, ул. Баумана, д. 10"},
		},
	})
	require.NoError(t, err)

	var org model.Organization
	require.NoError(t, testDB.First(&org, "id = ?", orgID).Error)
	assert.Equal(t, "г. Казань, ул. Баумана, д. 10", org.LegalAddress.Text)
	assert.False(t, org.LegalAddress.Point.Valid, "point must be absent")

	env := outboxEnvelope(t, "medincident.event.organization.v1.legal_address_changed")
	var msg orgv1.OrganizationLegalAddressChanged
	require.NoError(t, env.Payload.UnmarshalTo(&msg))
	assert.Equal(t, orgID.String(), env.AggregateId)
	assert.Equal(t, "г. Казань, ул. Баумана, д. 10", msg.GetLegalAddress().GetText())
	assert.Nil(t, msg.GetLegalAddress().GetPoint())
}

func TestOrganization_UpdateLegalAddress_HappyPath_WithPoint(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	err := orgSvc.UpdateLegalAddress(ctx, orgsvc.UpdateOrganizationLegalAddressCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateOrganizationLegalAddressPayload{
			ID: orgID.String(),
			Address: orgsvc.AddressInput{
				Text:  "г. Москва, ул. Пушкина, д. 1",
				Point: &orgsvc.PointInput{Longitude: 37.6, Latitude: 55.75},
			},
		},
	})
	require.NoError(t, err)

	var org model.Organization
	require.NoError(t, testDB.First(&org, "id = ?", orgID).Error)
	assert.Equal(t, "г. Москва, ул. Пушкина, д. 1", org.LegalAddress.Text)
	require.True(t, org.LegalAddress.Point.Valid)
	assert.InDelta(t, 37.6, org.LegalAddress.Point.V.Longitude, 0.0001)
	assert.InDelta(t, 55.75, org.LegalAddress.Point.V.Latitude, 0.0001)

	env := outboxEnvelope(t, "medincident.event.organization.v1.legal_address_changed")
	var msg orgv1.OrganizationLegalAddressChanged
	require.NoError(t, env.Payload.UnmarshalTo(&msg))
	assert.Equal(t, orgID.String(), env.AggregateId)
	assert.Equal(t, "г. Москва, ул. Пушкина, д. 1", msg.GetLegalAddress().GetText())
	require.NotNil(t, msg.GetLegalAddress().GetPoint())
	assert.InDelta(t, 37.6, msg.GetLegalAddress().GetPoint().GetLongitude(), 0.0001)
	assert.InDelta(t, 55.75, msg.GetLegalAddress().GetPoint().GetLatitude(), 0.0001)
}

func TestOrganization_UpdateLegalAddress_NoOp(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrganization(t)

	err := orgSvc.UpdateLegalAddress(ctx, orgsvc.UpdateOrganizationLegalAddressCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateOrganizationLegalAddressPayload{
			ID:      orgID.String(),
			Address: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, 0, countOutboxEvents(t, "medincident.event.organization.v1.legal_address_changed"))
}

func TestOrganization_UpdateLegalAddress_OrganizationNotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	err := orgSvc.UpdateLegalAddress(ctx, orgsvc.UpdateOrganizationLegalAddressCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateOrganizationLegalAddressPayload{
			ID:      uuid.New().String(),
			Address: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
		},
	})
	require.Error(t, err)
	assert.Equal(t, orgsvc.ErrCodeOrganizationNotFound, codeOf(t, err))
	assert.Equal(t, 0, countOutboxEvents(t, "medincident.event.organization.v1.legal_address_changed"))
	assert.Equal(t, 0, countOrganizations(t))
}
