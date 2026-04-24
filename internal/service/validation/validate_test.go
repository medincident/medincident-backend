package validation_test

import (
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
			require.Error(t, err)
			leaves := flattenJoin(err)
			require.Len(t, leaves, 1)
			o, ok := oops.AsOops(leaves[0])
			require.True(t, ok)
			assert.Equal(t, validation.CodeStringExtraWhitespace, o.Code())
			assert.Equal(t, "value", o.Context()["field"])
		})
	}
}

func flattenJoin(err error) []error {
	if u, ok := err.(interface{ Unwrap() []error }); ok {
		return u.Unwrap()
	}
	return []error{err}
}

type orderPayload struct {
	Value string `validate:"required,no_extra_ws,min=4"`
}

func TestStruct_OrderingWithOtherRules(t *testing.T) {
	t.Parallel()

	err := validation.Struct(orderPayload{Value: ""})
	require.Error(t, err)
	leaves := flattenJoin(err)
	require.Len(t, leaves, 1)
	o, _ := oops.AsOops(leaves[0])
	assert.Equal(t, validation.CodeStringRequired, o.Code())

	err = validation.Struct(orderPayload{Value: " abc"})
	require.Error(t, err)
	leaves = flattenJoin(err)
	require.Len(t, leaves, 1)
	o, _ = oops.AsOops(leaves[0])
	assert.Equal(t, validation.CodeStringExtraWhitespace, o.Code())

	err = validation.Struct(orderPayload{Value: "abc"})
	require.Error(t, err)
	leaves = flattenJoin(err)
	require.Len(t, leaves, 1)
	o, _ = oops.AsOops(leaves[0])
	assert.Equal(t, validation.CodeStringTooShort, o.Code())

	err = validation.Struct(orderPayload{Value: "alpha beta"})
	require.NoError(t, err)
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

	t.Run("leading space on non-nil fails with extra_whitespace", func(t *testing.T) {
		t.Parallel()
		err := validation.Struct(pointerPayload{Value: strPtr(" abc")})
		require.Error(t, err)
		leaves := flattenJoin(err)
		require.Len(t, leaves, 1)
		o, ok := oops.AsOops(leaves[0])
		require.True(t, ok)
		assert.Equal(t, validation.CodeStringExtraWhitespace, o.Code())
		assert.Equal(t, "value", o.Context()["field"])
	})

	t.Run("double inner space on non-nil fails with extra_whitespace", func(t *testing.T) {
		t.Parallel()
		err := validation.Struct(pointerPayload{Value: strPtr("a  b")})
		require.Error(t, err)
		leaves := flattenJoin(err)
		require.Len(t, leaves, 1)
		o, _ := oops.AsOops(leaves[0])
		assert.Equal(t, validation.CodeStringExtraWhitespace, o.Code())
	})

	t.Run("empty non-nil pointer passes no_extra_ws but trips min", func(t *testing.T) {
		t.Parallel()
		err := validation.Struct(pointerPayload{Value: strPtr("")})
		require.Error(t, err)
		leaves := flattenJoin(err)
		require.Len(t, leaves, 1)
		o, _ := oops.AsOops(leaves[0])
		assert.Equal(t, validation.CodeStringTooShort, o.Code())
	})
}
