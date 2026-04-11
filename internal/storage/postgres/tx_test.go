package postgres

import (
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

func TestQuoteSavepointNameAccepts(t *testing.T) {
	cases := []string{
		"sp",
		"sp_1",
		"before_item",
		"_leading_underscore",
		"A",
		"CamelCase_42",
	}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			ident, err := quoteSavepointName(name)
			require.NoError(t, err)
			require.Equal(t, "\""+name+"\"", ident)
		})
	}
}

func TestQuoteSavepointNameRejects(t *testing.T) {
	cases := map[string]string{
		"empty":              "",
		"leading digit":      "1sp",
		"space":              "sp 1",
		"dash":               "sp-1",
		"semicolon":          "sp;DROP TABLE",
		"quote":              `sp"evil`,
		"backslash":          `sp\evil`,
		"unicode":            "спайнт",
		"too long (64 char)": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	for label, name := range cases {
		t.Run(label, func(t *testing.T) {
			_, err := quoteSavepointName(name)
			require.Error(t, err)
			oe, ok := oops.AsOops(err)
			require.True(t, ok)
			require.Equal(t, CodeSavepointInvalidName, oe.Code())
		})
	}
}

func TestUnwrapRejectsNil(t *testing.T) {
	_, err := unwrap(nil)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, CodeNilTx, oe.Code())
}

type foreignTx struct{ tx.Tx }

func TestUnwrapRejectsForeign(t *testing.T) {
	_, err := unwrap(&foreignTx{})
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, CodeForeignTx, oe.Code())
}
