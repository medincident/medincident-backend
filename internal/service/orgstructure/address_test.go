package orgstructure

import (
	"errors"
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
)

func oopsCode(t *testing.T, err error) string {
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

func TestValidateAddressInput_TextEmpty(t *testing.T) {
	err := validateAddressInput(AddressInput{Text: ""})
	assert.Equal(t, ErrCodeAddressTextEmpty, oopsCode(t, err))
}

func TestValidateAddressInput_TextOnlyWhitespace(t *testing.T) {
	err := validateAddressInput(AddressInput{Text: "   "})
	assert.Equal(t, ErrCodeAddressTextEmpty, oopsCode(t, err))
}

func TestValidateAddressInput_TextTooShort(t *testing.T) {
	err := validateAddressInput(AddressInput{Text: "abc"}) // 3 runes, min is 4
	assert.Equal(t, ErrCodeAddressTextTooShort, oopsCode(t, err))
}

func TestValidateAddressInput_TextTooLong(t *testing.T) {
	long := make([]rune, addressMaxTextLen+1)
	for i := range long {
		long[i] = 'a'
	}
	err := validateAddressInput(AddressInput{Text: string(long)})
	assert.Equal(t, ErrCodeAddressTextTooLong, oopsCode(t, err))
}

func TestValidateAddressInput_Valid(t *testing.T) {
	err := validateAddressInput(AddressInput{Text: "г. Москва, ул. Пушкина"})
	assert.NoError(t, err)
}

func TestValidatePointInput_LongitudeOutOfRange(t *testing.T) {
	err := validatePointInput(PointInput{Longitude: 200, Latitude: 0})
	assert.Equal(t, ErrCodeAddressLongitudeOutOfRange, oopsCode(t, err))
}

func TestValidatePointInput_LatitudeOutOfRange(t *testing.T) {
	err := validatePointInput(PointInput{Longitude: 0, Latitude: -95})
	assert.Equal(t, ErrCodeAddressLatitudeOutOfRange, oopsCode(t, err))
}

func TestValidatePointInput_BothOutOfRange_JoinsErrors(t *testing.T) {
	err := validatePointInput(PointInput{Longitude: 200, Latitude: -95})
	unwrapper, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected joined error, got %T", err)
	}
	leaves := unwrapper.Unwrap()
	if len(leaves) != 2 {
		t.Fatalf("expected 2 joined leaves, got %d", len(leaves))
	}
	codes := []string{oopsCode(t, leaves[0]), oopsCode(t, leaves[1])}
	assert.Contains(t, codes, ErrCodeAddressLongitudeOutOfRange)
	assert.Contains(t, codes, ErrCodeAddressLatitudeOutOfRange)
}

func TestValidatePointInput_Valid(t *testing.T) {
	err := validatePointInput(PointInput{Longitude: 37.6, Latitude: 55.75})
	assert.NoError(t, err)
}
