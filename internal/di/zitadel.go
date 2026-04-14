package di

import (
	"context"

	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/services/zitadel"
)

func provideZitadelService(injector do.Injector) (*zitadel.Service, error) {
	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		return nil, err
	}
	return zitadel.NewServiceFromKeyFile(context.Background(), cfg.Zitadel.Domain, cfg.Zitadel.KeyPath)
}
