package like_test

import (
	"testing"

	"github.com/medincident/medincident-backend/internal/util/like"
)

func TestEscapePattern(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain text", in: "hello", want: "hello"},
		{name: "percent", in: "100%", want: `100\%`},
		{name: "underscore", in: "a_b", want: `a\_b`},
		{name: "backslash", in: `a\b`, want: `a\\b`},
		{name: "all metacharacters", in: `%_\`, want: `\%\_\\`},
		{name: "empty", in: "", want: ""},
		{name: "mixed", in: `50% off_sale\end`, want: `50\% off\_sale\\end`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := like.EscapePattern(tt.in)
			if got != tt.want {
				t.Errorf("EscapePattern(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
