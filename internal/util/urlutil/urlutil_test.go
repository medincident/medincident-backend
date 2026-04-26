package urlutil_test

import (
	"testing"

	"github.com/medincident/medincident-backend/internal/util/urlutil"
)

func TestRedact(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no userinfo",
			in:   "nats://localhost:4222",
			want: "nats://localhost:4222",
		},
		{ //nolint:gosec // test fixture, not a real credential
			name: "with user and password",
			in:   "nats://admin:secret@localhost:4222",
			want: "nats://%2A%2A%2A:%2A%2A%2A@localhost:4222",
		},
		{
			name: "user only",
			in:   "nats://admin@localhost:4222",
			want: "nats://%2A%2A%2A:%2A%2A%2A@localhost:4222",
		},
		{
			name: "empty string",
			in:   "",
			want: "",
		},
		{
			name: "unparseable URL",
			in:   "://bad\x7f",
			want: "***",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlutil.Redact(tt.in)
			if got != tt.want {
				t.Errorf("Redact(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
