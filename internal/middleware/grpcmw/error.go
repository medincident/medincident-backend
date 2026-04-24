// Package middleware hosts cross-cutting gRPC interceptors for the
// command service. Each file in the package exposes one interceptor
// constructor (e.g. ErrorInterceptor) that can be chained at server
// construction time.
//
// The error interceptor converts oops-wrapped service errors into gRPC
// status errors with rich google.rpc.errdetails payloads, so that
// clients can distinguish validation failures, missing entities,
// precondition violations and internal faults without parsing
// free-form strings.
//
// Mapping rules:
//
//   - `errors.Join` results (multi-error validation) are flattened into
//     a single codes.InvalidArgument status with one
//     BadRequest.FieldViolation per leaf. Field name comes from the
//     leaf's oops context under key "field" when set, otherwise the
//     leaf's Code() string.
//   - A single oops leaf is mapped to a gRPC code via a suffix-based
//     table with explicit overrides for the handful of codes that
//     don't fit the pattern. The leaf's Public() string becomes the
//     status message for client-visible codes; for Internal/Unavailable
//     the client only sees "internal error" while the server log keeps
//     the full context.
package grpcmw

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorInterceptor returns a grpc.UnaryServerInterceptor that runs the
// handler and translates any non-nil error into a gRPC status error.
func ErrorInterceptor(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}
		return nil, translateError(logger, info.FullMethod, err)
	}
}

// translateError is the pure function half of ErrorInterceptor; kept
// unexported and tested directly.
func translateError(logger *zerolog.Logger, method string, err error) error {
	leaves := flattenErrorLeaves(err)

	if len(leaves) == 0 {
		logger.Error().Err(err).Str("grpc_method", method).Msg("handler returned non-oops error")
		return status.Error(codes.Internal, "internal error")
	}

	if len(leaves) > 1 {
		return buildBadRequestStatus(logger, method, leaves)
	}
	return buildSingleStatus(logger, method, leaves[0])
}

// flattenErrorLeaves walks err and collects oops leaves. `errors.Join`
// produces an error whose Unwrap() returns []error; we recurse into
// those branches. Everything else is probed with oops.AsOops and, if it
// matches, added as a leaf. Leaves are stored as pointers to avoid
// copying the 280-byte OopsError struct on every operation.
func flattenErrorLeaves(err error) []*oops.OopsError {
	var leaves []*oops.OopsError
	var walk func(error)
	walk = func(e error) {
		if e == nil {
			return
		}
		if joined, ok := e.(interface{ Unwrap() []error }); ok {
			for _, child := range joined.Unwrap() {
				walk(child)
			}
			return
		}
		if o, ok := oops.AsOops(e); ok {
			leaves = append(leaves, &o)
		}
	}
	walk(err)
	return leaves
}

func buildBadRequestStatus(logger *zerolog.Logger, method string, leaves []*oops.OopsError) error {
	violations := make([]*errdetails.BadRequest_FieldViolation, 0, len(leaves))
	for _, leaf := range leaves {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       errorFieldName(leaf),
			Description: errorDescription(leaf),
			Reason:      errorCodeString(leaf),
		})
	}
	st := status.New(codes.InvalidArgument, "request is invalid")
	withDetails, detailErr := st.WithDetails(&errdetails.BadRequest{FieldViolations: violations})
	if detailErr != nil {
		logger.Error().Err(detailErr).Str("grpc_method", method).Msg("failed to attach BadRequest details")
		return st.Err()
	}
	logBadRequestError(logger, method, leaves)
	return withDetails.Err()
}

func buildSingleStatus(logger *zerolog.Logger, method string, leaf *oops.OopsError) error {
	code := errorCodeString(leaf)
	grpcCode := gRPCCodeForError(code)
	// For codes that indicate server-side faults, never leak internal
	// details (error code, wrapped message) to the client. The full
	// error is still logged below with its code and domain.
	var msg string
	switch grpcCode {
	case codes.Internal, codes.Unavailable, codes.Unknown:
		msg = "internal error"
	default:
		msg = errorDescription(leaf)
	}

	st := status.New(grpcCode, msg)
	info := &errdetails.ErrorInfo{
		Reason:   code,
		Domain:   leaf.Domain(),
		Metadata: errorContextToMetadata(leaf.Context()),
	}
	withDetails, detailErr := st.WithDetails(info)
	if detailErr != nil {
		logger.Error().Err(detailErr).Str("grpc_method", method).Msg("failed to attach ErrorInfo details")
		logSingleError(logger, method, leaf, grpcCode)
		return st.Err()
	}
	logSingleError(logger, method, leaf, grpcCode)
	return withDetails.Err()
}

