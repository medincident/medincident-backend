package grpcmw

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// fakeAuthorizer is a test double for authorization.AuthorizationChecker.
// It records the last token passed to CheckAuthorization so tests can
// assert the interceptor forwards the correct format.
type fakeAuthorizer struct {
	receivedToken string
	userID        string
	err           error
}

func (f *fakeAuthorizer) CheckAuthorization(_ context.Context, token string, _ ...authorization.CheckOption) (*oauth.IntrospectionContext, error) {
	f.receivedToken = token
	if f.err != nil {
		return nil, f.err
	}
	ic := &oauth.IntrospectionContext{}
	ic.IntrospectionResponse = oidc.IntrospectionResponse{
		Active:  true,
		Subject: f.userID,
	}
	return ic, nil
}

// incomingCtx builds a context with gRPC incoming metadata carrying
// the given Authorization header value. An empty value means no header.
func incomingCtx(authorizationValue string) context.Context {
	if authorizationValue == "" {
		return metadata.NewIncomingContext(context.Background(), metadata.MD{})
	}
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", authorizationValue))
}

// callInterceptor runs AuthnInterceptor chained after ErrorInterceptor
// (matching production wiring) so errors are returned as gRPC statuses.
func callInterceptor(
	ctx context.Context,
	authn grpc.UnaryServerInterceptor,
	method string,
) (context.Context, error) {
	var handlerCtx context.Context
	info := &grpc.UnaryServerInfo{FullMethod: method}
	errI := ErrorInterceptor(silentLogger())

	_, err := errI(ctx, nil, info, func(ctx context.Context, _ any) (any, error) {
		return authn(ctx, nil, info, func(ctx context.Context, _ any) (any, error) {
			handlerCtx = ctx
			return nil, nil
		})
	})
	return handlerCtx, err
}

// ─── authorizationFromMD ──────────────────────────────────────────────────────

