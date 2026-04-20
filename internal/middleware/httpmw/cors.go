package httpmw

import (
	"net/http"

	"github.com/rs/cors"
)

// CORS returns a CORS middleware backed by rs/cors. Caller assembles
// rs/cors.Options from its own config shape — the middleware package
// holds no config-schema knowledge.
func CORS(opts *cors.Options) func(http.Handler) http.Handler {
	return cors.New(*opts).Handler
}
