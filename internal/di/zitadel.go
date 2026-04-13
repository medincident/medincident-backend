package di

import (
	"context"

	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/zitadel"
)

// ProvideUserVerifier wires the real JWT Profile verifier. Only builds
// the client — no health-check.
func ProvideUserVerifier(injector do.Injector) (zitadel.UserVerifier, error) {
	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		return nil, err
	}
	return zitadel.NewJWTProfileVerifier(context.Background(), cfg.Zitadel.Domain, cfg.Zitadel.KeyPath)
}
