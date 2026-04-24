package validation_test

import (
	"errors"
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/validation"
)

type noExtraWSPayload struct {
	Value string `validate:"no_extra_ws"`
}

func TestStruct_NoExtraWhitespace(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		fail bool
	}{
		{"empty passes", "", false},
		{"single word", "alpha", false},
		{"single space between words", "alpha beta", false},
		{"multiple single-space separated words", "one two three", false},
		{"leading space", " alpha", true},
		{"trailing space", "alpha ", true},
		{"only spaces", "   ", true},
		{"double space inside", "alpha  beta", true},
		{"tab inside run with space", "alpha \tbeta", true},
		{"leading tab", "\talpha", true},
		{"trailing newline", "alpha\n", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := validation.Struct(noExtraWSPayload{Value: tc.in})
			if !tc.fail {
				require.NoError(t, err)
				return
			}
			vs := violationsOf(t, err)
			require.Len(t, vs, 1)
			assert.Equal(t, "value", vs[0].Field)
			assert.Equal(t, validation.TagNoExtraWhitespace, vs[0].Rule)
		})
	}
}

// violationsOf pulls the []Violation slice from an oops error emitted
// by validation.Struct. Fails the test if err is not a single
// validation_failed oops leaf.
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

type orderPayload struct {
	Value string `validate:"required,no_extra_ws,min=4"`
}

func TestStruct_OrderingWithOtherRules(t *testing.T) {
	t.Parallel()

	vs := violationsOf(t, validation.Struct(orderPayload{Value: ""}))
	require.Len(t, vs, 1)
	assert.Equal(t, "value", vs[0].Field)
	assert.Equal(t, "required", vs[0].Rule)

	vs = violationsOf(t, validation.Struct(orderPayload{Value: " abc"}))
	require.Len(t, vs, 1)
	assert.Equal(t, "value", vs[0].Field)
	assert.Equal(t, validation.TagNoExtraWhitespace, vs[0].Rule)

	vs = violationsOf(t, validation.Struct(orderPayload{Value: "abc"}))
	require.Len(t, vs, 1)
	assert.Equal(t, "value", vs[0].Field)
	assert.Equal(t, "min", vs[0].Rule)

	require.NoError(t, validation.Struct(orderPayload{Value: "alpha beta"}))
}

type pointerPayload struct {
	Value *string `validate:"omitnil,no_extra_ws,min=2"`
}

func TestStruct_NoExtraWhitespace_PointerField(t *testing.T) {
	t.Parallel()

	strPtr := func(s string) *string { return &s }

	t.Run("nil pointer skips rule", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, validation.Struct(pointerPayload{Value: nil}))
	})

	t.Run("clean non-nil passes", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, validation.Struct(pointerPayload{Value: strPtr("alpha beta")}))
	})

	t.Run("leading space on non-nil fails with no_extra_ws", func(t *testing.T) {
		t.Parallel()
		vs := violationsOf(t, validation.Struct(pointerPayload{Value: strPtr(" abc")}))
		require.Len(t, vs, 1)
		assert.Equal(t, "value", vs[0].Field)
		assert.Equal(t, validation.TagNoExtraWhitespace, vs[0].Rule)
	})

	t.Run("double inner space on non-nil fails with no_extra_ws", func(t *testing.T) {
		t.Parallel()
		vs := violationsOf(t, validation.Struct(pointerPayload{Value: strPtr("a  b")}))
		require.Len(t, vs, 1)
		assert.Equal(t, validation.TagNoExtraWhitespace, vs[0].Rule)
	})

	t.Run("empty non-nil pointer passes no_extra_ws but trips min", func(t *testing.T) {
		t.Parallel()
		vs := violationsOf(t, validation.Struct(pointerPayload{Value: strPtr("")}))
		require.Len(t, vs, 1)
		assert.Equal(t, "min", vs[0].Rule)
	})
}
