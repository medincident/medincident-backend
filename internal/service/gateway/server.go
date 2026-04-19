package gateway

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/oops"
	"google.golang.org/grpc"

	cmdclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
	cmdmembershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
	cmdorgv1 "github.com/medincident/medincident-command-service/pkg/command/orgstructure/v1"
	qidentityv1 "github.com/medincident/medincident-command-service/pkg/query/identity/v1"
	qclassifierv1 "github.com/medincident/medincident-command-service/pkg/query/incident/classifier/v1"
	qmembershipv1 "github.com/medincident/medincident-command-service/pkg/query/membership/v1"
	qorgv1 "github.com/medincident/medincident-command-service/pkg/query/orgstructure/v1"
	qstatsv1 "github.com/medincident/medincident-command-service/pkg/query/stats/v1"
)

// Error codes emitted by this file.
const (
	ErrCodeGatewayRegisterFailed = "gateway_register_failed"
)

// ServerWrapper owns the *http.Server lifecycle. Implements the
// samber/do Shutdowner protocol.
type ServerWrapper struct {
	Server *http.Server
}

// Shutdown forwards to http.Server.Shutdown bounded by ctx.
func (w *ServerWrapper) Shutdown(ctx context.Context) error {
	return w.Server.Shutdown(ctx)
}

// NewServeMux registers every command- and query-side grpc-gateway
// handler against the supplied ClientConns, wires the
// IncomingHeaderMatcher, and returns the mux. Caller composes it
// behind the health/middleware chain via BuildHandler.
func NewServeMux(ctx context.Context, command, query *grpc.ClientConn) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(IncomingHeaderMatcher),
	)

	errs := []error{
		// Command side.
		cmdorgv1.RegisterOrgStructureCommandServiceHandler(ctx, mux, command),
		cmdmembershipv1.RegisterMembershipCommandServiceHandler(ctx, mux, command),
		cmdclassifierv1.RegisterIncidentClassifierCommandServiceHandler(ctx, mux, command),

		// Query side.
		qorgv1.RegisterOrgStructureQueryServiceHandler(ctx, mux, query),
		qmembershipv1.RegisterMembershipQueryServiceHandler(ctx, mux, query),
		qclassifierv1.RegisterIncidentClassifierQueryServiceHandler(ctx, mux, query),
		qstatsv1.RegisterStatsQueryServiceHandler(ctx, mux, query),
		qidentityv1.RegisterIdentityQueryServiceHandler(ctx, mux, query),
	}
	if err := errors.Join(errs...); err != nil {
		return nil, oops.In("gateway").Code(ErrCodeGatewayRegisterFailed).Wrap(err)
	}
	return mux, nil
}

// BuildHandler composes the final http.Handler: the outer chain is
// CORS → access log; inside that, health endpoints are routed off
// /healthz and /readyz, and everything else falls through to the
// grpc-gateway mux. CORS is only applied when corsMW is non-nil.
func BuildHandler(
	gatewayMux *runtime.ServeMux,
	liveness http.Handler,
	readiness http.Handler,
	accessLogMW Middleware,
	corsMW Middleware,
) http.Handler {
	router := http.NewServeMux()
	router.Handle("/healthz", liveness)
	router.Handle("/readyz", readiness)
	router.Handle("/", gatewayMux)

	var h http.Handler = router
	if accessLogMW != nil {
		h = accessLogMW(h)
	}
	if corsMW != nil {
		h = corsMW(h)
	}
	return h
}
