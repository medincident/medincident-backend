// Package middleware hosts cross-cutting gRPC interceptors for the
// command service.
//
// The error interceptor converts oops-wrapped service errors into gRPC
// status errors with custom errorv1 detail payloads:
//
//   - ErrorCode{code} is always attached for client-visible errors.
//   - ValidationFailedDetails is attached when code == "validation_failed"
//     (struct-tag validation or errors.Join of domain leaves).
//   - Server-fault codes (Internal, Unavailable, Unknown) get no details;
//     the message is always "internal error".
package grpcmw

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/medincident/medincident-backend/internal/service/validation"
	errorv1 "github.com/medincident/medincident-backend/pkg/error/v1"
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

func translateError(logger *zerolog.Logger, method string, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		logger.Warn().Err(err).Str("grpc_method", method).Str("grpc_code", codes.DeadlineExceeded.String()).Msg("handler error")
		return status.Error(codes.DeadlineExceeded, "deadline exceeded")
	}
	if errors.Is(err, context.Canceled) {
		logger.Debug().Err(err).Str("grpc_method", method).Str("grpc_code", codes.Canceled.String()).Msg("handler error")
		return status.Error(codes.Canceled, "request canceled")
	}

	leaves := flattenErrorLeaves(err)

	if len(leaves) == 0 {
		logger.Error().Err(err).Str("grpc_method", method).Msg("handler returned non-oops error")
		return status.Error(codes.Internal, "internal error")
	}

	if len(leaves) == 1 && errorCodeString(leaves[0]) == validation.CodeValidationFailed {
		if vs, ok := leaves[0].Context()[validation.ContextKeyViolations].([]validation.Violation); ok && len(vs) > 0 {
			return buildValidationStatus(logger, method, leaves[0], vs)
		}
	}

	if len(leaves) > 1 {
		return buildMultiLeafStatus(logger, method, leaves)
	}
	return buildSingleStatus(logger, method, leaves[0])
}

// flattenErrorLeaves walks err and collects oops leaves. errors.Join produces
// an error whose Unwrap() returns []error; we recurse into those branches.
// Everything else is probed with oops.AsOops and, if it matches, added as a
// leaf. Leaves are stored as pointers to avoid copying the 280-byte OopsError
// struct on every operation.
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

// buildValidationStatus emits ErrorCode{validation_failed} + ValidationFailedDetails
// from the struct-tag violations carried on the single oops leaf.
func buildValidationStatus(logger *zerolog.Logger, method string, leaf *oops.OopsError, violations []validation.Violation) error {
	fvs := make([]*errorv1.ValidationFailedDetails_FieldViolation, 0, len(violations))
	for _, v := range violations {
		fv := &errorv1.ValidationFailedDetails_FieldViolation{
			Field:   v.Field,
			Rule:    v.Rule,
			Message: v.Message,
		}
		if v.Param != "" {
			p := v.Param
			fv.Param = &p
		}
		fvs = append(fvs, fv)
	}
	ec := &errorv1.ErrorCode{Code: validation.CodeValidationFailed}
	vfd := &errorv1.ValidationFailedDetails{Violations: fvs}
	st := status.New(codes.InvalidArgument, errorDescription(leaf))
	withDetails, detailErr := st.WithDetails(ec, vfd)
	if detailErr != nil {
		logger.Error().Err(detailErr).Str("grpc_method", method).Msg("failed to attach validation details")
		return st.Err()
	}
	logValidationStatus(logger, method, leaf, violations)
	return withDetails.Err()
}

// buildMultiLeafStatus handles errors.Join leaves: emits ErrorCode{validation_failed}
// + ValidationFailedDetails with one FieldViolation per leaf, using each leaf's
// field name and oops code as the violation coordinates.
func buildMultiLeafStatus(logger *zerolog.Logger, method string, leaves []*oops.OopsError) error {
	fvs := make([]*errorv1.ValidationFailedDetails_FieldViolation, 0, len(leaves))
	for _, leaf := range leaves {
		fvs = append(fvs, &errorv1.ValidationFailedDetails_FieldViolation{
			Field:   errorFieldName(leaf),
			Rule:    errorCodeString(leaf),
			Message: errorDescription(leaf),
		})
	}
	ec := &errorv1.ErrorCode{Code: validation.CodeValidationFailed}
	vfd := &errorv1.ValidationFailedDetails{Violations: fvs}
	st := status.New(codes.InvalidArgument, "request is invalid")
	withDetails, detailErr := st.WithDetails(ec, vfd)
	if detailErr != nil {
		logger.Error().Err(detailErr).Str("grpc_method", method).Msg("failed to attach multi-leaf validation details")
		return st.Err()
	}
	logMultiLeafError(logger, method, leaves)
	return withDetails.Err()
}

// buildSingleStatus maps a single oops leaf to a gRPC status. Server-fault codes
// (Internal, Unavailable, Unknown) get no details and mask the message as
// "internal error" so implementation details never leak to the client.
func buildSingleStatus(logger *zerolog.Logger, method string, leaf *oops.OopsError) error {
	code := errorCodeString(leaf)
	grpcCode := gRPCCodeForError(code)

	logSingleError(logger, method, leaf, grpcCode)

	switch grpcCode { //nolint:exhaustive // only server-fault codes need special handling
	case codes.Internal, codes.Unavailable, codes.Unknown:
		return status.New(grpcCode, "internal error").Err()
	}

	msg := errorDescription(leaf)
	st := status.New(grpcCode, msg)
	withDetails, detailErr := st.WithDetails(&errorv1.ErrorCode{Code: code})
	if detailErr != nil {
		logger.Error().Err(detailErr).Str("grpc_method", method).Msg("failed to attach ErrorCode detail")
		return st.Err()
	}
	return withDetails.Err()
}

