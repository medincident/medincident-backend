package middleware

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/medincident/medincident-command-service/internal/config"
)

// HTTPCORS returns a CORS middleware backed by rs/cors when cfg is
// present, or nil when cfg is nil (so callers can skip the wrap).
func HTTPCORS(cfg *config.GatewayCORSConfig) func(http.Handler) http.Handler {
	if cfg == nil {
		return nil
	}
	return cors.New(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: cfg.AllowedMethods,
		AllowedHeaders: cfg.AllowedHeaders,
		MaxAge:         cfg.MaxAgeSeconds,
	}).Handler
}
