package httpmw

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// AccessLog emits one structured zerolog line per HTTP request.
// It delegates status and byte accounting to rs/zerolog/hlog, which
// wraps the ResponseWriter and captures both transparently.
func AccessLog(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", status).
			Int("bytes", size).
			Str("remote_addr", r.RemoteAddr).
			Int64("duration_ms", duration.Milliseconds()).
			Msg("http request")
	})
}
