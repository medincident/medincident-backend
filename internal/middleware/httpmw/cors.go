package httpmw

import (
	"net/http"

	"github.com/rs/cors"
)

// CORS returns a CORS middleware backed by rs/cors. Caller assembles
// rs/cors.Options from its own config shape — the middleware package
// holds no config-schema knowledge. opts must be non-nil; the caller
// decides whether to wrap with CORS at all (a nil-guard belongs in
// the caller so it can skip the wrap entirely when CORS is disabled).
func CORS(opts *cors.Options) func(http.Handler) http.Handler {
	return cors.New(*opts).Handler
}
