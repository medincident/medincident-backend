package di

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/services/zitadel"
)

// zitadelInitTimeout bounds the implicit OIDC discovery call performed
// by zitadel-go's client.New during service construction. Without it,
// an unreachable Zitadel instance would hang the DI bootstrap forever.
const zitadelInitTimeout = 30 * time.Second

func provideZitadelService(injector do.Injector) (*zitadel.Service, error) {
	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), zitadelInitTimeout)
	defer cancel()
	return zitadel.NewServiceFromKeyFile(ctx, logger, cfg.Zitadel.Domain, cfg.Zitadel.KeyPath)
}