func TestAuthorizationFromMD(t *testing.T) {
	cases := []struct {
		name  string
		ctx   context.Context // explicit context so each case controls the code path
		input string          // header value for display in error messages; "" when not applicable
		want  string          // expected return value of authorizationFromMD
	}{
		{
			name:  "no metadata",
			ctx:   context.Background(), // exercises the !ok branch in metadata.FromIncomingContext
			input: "(no metadata)",
			want:  "",
		},
		{
			name:  "empty md map",
			ctx:   incomingCtx(""), // exercises the len(values)==0 branch
			input: "",
			want:  "",
		},
		{
			name:  "whitespace-only value",
			ctx:   incomingCtx("   "),
			input: "   ",
			want:  "",
		},
		{
			name:  "wrong scheme basic",
			ctx:   incomingCtx("Basic dXNlcjpwYXNz"),
			input: "Basic dXNlcjpwYXNz",
			want:  "",
		},
		{
			name:  "bearer keyword only, no token",
			ctx:   incomingCtx("Bearer"),
			input: "Bearer",
			want:  "",
		},
		{
			name:  "bearer with space but no token",
			ctx:   incomingCtx("Bearer "),
			input: "Bearer ",
			want:  "",
		},
		{
			name:  "canonical bearer token",
			ctx:   incomingCtx("Bearer mytoken"),
			input: "Bearer mytoken",
			want:  "Bearer mytoken",
		},
		{
			name:  "lowercase bearer accepted and normalized",
			ctx:   incomingCtx("bearer mytoken"),
			input: "bearer mytoken",
			want:  "Bearer mytoken",
		},
		{
			name:  "uppercase bearer accepted and normalized",
			ctx:   incomingCtx("BEARER mytoken"),
			input: "BEARER mytoken",
			want:  "Bearer mytoken",
		},
		{
			name:  "mixed case bearer accepted and normalized",
			ctx:   incomingCtx("bEaReR mytoken"),
			input: "bEaReR mytoken",
			want:  "Bearer mytoken",
		},
		{
			name:  "leading whitespace trimmed",
			ctx:   incomingCtx("  Bearer mytoken"),
			input: "  Bearer mytoken",
			want:  "Bearer mytoken",
		},
		{
			name:  "trailing whitespace stripped from header value",
			ctx:   incomingCtx("Bearer mytoken  "),
			input: "Bearer mytoken  ",
			// TrimSpace is applied to the whole header value before prefix
			// extraction, so trailing whitespace on the token is removed.
			// The SDK receives "Bearer mytoken" and does not encounter
			// spurious whitespace in the token portion.
			want: "Bearer mytoken",
		},
		{
			name:  "double space after bearer keyword",
			ctx:   incomingCtx("Bearer  tok"),
			input: "Bearer  tok",
			// TrimSpace does not collapse internal whitespace. The prefix
			// check passes ("Bearer " matches the first 7 bytes), and the
			// remainder " tok" (with a leading space) becomes the token
			// portion, yielding "Bearer  tok". The SDK passes " tok" to
			// the introspection endpoint as-is. This documents the current
			// behaviour; a future change that normalises internal whitespace
			// would need to update this expectation.
			want: "Bearer  tok",
		},
		{
			// Use a simple opaque token value; JWT-shaped strings trigger
			// the gosec G101 hardcoded-credentials false positive.
			name:  "opaque token preserved verbatim",
			ctx:   incomingCtx("Bearer opaque-abc-123"),
			input: "Bearer opaque-abc-123",
			want:  "Bearer opaque-abc-123",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := authorizationFromMD(tc.ctx)
			if got != tc.want {
				t.Errorf("authorizationFromMD(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestAuthorizationFromMD_ReturnedValueStartsWithBearer asserts that any
// non-empty return value always begins with the canonical "Bearer " prefix.
// This is the contract required by CheckAuthorization.
func TestAuthorizationFromMD_ReturnedValueStartsWithBearer(t *testing.T) {
	variants := []string{
		"Bearer token",
		"bearer token",
		"BEARER token",
	}
	for _, v := range variants {
		got := authorizationFromMD(incomingCtx(v))
		if got == "" {
			t.Errorf("input %q: expected non-empty result", v)
			continue
		}
		if !strings.HasPrefix(got, "Bearer ") {
			t.Errorf("input %q: result %q does not start with canonical 'Bearer '", v, got)
		}
	}
}

// ─── AuthnInterceptor ────────────────────────────────────────────────────────

func TestAuthnInterceptor_SkippedMethod_BypassesAuth(t *testing.T) {
	fake := &fakeAuthorizer{err: errors.New("must not be called")}
	skip := map[string]struct{}{"/health.v1.Health/Check": {}}
	interceptor := AuthnInterceptor(fake, skip)

	// No authorization header — would fail if auth ran.
	_, err := callInterceptor(incomingCtx(""), interceptor, "/health.v1.Health/Check")
	if err != nil {
		t.Errorf("skipped method must pass through without error, got: %v", err)
	}
	if fake.receivedToken != "" {
		t.Error("CheckAuthorization must not be called for skipped methods")
	}
}

func TestAuthnInterceptor_MissingHeader_ReturnsUnauthenticated(t *testing.T) {
	fake := &fakeAuthorizer{}
	interceptor := AuthnInterceptor(fake, nil)

	_, err := callInterceptor(incomingCtx(""), interceptor, "/svc/Method")
	if err == nil {
		t.Fatal("expected error for missing header, got nil")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", st.Code())
	}
	if fake.receivedToken != "" {
		t.Error("CheckAuthorization must not be called when header is absent")
	}
}

func TestAuthnInterceptor_WrongScheme_ReturnsUnauthenticated(t *testing.T) {
	fake := &fakeAuthorizer{}
	interceptor := AuthnInterceptor(fake, nil)

	_, err := callInterceptor(incomingCtx("Basic dXNlcjpwYXNz"), interceptor, "/svc/Method")
	if err == nil {
		t.Fatal("expected error for non-Bearer scheme, got nil")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", st.Code())
	}
	if fake.receivedToken != "" {
		t.Error("CheckAuthorization must not be called for non-Bearer scheme")
	}
}

// TestAuthnInterceptor_ForwardsFullBearerHeader is the regression test for
// the double-stripping bug (issue #145). Before the fix, the interceptor
// stripped "Bearer " from the header and passed just the raw token to
// CheckAuthorization. The Zitadel SDK then did its own CutPrefix("Bearer ", ...)
// which returned ok=false, causing every valid token to be rejected.
//
// This test proves that CheckAuthorization always receives the full
// "Bearer <token>" string, not just the raw token.
func TestAuthnInterceptor_ForwardsFullBearerHeader(t *testing.T) {
	cases := []struct {
		name      string
		header    string
		wantToken string
	}{
		{
			name:      "canonical bearer",
			header:    "Bearer abc123",
			wantToken: "Bearer abc123",
		},
		{
			name:      "lowercase bearer normalized",
			header:    "bearer abc123",
			wantToken: "Bearer abc123",
		},
		{ //nolint:gosec // G101 false positive: test fixture, not a real credential
			name:      "opaque value",
			header:    "Bearer opaque-abc-xyz",
			wantToken: "Bearer opaque-abc-xyz",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeAuthorizer{err: errors.New("stop after token capture")}
			interceptor := AuthnInterceptor(fake, nil)

			_, _ = callInterceptor(incomingCtx(tc.header), interceptor, "/svc/Method")

			if fake.receivedToken != tc.wantToken {
				t.Errorf(
					"CheckAuthorization received %q, want %q — the raw token must not be stripped before forwarding",
					fake.receivedToken, tc.wantToken,
				)
			}
		})
	}
}

func TestAuthnInterceptor_InvalidToken_ReturnsUnauthenticated(t *testing.T) {
	fake := &fakeAuthorizer{err: errors.New("introspection rejected")}
	interceptor := AuthnInterceptor(fake, nil)

	_, err := callInterceptor(incomingCtx("Bearer badtoken"), interceptor, "/svc/Method")
	if err == nil {
		t.Fatal("expected error for rejected token, got nil")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", st.Code())
	}
}

func TestAuthnInterceptor_ValidToken_SetsCallerID(t *testing.T) {
	const wantUserID = "zitadel-user-abc"
	fake := &fakeAuthorizer{userID: wantUserID}
	interceptor := AuthnInterceptor(fake, nil)

	handlerCtx, err := callInterceptor(incomingCtx("Bearer validtoken"), interceptor, "/svc/Method")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	callerID, callerErr := CallerID(handlerCtx)
	if callerErr != nil {
		t.Fatalf("CallerID returned error: %v", callerErr)
	}
	if callerID != wantUserID {
		t.Errorf("CallerID = %q, want %q", callerID, wantUserID)
	}
}

func TestAuthnInterceptor_ValidToken_NonSkippedMethodRunsAuth(t *testing.T) {
	const wantUserID = "zitadel-user-xyz"
	fake := &fakeAuthorizer{userID: wantUserID}
	skip := map[string]struct{}{"/health.v1.Health/Check": {}}
	interceptor := AuthnInterceptor(fake, skip)

	// A non-skipped method must go through auth.
	handlerCtx, err := callInterceptor(incomingCtx("Bearer tok"), interceptor, "/svc/AnyOtherMethod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	callerID, _ := CallerID(handlerCtx)
	if callerID != wantUserID {
		t.Errorf("callerID = %q, want %q", callerID, wantUserID)
	}
}
