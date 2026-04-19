package gateway

import (
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/medincident/medincident-command-service/internal/config"
)

// Middleware is a standard net/http decorator.
type Middleware func(http.Handler) http.Handler

// AccessLog emits one structured zerolog line per request, capturing
// method, path, final status code, response size, remote address, and
// wall-clock duration in milliseconds. Status defaults to 200 when
// the inner handler writes the body without an explicit WriteHeader.
func AccessLog(logger *zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", rec.status).
				Int("bytes", rec.bytes).
				Str("remote_addr", r.RemoteAddr).
				Dur("duration_ms", time.Since(start)).
				Msg("http request")
		})
	}
}

// CORSMiddleware returns a Middleware backed by rs/cors when cfg is
// present, or nil when cfg is nil (so callers can skip the wrap).
//
// rs/cors v1.11 stores allowed-header names lowercased and matches them
// against the Access-Control-Request-Headers request values. The Fetch
// spec guarantees browsers send lowercase names there, but unit tests
// (and some gateways) may send mixed case. To keep matching symmetric,
// we lowercase AllowedHeaders at config time and normalise the
// Access-Control-Request-Headers value to lowercase before rs/cors sees
// it — exactly the same transformation the spec mandates for browsers.
func CORSMiddleware(cfg *config.GatewayCORSConfig) Middleware {
	if cfg == nil {
		return nil
	}
	lowHeaders := make([]string, len(cfg.AllowedHeaders))
	for i, h := range cfg.AllowedHeaders {
		lowHeaders[i] = strings.ToLower(h)
	}
	c := cors.New(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: cfg.AllowedMethods,
		AllowedHeaders: lowHeaders,
		MaxAge:         cfg.MaxAgeSeconds,
	})
	return func(next http.Handler) http.Handler {
		corsHandler := c.Handler(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if vals, ok := r.Header["Access-Control-Request-Headers"]; ok {
				normalized := make([]string, len(vals))
				for i, v := range vals {
					normalized[i] = strings.ToLower(v)
				}
				r = r.Clone(r.Context())
				r.Header["Access-Control-Request-Headers"] = normalized
			}
			corsHandler.ServeHTTP(w, r)
		})
	}
}

// statusRecorder shims http.ResponseWriter so AccessLog can see the
// final status code and bytes written without depending on the
// handler's cooperation.
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
	wrote  bool
}

func (s *statusRecorder) WriteHeader(code int) {
	if !s.wrote {
		s.status = code
		s.wrote = true
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if !s.wrote {
		s.wrote = true
	}
	n, err := s.ResponseWriter.Write(b)
	s.bytes += n
	return n, err
}
