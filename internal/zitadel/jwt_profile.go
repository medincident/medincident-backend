package zitadel

import (
	"context"
	"net/url"

	"github.com/samber/oops"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	user_v2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
	zitadelcfg "github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error codes emitted by NewJWTProfileVerifier and JWTProfileVerifier.Verify.
const (
	ErrCodeZitadelClientBuildFailed = "zitadel_client_build_failed"
	ErrCodeZitadelVerifyFailed      = "zitadel_verify_failed"
)

// JWTProfileVerifier authenticates to Zitadel using a service user's
// JWT profile key and queries the v2 UserService to check existence.
type JWTProfileVerifier struct {
	cl *client.Client
}

// NewJWTProfileVerifier builds a Zitadel client authenticated via JWT
// Profile (service user key file). domain must be a URL such as
// https://auth.example.com; the host component is extracted and passed
// to the SDK. This constructor does NOT ping Zitadel — it only wires
// (convention: DI factories don't health-check).
func NewJWTProfileVerifier(ctx context.Context, domain, keyPath string) (*JWTProfileVerifier, error) {
	u, err := url.Parse(domain)
	if err != nil {
		return nil, oops.In("zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			Wrap(err)
	}

	authOpt := client.DefaultServiceUserAuthentication(
		keyPath,
		oidc.ScopeOpenID,
		client.ScopeZitadelAPI(),
	)

	cl, err := client.New(
		ctx,
		zitadelcfg.New(u.Host),
		client.WithAuth(authOpt),
	)
	if err != nil {
		return nil, oops.In("zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			Wrap(err)
	}

	return &JWTProfileVerifier{cl: cl}, nil
}

// Verify returns nil if the user exists, ErrUserNotFound if Zitadel
// reports NotFound, or a wrapped error otherwise.
func (v *JWTProfileVerifier) Verify(ctx context.Context, zitadelUserID string) error {
	_, err := v.cl.UserServiceV2().GetUserByID(ctx, &user_v2.GetUserByIDRequest{
		UserId: zitadelUserID,
	})
	if err == nil {
		return nil
	}
	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		return ErrUserNotFound
	}
	return oops.In("zitadel").
		Code(ErrCodeZitadelVerifyFailed).
		With("zitadel_user_id", zitadelUserID).
		Wrap(err)
}

// Compile-time assertion.
var _ UserVerifier = (*JWTProfileVerifier)(nil)
