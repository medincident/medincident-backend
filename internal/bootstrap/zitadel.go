package bootstrap

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"

	"github.com/medincident/medincident-backend/internal/config"
	zitadelsvc "github.com/medincident/medincident-backend/internal/service/zitadel"
)

// ZitadelInitTimeout bounds the implicit OIDC discovery call performed
// by zitadel-go's client.New during startup. Without it, an unreachable
// Zitadel instance would hang the boot sequence forever.
const ZitadelInitTimeout = 30 * time.Second

// NewZitadelAuthorizer builds the JWT-introspection authorizer used by
// grpcmw.AuthnInterceptor. It is shared by command-server and
// query-server; gateway-server does not validate tokens.
func NewZitadelAuthorizer(ctx context.Context, cfg *config.ZitadelConfig) (*authorization.Authorizer[*oauth.IntrospectionContext], error) {
	hostname, port, tls, _, err := zitadelsvc.ParseDomain(cfg.Domain)
	if err != nil {
		return nil, err
	}
	opts := zitadelsvc.ZitadelOptsFromParsed(port, tls)

	initCtx, cancel := context.WithTimeout(ctx, ZitadelInitTimeout)
	defer cancel()

	return authorization.New[*oauth.IntrospectionContext](
		initCtx,
		zitadel.New(hostname, opts...),
		oauth.DefaultAuthorization(cfg.IntrospectionKeyPath),
	)
}

// NewZitadelService builds the richer Zitadel client used by the
// command-side EmployeeService to verify that a Zitadel user exists
// before hiring. Query-server and gateway-server do not need this.
//
// management_key_path is optional in ZitadelConfig (query-server omits it),
// so we validate it explicitly here to surface a clear error at startup
// rather than an opaque "open : no such file" from the SDK.
func NewZitadelService(ctx context.Context, cfg *config.ZitadelConfig, logger *zerolog.Logger) (*zitadelsvc.Service, error) {
	if cfg.ManagementKeyPath == "" {
		return nil, errors.New("zitadel management_key_path is required for command-server but is not set")
	}
	initCtx, cancel := context.WithTimeout(ctx, ZitadelInitTimeout)
	defer cancel()
	return zitadelsvc.NewServiceFromKeyFile(initCtx, logger, cfg.Domain, cfg.ManagementKeyPath)
}
