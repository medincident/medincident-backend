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

	// Projection row written atomically alongside the domain row.
	assert.Equal(t, 1, countProjectionOrganizations(t))

	var projName string
	require.NoError(t, testDB.Raw(
		`SELECT name FROM projections.organizations WHERE id = ?`, result.ID,
	).Row().Scan(&projName))
	assert.Equal(t, "Клиника Пушкина", projName)
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

	unwrapper, ok := err.(interface{ Unwrap() []error })
	require.True(t, ok, "expected errors.Join-style error, got %T", err)
	leaves := unwrapper.Unwrap()

	codes := make(map[string]bool)
	var collect func(error)
	collect = func(e error) {
		if inner, ok := e.(interface{ Unwrap() []error }); ok {
			for _, sub := range inner.Unwrap() {
				collect(sub)
			}
			return
		}
		codes[codeOf(t, e)] = true
	}
	for _, leaf := range leaves {
		collect(leaf)
	}

	assert.True(t, codes[validation.CodeStringRequired], "string_required expected for empty name")
	assert.True(t, codes[validation.CodeStringTooShort], "string_too_short expected for short description / address text")
	assert.True(t, codes[validation.CodeFloatOutOfRange], "float_out_of_range expected for bad coordinates")

	assert.Equal(t, 0, countOrganizations(t))
	assert.Equal(t, 0, countProjectionOrganizations(t))
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

	// Capture the projection's updated_at before the no-op.
	var firstUpdatedAt string
	require.NoError(t, testDB.Raw(
		`SELECT updated_at::text FROM projections.organizations WHERE id = ?`, created.ID,
	).Row().Scan(&firstUpdatedAt))

	// Re-send the same values. No-op means no projection write.
	require.NoError(t, orgSvc.UpdateDetails(ctx, orgsvc.UpdateOrganizationDetailsCommand{
		Caller: sysadminCaller,
		Payload: orgsvc.UpdateOrganizationDetailsPayload{
			ID:          created.ID.String(),
			Name:        "Тестовая организация",
			Description: &desc,
		},
	}))

	var secondUpdatedAt string
	require.NoError(t, testDB.Raw(
		`SELECT updated_at::text FROM projections.organizations WHERE id = ?`, created.ID,
	).Row().Scan(&secondUpdatedAt))
	assert.Equal(t, firstUpdatedAt, secondUpdatedAt, "expected no-op to leave updated_at alone")
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

	var projName, projDesc string
	require.NoError(t, testDB.Raw(
		`SELECT name, description FROM projections.organizations WHERE id = ?`, created.ID,
	).Row().Scan(&projName, &projDesc))
	assert.Equal(t, "Орг А обновлённая", projName)
	assert.Equal(t, changed, projDesc)
}
