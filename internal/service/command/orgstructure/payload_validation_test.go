package orgstructure_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// violationsOf extracts the []Violation slice from a validation_failed
// oops error returned by validation.Struct.
func violationsOf(t *testing.T, err error) []validation.Violation {
	t.Helper()
	require.Error(t, err)
	var oe oops.OopsError
	require.True(t, errors.As(err, &oe), "error is not an oops error: %v", err)
	require.Equal(t, validation.CodeValidationFailed, oe.Code())
	vs, ok := oe.Context()[validation.ContextKeyViolations].([]validation.Violation)
	require.True(t, ok, "expected []validation.Violation in context, got %T", oe.Context()[validation.ContextKeyViolations])
	return vs
}

// hasViolation reports whether vs contains a violation matching field
// path and rule.
func hasViolation(vs []validation.Violation, field, rule string) bool {
	for _, v := range vs {
		if v.Field == field && v.Rule == rule {
			return true
		}
	}
	return false
}

// TestCreateOrganizationPayload_NestedAddressValidation guards the
// regression where AddressInput's struct-tag rules were silently
// skipped because the parent field had no `validate:"required"` tag.
// validator.WithRequiredStructEnabled requires the nested struct to be
// reachable for descent; without `required` on LegalAddress, an empty
// AddressInput would slip past validation entirely.
func TestCreateOrganizationPayload_NestedAddressValidation(t *testing.T) {
	t.Parallel()

	t.Run("zero address fails parent required", func(t *testing.T) {
		t.Parallel()
		// A zero AddressInput stops at the parent `required` rule on
		// LegalAddress — without that tag the validator would not even
		// reach the nested fields, which is the regression this test
		// guards against.
		err := validation.Struct(orgsvc.CreateOrganizationPayload{
			Name:         "Valid Name",
			LegalAddress: orgsvc.AddressInput{},
		})
		vs := violationsOf(t, err)
		assert.True(t, hasViolation(vs, "legal_address", "required"),
			"expected legal_address required violation, got %+v", vs)
	})

	t.Run("too-short text fails min", func(t *testing.T) {
		t.Parallel()
		err := validation.Struct(orgsvc.CreateOrganizationPayload{
			Name:         "Valid Name",
			LegalAddress: orgsvc.AddressInput{Text: "abc"},
		})
		vs := violationsOf(t, err)
		assert.True(t, hasViolation(vs, "legal_address.text", "min"),
			"expected legal_address.text min violation, got %+v", vs)
	})

	t.Run("out-of-range point coords fail max", func(t *testing.T) {
		t.Parallel()
		err := validation.Struct(orgsvc.CreateOrganizationPayload{
			Name: "Valid Name",
			LegalAddress: orgsvc.AddressInput{
				Text:  "г. Москва, ул. Тверская, д. 1",
				Point: &orgsvc.PointInput{Longitude: 200, Latitude: 100},
			},
		})
		vs := violationsOf(t, err)
		assert.True(t, hasViolation(vs, "legal_address.point.longitude", "max"),
			"expected legal_address.point.longitude max violation, got %+v", vs)
		assert.True(t, hasViolation(vs, "legal_address.point.latitude", "max"),
			"expected legal_address.point.latitude max violation, got %+v", vs)
	})

	t.Run("valid payload passes", func(t *testing.T) {
		t.Parallel()
		err := validation.Struct(orgsvc.CreateOrganizationPayload{
			Name: "Valid Name",
			LegalAddress: orgsvc.AddressInput{
				Text:  "г. Москва, ул. Тверская, д. 1",
				Point: &orgsvc.PointInput{Longitude: 37.6, Latitude: 55.75},
			},
		})
		require.NoError(t, err)
	})
}

// TestUpdateOrganizationLegalAddressPayload_NestedAddressValidation
// covers the same regression on the update path.
func TestUpdateOrganizationLegalAddressPayload_NestedAddressValidation(t *testing.T) {
	t.Parallel()

	err := validation.Struct(orgsvc.UpdateOrganizationLegalAddressPayload{
		ID:      uuid.New().String(),
		Address: orgsvc.AddressInput{Text: "abc"},
	})
	vs := violationsOf(t, err)
	assert.True(t, hasViolation(vs, "address.text", "min"),
		"expected address.text min violation, got %+v", vs)
}

// TestCreateClinicPayload_NestedAddressValidation covers the clinic
// create path.
func TestCreateClinicPayload_NestedAddressValidation(t *testing.T) {
	t.Parallel()

	err := validation.Struct(orgsvc.CreateClinicPayload{
		OrganizationID:  uuid.New().String(),
		Name:            "Valid Clinic",
		PhysicalAddress: orgsvc.AddressInput{Text: "abc"},
	})
	vs := violationsOf(t, err)
	assert.True(t, hasViolation(vs, "physical_address.text", "min"),
		"expected physical_address.text min violation, got %+v", vs)
}

// TestUpdateClinicPhysicalAddressPayload_NestedAddressValidation covers
// the clinic update-address path.
func TestUpdateClinicPhysicalAddressPayload_NestedAddressValidation(t *testing.T) {
	t.Parallel()

	err := validation.Struct(orgsvc.UpdateClinicPhysicalAddressPayload{
		ID:      uuid.New().String(),
		Address: orgsvc.AddressInput{Text: "abc"},
	})
	vs := violationsOf(t, err)
	assert.True(t, hasViolation(vs, "address.text", "min"),
		"expected address.text min violation, got %+v", vs)
}
