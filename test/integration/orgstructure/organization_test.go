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
	"google.golang.org/protobuf/proto"

	organizationv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/organization/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	orgsvc "github.com/medincident/medincident-command-service/internal/services/orgstructure"
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
		Name:        "Клиника Пушкина",
		Description: &desc,
		LegalAddress: orgsvc.AddressInput{
			Text:  "г. Москва, ул. Пушкина, д. Колотушкина",
			Point: &orgsvc.PointInput{Longitude: 37.6, Latitude: 55.75},
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
	assert.True(t, row.LegalAddress.Point.Longitude.Valid)
	assert.InDelta(t, 37.6, row.LegalAddress.Point.Longitude.Float64, 0.0001)

	assert.Equal(t, int64(1), countRows(t, "outbox.events"))

	var event model.OutboxEvent
	require.NoError(t, testDB.First(&event, "subject = ?", orgsvc.SubjectOrganizationCreated).Error)

	var env envelopev1.Envelope
	require.NoError(t, proto.Unmarshal(event.Payload, &env))
	assert.Equal(t, "organization", env.AggregateType)
	assert.Equal(t, result.ID.String(), env.AggregateId)
	assert.NotNil(t, env.OccurredAt)

	var payload organizationv1.OrganizationCreated
	require.NoError(t, env.Payload.UnmarshalTo(&payload))
	assert.Equal(t, "Клиника Пушкина", payload.Name)
	require.NotNil(t, payload.Description)
	assert.Equal(t, desc, *payload.Description)
	require.NotNil(t, payload.LegalAddress)
	assert.Equal(t, "г. Москва, ул. Пушкина, д. Колотушкина", payload.LegalAddress.Text)
	require.NotNil(t, payload.LegalAddress.Point)
	assert.InDelta(t, 37.6, payload.LegalAddress.Point.Longitude, 0.0001)
}

func TestOrganization_Create_MultiFieldViolations(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	tooShortDesc := "tiny"

	_, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Name:        "",
		Description: &tooShortDesc,
		LegalAddress: orgsvc.AddressInput{
			Text:  "abc",
			Point: &orgsvc.PointInput{Longitude: 200, Latitude: -95},
		},
	})
	require.Error(t, err)

	unwrapper, ok := err.(interface{ Unwrap() []error })
	require.True(t, ok, "expected errors.Join-style error, got %T", err)
	leaves := unwrapper.Unwrap()

	codes := make(map[string]bool)
	for _, leaf := range leaves {
		// leaf may itself be a joined error (point validator returns
		// errors.Join of two leaves); flatten one level.
		if inner, ok := leaf.(interface{ Unwrap() []error }); ok {
			for _, sub := range inner.Unwrap() {
				codes[codeOf(t, sub)] = true
			}
			continue
		}
		codes[codeOf(t, leaf)] = true
	}

	assert.True(t, codes[orgsvc.ErrCodeOrganizationNameEmpty], "name_empty expected")
	assert.True(t, codes[orgsvc.ErrCodeOrganizationDescriptionTooShort], "description_too_short expected")
	assert.True(t, codes[orgsvc.ErrCodeAddressTextTooShort], "address_text_too_short expected")
	assert.True(t, codes[orgsvc.ErrCodeAddressLongitudeOutOfRange], "longitude_out_of_range expected")
	assert.True(t, codes[orgsvc.ErrCodeAddressLatitudeOutOfRange], "latitude_out_of_range expected")

	assert.Equal(t, int64(0), countRows(t, "domain.organizations"))
	assert.Equal(t, int64(0), countRows(t, "outbox.events"))
}

func TestOrganization_UpdateDetails_NoOp(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	desc := "Первое описание"

	created, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Name:         "Тестовая организация",
		Description:  &desc,
		LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), countRows(t, "outbox.events"))

	// Re-send the same values. Should be a no-op: no new outbox row.
	require.NoError(t, orgSvc.UpdateDetails(ctx, orgsvc.UpdateOrganizationDetailsCommand{
		ID:          created.ID,
		Name:        "Тестовая организация",
		Description: &desc,
	}))
	assert.Equal(t, int64(1), countRows(t, "outbox.events"))
}

func TestOrganization_UpdateDetails_RealChange(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	initial := "Начальное описание"
	changed := "Обновлённое описание с нормальной длиной"

	created, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Name:         "Орг А тестовая",
		Description:  &initial,
		LegalAddress: orgsvc.AddressInput{Text: "г. Москва, ул. Ленина, д. 1"},
	})
	require.NoError(t, err)

	require.NoError(t, orgSvc.UpdateDetails(ctx, orgsvc.UpdateOrganizationDetailsCommand{
		ID:          created.ID,
		Name:        "Орг А обновлённая",
		Description: &changed,
	}))
	assert.Equal(t, int64(2), countRows(t, "outbox.events"),
		"expected create + details_changed rows")

	var latest model.OutboxEvent
	require.NoError(t, testDB.Where("subject = ?", orgsvc.SubjectOrganizationDetailsChanged).
		First(&latest).Error)
	var env envelopev1.Envelope
	require.NoError(t, proto.Unmarshal(latest.Payload, &env))
	assert.Equal(t, "organization", env.AggregateType)
	assert.Equal(t, created.ID.String(), env.AggregateId)
}
