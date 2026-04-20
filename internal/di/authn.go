package di

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"

	"github.com/medincident/medincident-backend/internal/config"
	zitadelsvc "github.com/medincident/medincident-backend/internal/service/zitadel"
)

func provideAuthorizer(injector do.Injector) (*authorization.Authorizer[*oauth.IntrospectionContext], error) {
	cfg, err := do.Invoke[*config.CommandServerConfig](injector)
	if err != nil {
		return nil, err
	}

	hostname, port, tls, _, err := zitadelsvc.ParseDomain(cfg.Zitadel.Domain)
	if err != nil {
		return nil, err
	}
	opts := zitadelsvc.ZitadelOptsFromParsed(port, tls)

	ctx, cancel := context.WithTimeout(context.Background(), zitadelInitTimeout)
	defer cancel()

	return authorization.New[*oauth.IntrospectionContext](
		ctx,
		zitadel.New(hostname, opts...),
		oauth.DefaultAuthorization(cfg.Zitadel.KeyPath),
	)
}