func logBadRequestError(logger *zerolog.Logger, method string, leaves []*oops.OopsError) {
	fields := make([]string, len(leaves))
	for i, leaf := range leaves {
		fields[i] = errorFieldName(leaf) + "=" + errorCodeString(leaf)
	}
	logger.Debug().
		Str("grpc_method", method).
		Str("grpc_code", codes.InvalidArgument.String()).
		Strs("violations", fields).
		Msg("invalid request")
}

func logSingleError(logger *zerolog.Logger, method string, leaf *oops.OopsError, grpcCode codes.Code) {
	// Build the full event in one chain so zerologlint is happy and
	// the level matches the severity of the translated gRPC code.
	switch grpcCode {
	case codes.Internal, codes.Unavailable, codes.Unknown:
		logger.Error().
			Err(leaf).
			Str("grpc_method", method).
			Str("grpc_code", grpcCode.String()).
			Str("error_code", errorCodeString(leaf)).
			Str("error_domain", leaf.Domain()).
			Msg("handler error")
	case codes.DeadlineExceeded:
		logger.Warn().
			Err(leaf).
			Str("grpc_method", method).
			Str("grpc_code", grpcCode.String()).
			Str("error_code", errorCodeString(leaf)).
			Str("error_domain", leaf.Domain()).
			Msg("handler error")
	default:
		logger.Debug().
			Err(leaf).
			Str("grpc_method", method).
			Str("grpc_code", grpcCode.String()).
			Str("error_code", errorCodeString(leaf)).
			Str("error_domain", leaf.Domain()).
			Msg("handler error")
	}
}

func errorFieldName(leaf *oops.OopsError) string {
	if f, ok := leaf.Context()["field"].(string); ok && f != "" {
		return f
	}
	return errorCodeString(leaf)
}

func errorDescription(leaf *oops.OopsError) string {
	if pub := leaf.Public(); pub != "" {
		return pub
	}
	return errorCodeString(leaf)
}

func errorCodeString(leaf *oops.OopsError) string {
	if s, ok := leaf.Code().(string); ok {
		return s
	}
	return ""
}

func errorContextToMetadata(ctx map[string]any) map[string]string {
	if len(ctx) == 0 {
		return nil
	}
	md := make(map[string]string, len(ctx))
	for k, v := range ctx {
		if s, ok := v.(string); ok {
			md[k] = s
			continue
		}
		md[k] = errorContextValueToString(v)
	}
	return md
}

// errorContextValueToString converts the scalar types we actually stash
// into oops context (IDs, counts, flags). Anything else falls through
// to the empty string to keep the metadata map clean.
func errorContextValueToString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case error:
		return x.Error()
	case interface{ String() string }:
		return x.String()
	}
	return ""
}

// gRPCCodeForError is pure and table-driven. Override wins over suffix;
// if nothing matches, we default to Internal because an unclassified
// error is never a client problem.
func gRPCCodeForError(code string) codes.Code {
	if override, ok := errorCodeOverrides[code]; ok {
		return override
	}
	for _, rule := range errorCodeSuffixes {
		if strings.HasSuffix(code, rule.suffix) {
			return rule.grpcCode
		}
	}
	return codes.Internal
}

// errorCodeOverrides lists codes that either don't match any suffix
// rule or need a different mapping than the suffix would produce.
var errorCodeOverrides = map[string]codes.Code{
	"zitadel_verify_failed":       codes.Unavailable,
	"zitadel_client_build_failed": codes.Internal,
	"read_failed":                 codes.Internal,
	"validate_failed":             codes.InvalidArgument,
	"unmarshal_failed":            codes.Internal,
	"cleanup_failed":              codes.Internal,
	"unauthenticated":             codes.Unauthenticated,
	"permission_denied":           codes.PermissionDenied,
	// Each authz.* check failure is a DB-level fault, not a client
	// error — the generic suffix rules below would route them to
	// Internal too, but the explicit overrides document the policy.
	"authz_system_admin_check_failed":      codes.Internal,
	"authz_org_access_check_failed":        codes.Internal,
	"authz_clinic_access_check_failed":     codes.Internal,
	"authz_department_access_check_failed": codes.Internal,
	"authz_employee_access_check_failed":   codes.Internal,
	"authz_vacation_access_check_failed":   codes.Internal,
	"authz_category_access_check_failed":   codes.Internal,
	"authz_type_access_check_failed":       codes.Internal,
}

