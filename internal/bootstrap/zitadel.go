package bootstrap

import (
	"context"
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
		oauth.DefaultAuthorization(cfg.KeyPath),
	)
}

// NewZitadelService builds the richer Zitadel client used by the
// command-side EmployeeService to verify that a Zitadel user exists
// before hiring. Query-server and gateway-server do not need this.
func NewZitadelService(ctx context.Context, cfg *config.ZitadelConfig, logger *zerolog.Logger) (*zitadelsvc.Service, error) {
	initCtx, cancel := context.WithTimeout(ctx, ZitadelInitTimeout)
	defer cancel()
	return zitadelsvc.NewServiceFromKeyFile(initCtx, logger, cfg.Domain, cfg.KeyPath)
}
