package gateway

import "testing"

func TestIncomingHeaderMatcher_Authorization(t *testing.T) {
	got, ok := IncomingHeaderMatcher("Authorization")
	if !ok {
		t.Fatal("Authorization must be forwarded")
	}
	if got != "authorization" {
		t.Errorf("grpc metadata key = %q, want lowercase authorization", got)
	}
}

func TestIncomingHeaderMatcher_AuthorizationCaseInsensitive(t *testing.T) {
	for _, h := range []string{"authorization", "AUTHORIZATION", "Authorization"} {
		if _, ok := IncomingHeaderMatcher(h); !ok {
			t.Errorf("%q should be forwarded", h)
		}
	}
}

func TestIncomingHeaderMatcher_DefaultsAlsoPass(t *testing.T) {
	// grpc-gateway's default matcher forwards headers prefixed with
	// "Grpc-Metadata-" and permanent IANA HTTP headers (Accept, Cookie, etc.)
	// with the "grpcgateway-" prefix. Our matcher delegates to it for
	// anything that isn't Authorization.
	if _, ok := IncomingHeaderMatcher("Grpc-Metadata-X-Forwarded-For"); !ok {
		t.Error("Grpc-Metadata-X-Forwarded-For should be forwarded via default matcher")
	}
}

func TestIncomingHeaderMatcher_DropsUnknown(t *testing.T) {
	// Arbitrary custom headers not in the default passthrough set must
	// not be forwarded by our matcher.
	if _, ok := IncomingHeaderMatcher("X-Forwarded-For"); ok {
		t.Error("X-Forwarded-For should not be forwarded (not in default passthrough set)")
	}
}