func logValidationStatus(logger *zerolog.Logger, method string, leaf *oops.OopsError, violations []validation.Violation) {
	fields := make([]string, len(violations))
	for i, v := range violations {
		fields[i] = v.Field + "=" + v.Rule
	}
	logger.Debug().
		Str("grpc_method", method).
		Str("grpc_code", codes.InvalidArgument.String()).
		Str("error_code", errorCodeString(leaf)).
		Str("error_domain", leaf.Domain()).
		Strs("violations", fields).
		Msg("invalid request")
}

func logMultiLeafError(logger *zerolog.Logger, method string, leaves []*oops.OopsError) {
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

// gRPCCodeForError is pure and table-driven. Override wins over suffix; if
// nothing matches, we default to Internal because an unclassified error is
// never a client problem.
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

var errorCodeOverrides = map[string]codes.Code{
	"zitadel_verify_failed":                codes.Unavailable,
	"zitadel_client_build_failed":          codes.Internal,
	"read_failed":                          codes.Internal,
	"validation_failed":                    codes.InvalidArgument,
	"unmarshal_failed":                     codes.Internal,
	"cleanup_failed":                       codes.Internal,
	"unauthenticated":                      codes.Unauthenticated,
	"permission_denied":                    codes.PermissionDenied,
	"authz_system_admin_check_failed":      codes.Internal,
	"authz_org_access_check_failed":        codes.Internal,
	"authz_clinic_access_check_failed":     codes.Internal,
	"authz_department_access_check_failed": codes.Internal,
	"authz_employee_access_check_failed":   codes.Internal,
	"authz_vacation_access_check_failed":   codes.Internal,
	"authz_category_access_check_failed":   codes.Internal,
	"authz_type_access_check_failed":       codes.Internal,
	"authz_check_failed":                   codes.Internal,
	"consume_start_failed":                 codes.Internal,
	"consumer_create_failed":               codes.Internal,
	"buffer_not_patient_owner":             codes.PermissionDenied,
	"buffer_not_pending":                   codes.FailedPrecondition,
	"buffer_type_not_allowed_for_patients": codes.FailedPrecondition,
	"announcement_archived":                codes.FailedPrecondition,
	"announcement_query_bad_cursor":        codes.InvalidArgument,
	"incident_not_cancellable":             codes.FailedPrecondition,
	"incident_not_reopenable":              codes.FailedPrecondition,
}

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
	{suffix: "_read_failed", grpcCode: codes.Internal},
	{suffix: "_count_failed", grpcCode: codes.Internal},
	{suffix: "_lock_failed", grpcCode: codes.Internal},
	{suffix: "_malformed", grpcCode: codes.Internal},
	// Client input validation (InvalidArgument).
	{suffix: "_required", grpcCode: codes.InvalidArgument},
	{suffix: "_out_of_range", grpcCode: codes.InvalidArgument},
	{suffix: "_end_before_start", grpcCode: codes.InvalidArgument},
	{suffix: "_end_in_past", grpcCode: codes.InvalidArgument},
	{suffix: "_start_in_past", grpcCode: codes.InvalidArgument},
	{suffix: "_in_future", grpcCode: codes.InvalidArgument},
	{suffix: "_too_old", grpcCode: codes.InvalidArgument},
	{suffix: "_too_long", grpcCode: codes.InvalidArgument},
	{suffix: "_invalid_scope", grpcCode: codes.InvalidArgument},
	{suffix: "_invalid_time_range", grpcCode: codes.InvalidArgument},
	// Existence.
	{suffix: "_not_found", grpcCode: codes.NotFound},
	// Duplicates (AlreadyExists).
	{suffix: "_already_hired", grpcCode: codes.AlreadyExists},
	{suffix: "_already_assigned", grpcCode: codes.AlreadyExists},
	{suffix: "_already_granted", grpcCode: codes.AlreadyExists},
	{suffix: "_already_started", grpcCode: codes.AlreadyExists},
	{suffix: "_already_ended", grpcCode: codes.AlreadyExists},
	{suffix: "_name_conflict", grpcCode: codes.AlreadyExists},
	{suffix: "_overlap", grpcCode: codes.AlreadyExists},
	// Business-rule preconditions (FailedPrecondition).
	{suffix: "_invalid_status_transition", grpcCode: codes.FailedPrecondition},
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
	{suffix: "_type_inactive", grpcCode: codes.FailedPrecondition},
	{suffix: "_frozen", grpcCode: codes.FailedPrecondition},
	{suffix: "_mismatch", grpcCode: codes.FailedPrecondition},
	// Catch-all client input (InvalidArgument) — kept last so more specific suffixes win.
	{suffix: "_empty", grpcCode: codes.InvalidArgument},
	{suffix: "_invalid", grpcCode: codes.InvalidArgument},
}