// errorCodeSuffixes is checked in order; the first matching suffix
// wins. More specific suffixes must precede shorter, ambiguous ones.
var errorCodeSuffixes = []struct {
	suffix   string
	grpcCode codes.Code
}{
	// Internal infrastructure failures.
	{suffix: "_id_generation_failed", grpcCode: codes.Internal},
	{suffix: "_projection_failed", grpcCode: codes.Internal},
	{suffix: "_save_failed", grpcCode: codes.Internal},
	{suffix: "_load_failed", grpcCode: codes.Internal},
	{suffix: "_delete_failed", grpcCode: codes.Internal},
	{suffix: "_lookup_failed", grpcCode: codes.Internal},
	{suffix: "_append_failed", grpcCode: codes.Internal},
	{suffix: "_marshal_failed", grpcCode: codes.Internal},
	{suffix: "_open_failed", grpcCode: codes.Internal},
	{suffix: "_tune_failed", grpcCode: codes.Internal},

	// Client input validation (InvalidArgument). Generic codes
	// (string_required, uuid_required, float_out_of_range, …) emitted
	// by internal/service/validation/ match the _required / _too_short /
	// _too_long / _out_of_range suffixes below. Aggregate-specific
	// codes (vacation_start_required, vacation_end_before_start, …)
	// keep their own tail-matches.
	{suffix: "_required", grpcCode: codes.InvalidArgument},
	{suffix: "_too_short", grpcCode: codes.InvalidArgument},
	{suffix: "_too_long", grpcCode: codes.InvalidArgument},
	{suffix: "_out_of_range", grpcCode: codes.InvalidArgument},
	{suffix: "_end_before_start", grpcCode: codes.InvalidArgument},
	{suffix: "_end_in_past", grpcCode: codes.InvalidArgument},
	{suffix: "_start_in_past", grpcCode: codes.InvalidArgument},

	// Existence.
	{suffix: "_not_found", grpcCode: codes.NotFound},

	// Uniqueness / already-exists semantics.
	{suffix: "_already_hired", grpcCode: codes.AlreadyExists},
	{suffix: "_already_assigned", grpcCode: codes.AlreadyExists},
	{suffix: "_already_granted", grpcCode: codes.AlreadyExists},
	{suffix: "_already_started", grpcCode: codes.AlreadyExists},
	{suffix: "_already_ended", grpcCode: codes.AlreadyExists},
	{suffix: "_name_conflict", grpcCode: codes.AlreadyExists},
	{suffix: "_overlap", grpcCode: codes.AlreadyExists},

	// Business preconditions (FailedPrecondition).
	{suffix: "_not_in_department", grpcCode: codes.FailedPrecondition},
	{suffix: "_not_in_clinic", grpcCode: codes.FailedPrecondition},
	{suffix: "_not_in_organization", grpcCode: codes.FailedPrecondition},
	{suffix: "_not_in_same_organization", grpcCode: codes.FailedPrecondition},
	{suffix: "_not_started", grpcCode: codes.FailedPrecondition},
	{suffix: "_not_assigned", grpcCode: codes.FailedPrecondition},
	{suffix: "_is_holder", grpcCode: codes.FailedPrecondition},
	{suffix: "_max_depth_exceeded", grpcCode: codes.FailedPrecondition},
	{suffix: "_move_organization_mismatch", grpcCode: codes.FailedPrecondition},
	{suffix: "_move_would_create_cycle", grpcCode: codes.FailedPrecondition},
	{suffix: "_move_would_exceed_depth", grpcCode: codes.FailedPrecondition},
	{suffix: "_parent_organization_mismatch", grpcCode: codes.FailedPrecondition},
	{suffix: "_reactivate_inactive_ancestor", grpcCode: codes.FailedPrecondition},
	{suffix: "_reactivate_name_conflict", grpcCode: codes.FailedPrecondition},
	{suffix: "_parent_inactive", grpcCode: codes.FailedPrecondition},
	{suffix: "_category_inactive", grpcCode: codes.FailedPrecondition},

	// Fallback buckets kept last so the specific rules above win.
	// `_required` is declared once in the InvalidArgument group above;
	// `_empty` / `_invalid` live here as catch-alls for codes that
	// predate the unified validation vocabulary.
	{suffix: "_empty", grpcCode: codes.InvalidArgument},
	{suffix: "_invalid", grpcCode: codes.InvalidArgument},
}
