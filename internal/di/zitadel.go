package di

import (
	"context"

	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/services/zitadel"
)

// ProvideZitadelService wires the real JWT Profile-authenticated Zitadel
// service. Only builds the client — no health-check.
func ProvideZitadelService(injector do.Injector) (*zitadel.Service, error) {
	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		return nil, err
	}
	return zitadel.NewServiceFromKeyFile(context.Background(), cfg.Zitadel.Domain, cfg.Zitadel.KeyPath)
}
