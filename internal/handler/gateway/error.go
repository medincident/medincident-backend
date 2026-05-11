package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	errorv1 "github.com/medincident/medincident-backend/pkg/error/v1"
)

type errorResponseBody struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details *detailsBody `json:"details,omitempty"`
}

type detailsBody struct {
	Violations []violationBody `json:"violations"`
}

type violationBody struct {
	Field   string  `json:"field"`
	Rule    string  `json:"rule"`
	Message string  `json:"message"`
	Param   *string `json:"param,omitempty"`
}

// GatewayErrorHandler is a runtime.ErrorHandlerFunc that replaces the
// default grpc-gateway JSON error marshaler. It extracts ErrorCode and
// ValidationFailedDetails from gRPC status details via anypb.UnmarshalTo,
// maps our domain code to an HTTP status, and writes clean JSON without
// gRPC numeric codes or @type fields.
func GatewayErrorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	// grpc-gateway always passes status errors here; a raw error produces a
	// synthetic codes.Unknown status, which maps safely to 500+unexpected_error.
	st, _ := grpcstatus.FromError(err)

	var code string
	var vfd *errorv1.ValidationFailedDetails

	for _, detail := range st.Proto().Details {
		var ec errorv1.ErrorCode
		if unmarshalErr := detail.UnmarshalTo(&ec); unmarshalErr == nil {
			code = ec.Code
			continue
		}
		var vfdMsg errorv1.ValidationFailedDetails
		if unmarshalErr := detail.UnmarshalTo(&vfdMsg); unmarshalErr == nil {
			vfd = &vfdMsg
		}
	}

	// When no ErrorCode detail is present, fall back to the transport-level
	// gRPC code. DeadlineExceeded and Canceled are sent as plain statuses by
	// the error interceptor (no domain code exists for them).
	var httpStatus int
	if code == "" {
		switch st.Code() {
		case codes.DeadlineExceeded:
			code = "deadline_exceeded"
			httpStatus = http.StatusGatewayTimeout
		case codes.Canceled:
			code = "request_canceled"
			httpStatus = http.StatusRequestTimeout
		default:
			code = "unexpected_error"
			httpStatus = http.StatusInternalServerError
		}
	} else {
		httpStatus = codeToHTTPStatus(code)
	}

	msg := st.Message()
	if httpStatus == http.StatusInternalServerError {
		msg = "internal error"
	}

	body := errorResponseBody{
		Code:    code,
		Message: msg,
	}
	if vfd != nil {
		body.Details = toDetailsBody(vfd)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(body)
}

func toDetailsBody(vfd *errorv1.ValidationFailedDetails) *detailsBody {
	vs := make([]violationBody, 0, len(vfd.Violations))
	for _, v := range vfd.Violations {
		vb := violationBody{
			Field:   v.Field,
			Rule:    v.Rule,
			Message: v.Message,
		}
		if v.Param != nil {
			p := *v.Param
			vb.Param = &p
		}
		vs = append(vs, vb)
	}
	return &detailsBody{Violations: vs}
}

// codeToHTTPStatus maps a domain error code to an HTTP status code.
// Explicit overrides take precedence over suffix rules; unknown codes
// default to 500.
func codeToHTTPStatus(code string) int {
	if code == "" {
		return http.StatusInternalServerError
	}
	if s, ok := httpCodeOverrides[code]; ok {
		return s
	}
	for _, rule := range httpCodeSuffixes {
		if strings.HasSuffix(code, rule.suffix) {
			return rule.status
		}
	}
	return http.StatusInternalServerError
}

var httpCodeOverrides = map[string]int{
	"unauthenticated":                      http.StatusUnauthorized,
	"permission_denied":                    http.StatusForbidden,
	"buffer_not_patient_owner":             http.StatusForbidden,
	"validation_failed":                    http.StatusBadRequest,
	"announcement_query_bad_cursor":        http.StatusBadRequest,
	"buffer_type_not_allowed_for_patients": http.StatusUnprocessableEntity,
	"incident_not_cancellable":             http.StatusUnprocessableEntity,
	"incident_not_reopenable":              http.StatusUnprocessableEntity,
}

var httpCodeSuffixes = []struct {
	suffix string
	status int
}{
	// 400 — client input errors
	{suffix: "_required", status: http.StatusBadRequest},
	{suffix: "_out_of_range", status: http.StatusBadRequest},
	{suffix: "_end_before_start", status: http.StatusBadRequest},
	{suffix: "_end_in_past", status: http.StatusBadRequest},
	{suffix: "_start_in_past", status: http.StatusBadRequest},
	{suffix: "_in_future", status: http.StatusBadRequest},
	{suffix: "_too_old", status: http.StatusBadRequest},
	{suffix: "_too_long", status: http.StatusBadRequest},
	{suffix: "_invalid_scope", status: http.StatusBadRequest},
	{suffix: "_invalid_time_range", status: http.StatusBadRequest},
	{suffix: "_empty", status: http.StatusBadRequest},
	{suffix: "_invalid", status: http.StatusBadRequest},
	{suffix: "_bad_cursor", status: http.StatusBadRequest},
	// 404
	{suffix: "_not_found", status: http.StatusNotFound},
	// 409 — uniqueness / already-exists
	{suffix: "_already_hired", status: http.StatusConflict},
	{suffix: "_already_assigned", status: http.StatusConflict},
	{suffix: "_already_granted", status: http.StatusConflict},
	{suffix: "_already_started", status: http.StatusConflict},
	{suffix: "_already_ended", status: http.StatusConflict},
	{suffix: "_name_conflict", status: http.StatusConflict},
	{suffix: "_overlap", status: http.StatusConflict},
	// 422 — state and business preconditions (more specific before shorter suffixes)
	{suffix: "_invalid_status_transition", status: http.StatusUnprocessableEntity},
	{suffix: "_not_in_department", status: http.StatusUnprocessableEntity},
	{suffix: "_not_in_clinic", status: http.StatusUnprocessableEntity},
	{suffix: "_not_in_organization", status: http.StatusUnprocessableEntity},
	{suffix: "_not_in_same_organization", status: http.StatusUnprocessableEntity},
	{suffix: "_not_pending", status: http.StatusUnprocessableEntity},
	{suffix: "_not_started", status: http.StatusUnprocessableEntity},
	{suffix: "_not_assigned", status: http.StatusUnprocessableEntity},
	{suffix: "_is_holder", status: http.StatusUnprocessableEntity},
	{suffix: "_archived", status: http.StatusUnprocessableEntity},
	{suffix: "_frozen", status: http.StatusUnprocessableEntity},
	{suffix: "_mismatch", status: http.StatusUnprocessableEntity},
	{suffix: "_parent_inactive", status: http.StatusUnprocessableEntity},
	{suffix: "_category_inactive", status: http.StatusUnprocessableEntity},
	{suffix: "_type_inactive", status: http.StatusUnprocessableEntity},
	{suffix: "_max_depth_exceeded", status: http.StatusUnprocessableEntity},
	{suffix: "_move_organization_mismatch", status: http.StatusUnprocessableEntity},
	{suffix: "_move_would_create_cycle", status: http.StatusUnprocessableEntity},
	{suffix: "_move_would_exceed_depth", status: http.StatusUnprocessableEntity},
	{suffix: "_parent_organization_mismatch", status: http.StatusUnprocessableEntity},
	{suffix: "_reactivate_inactive_ancestor", status: http.StatusUnprocessableEntity},
	{suffix: "_reactivate_name_conflict", status: http.StatusUnprocessableEntity},
}
