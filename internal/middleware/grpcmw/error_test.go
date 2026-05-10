package grpcmw

import (
	"errors"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/medincident/medincident-backend/internal/service/validation"
	errorv1 "github.com/medincident/medincident-backend/pkg/error/v1"
)

func silentLogger() *zerolog.Logger {
	l := zerolog.New(io.Discard)
	return &l
}

func TestGRPCCodeForError_Overrides(t *testing.T) {
	cases := map[string]codes.Code{
		"zitadel_verify_failed":                codes.Unavailable,
		"zitadel_client_build_failed":          codes.Internal,
		"validation_failed":                    codes.InvalidArgument,
		"cleanup_failed":                       codes.Internal,
		"buffer_not_patient_owner":             codes.PermissionDenied,
		"buffer_not_pending":                   codes.FailedPrecondition,
		"buffer_type_not_allowed_for_patients": codes.FailedPrecondition,
		"announcement_archived":                codes.FailedPrecondition,
		"announcement_query_bad_cursor":        codes.InvalidArgument,
		"incident_not_cancellable":             codes.FailedPrecondition,
		"incident_not_reopenable":              codes.FailedPrecondition,
		"consume_start_failed":                 codes.Internal,
		"consumer_create_failed":               codes.Internal,
		"authz_check_failed":                   codes.Internal,
	}
	for code, want := range cases {
		if got := gRPCCodeForError(code); got != want {
			t.Errorf("gRPCCodeForError(%q) = %v, want %v", code, got, want)
		}
	}
}

func TestGRPCCodeForError_SuffixRules(t *testing.T) {
	cases := map[string]codes.Code{
		"employee_save_failed":                           codes.Internal,
		"employee_load_failed":                           codes.Internal,
		"organization_projection_failed":                 codes.Internal,
		"employee_projection_failed":                     codes.Internal,
		"vacation_id_generation_failed":                  codes.Internal,
		"department_lookup_failed":                       codes.Internal,
		"postgres_open_failed":                           codes.Internal,
		"announcement_query_read_failed":                 codes.Internal,
		"buffer_query_read_failed":                       codes.Internal,
		"incident_query_read_failed":                     codes.Internal,
		"clinic_count_failed":                            codes.Internal,
		"organization_count_failed":                      codes.Internal,
		"incident_classifier_lock_failed":                codes.Internal,
		"envelope_unmarshal_malformed":                   codes.Internal,
		"payload_unmarshal_malformed":                    codes.Internal,
		"list_limit_out_of_range":                        codes.InvalidArgument,
		"list_offset_out_of_range":                       codes.InvalidArgument,
		"vacation_start_required":                        codes.InvalidArgument,
		"vacation_end_before_start":                      codes.InvalidArgument,
		"vacation_end_in_past":                           codes.InvalidArgument,
		"vacation_start_in_past":                         codes.InvalidArgument,
		"buffer_occurred_at_in_future":                   codes.InvalidArgument,
		"incident_occurred_at_in_future":                 codes.InvalidArgument,
		"buffer_occurred_at_too_old":                     codes.InvalidArgument,
		"incident_occurred_at_too_old":                   codes.InvalidArgument,
		"employee_search_query_too_long":                 codes.InvalidArgument,
		"organization_search_query_too_long":             codes.InvalidArgument,
		"announcement_invalid_scope":                     codes.InvalidArgument,
		"announcement_invalid_time_range":                codes.InvalidArgument,
		"employee_not_found":                             codes.NotFound,
		"department_not_found":                           codes.NotFound,
		"zitadel_user_not_found":                         codes.NotFound,
		"incident_category_not_found":                    codes.NotFound,
		"employee_already_hired":                         codes.AlreadyExists,
		"clinic_head_already_assigned":                   codes.AlreadyExists,
		"system_admin_already_granted":                   codes.AlreadyExists,
		"vacation_already_started":                       codes.AlreadyExists,
		"vacation_already_ended":                         codes.AlreadyExists,
		"incident_category_name_conflict":                codes.AlreadyExists,
		"vacation_overlap":                               codes.AlreadyExists,
		"employee_not_in_department":                     codes.FailedPrecondition,
		"department_not_in_same_organization":            codes.FailedPrecondition,
		"vacation_not_started":                           codes.FailedPrecondition,
		"deputy_not_assigned":                            codes.FailedPrecondition,
		"deputy_is_holder":                               codes.FailedPrecondition,
		"incident_category_max_depth_exceeded":           codes.FailedPrecondition,
		"incident_category_move_would_create_cycle":      codes.FailedPrecondition,
		"incident_category_move_would_exceed_depth":      codes.FailedPrecondition,
		"incident_category_reactivate_inactive_ancestor": codes.FailedPrecondition,
		"incident_frozen":                                codes.FailedPrecondition,
		"service_request_frozen":                         codes.FailedPrecondition,
		"incident_invalid_status_transition":             codes.FailedPrecondition,
		"service_request_invalid_status_transition":      codes.FailedPrecondition,
		"incident_type_inactive":                         codes.FailedPrecondition,
		"service_request_type_inactive":                  codes.FailedPrecondition,
		"incident_type_category_mismatch":                codes.FailedPrecondition,
		"incident_type_organization_mismatch":            codes.FailedPrecondition,
		"service_request_type_org_mismatch":              codes.FailedPrecondition,
		"service_request_incident_org_mismatch":          codes.FailedPrecondition,
		"service_request_employee_dept_mismatch":         codes.FailedPrecondition,
	}
	for code, want := range cases {
		if got := gRPCCodeForError(code); got != want {
			t.Errorf("gRPCCodeForError(%q) = %v, want %v", code, got, want)
		}
	}
}

func TestGRPCCodeForError_Unknown(t *testing.T) {
	if got := gRPCCodeForError("something_nobody_thought_of"); got != codes.Internal {
		t.Errorf("unclassified code should map to Internal, got %v", got)
	}
}

func TestTranslateError_NonOops_ReturnsInternal(t *testing.T) {
	got := translateError(silentLogger(), "/svc/Method", errors.New("bare"))
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected status error, got %T", got)
	}
	if st.Code() != codes.Internal {
		t.Errorf("expected Internal, got %v", st.Code())
	}
	if st.Message() != "internal error" {
		t.Errorf("expected generic message, got %q", st.Message())
	}
}

func TestTranslateError_SingleLeaf_MapsCode(t *testing.T) {
	err := oops.In("service.employee").
		Code("employee_not_found").
		Public("Employee not found.").
		With("employee_id", "123e4567-e89b-12d3-a456-426614174000").
		Errorf("lookup miss")

	got := translateError(silentLogger(), "/svc/Method", err)
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected status error, got %T", got)
	}
	if st.Code() != codes.NotFound {
		t.Errorf("expected NotFound, got %v", st.Code())
	}
	if st.Message() != "Employee not found." {
		t.Errorf("expected Public message, got %q", st.Message())
	}

	var code *errorv1.ErrorCode
	for _, detail := range st.Details() {
		if ec, ok := detail.(*errorv1.ErrorCode); ok {
			code = ec
			break
		}
	}
	if code == nil {
		t.Fatal("expected ErrorCode detail")
	}
	if code.Code != "employee_not_found" {
		t.Errorf("expected code=employee_not_found, got %q", code.Code)
	}

	// No ValidationFailedDetails for a simple domain error.
	for _, detail := range st.Details() {
		if _, ok := detail.(*errorv1.ValidationFailedDetails); ok {
			t.Error("domain error must not carry ValidationFailedDetails")
		}
	}
}

func TestTranslateError_ValidationFailed_EmitsErrorCodeAndDetails(t *testing.T) {
	violations := []validation.Violation{
		{Field: "zitadel_user_id", Rule: "required", Message: "required"},
		{Field: "position", Rule: "max", Param: "256", Message: "length must be at most 256 characters"},
	}
	err := oops.In("validation").
		Code(validation.CodeValidationFailed).
		With(validation.ContextKeyViolations, violations).
		Public("request is invalid").
		Errorf("validation failed: %d violation(s)", len(violations))

	got := translateError(silentLogger(), "/svc/Hire", err)
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected status error, got %T", got)
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}
	if st.Message() != "request is invalid" {
		t.Errorf("expected 'request is invalid', got %q", st.Message())
	}

	var code *errorv1.ErrorCode
	var vfd *errorv1.ValidationFailedDetails
	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errorv1.ErrorCode:
			code = d
		case *errorv1.ValidationFailedDetails:
			vfd = d
		}
	}

	if code == nil || code.Code != "validation_failed" {
		t.Fatalf("expected ErrorCode{validation_failed}, got %+v", code)
	}
	if vfd == nil {
		t.Fatal("expected ValidationFailedDetails")
	}
	if len(vfd.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(vfd.Violations))
	}

	byField := map[string]*errorv1.ValidationFailedDetails_FieldViolation{}
	for _, v := range vfd.Violations {
		byField[v.Field] = v
	}
	if got := byField["zitadel_user_id"]; got == nil || got.Rule != "required" {
		t.Errorf("zitadel_user_id violation: got %+v", got)
	}
	pos := byField["position"]
	if pos == nil || pos.Rule != "max" || pos.Message == "" {
		t.Errorf("position violation: got %+v", pos)
	}
	if pos != nil && (pos.Param == nil || *pos.Param != "256") {
		t.Errorf("position violation: expected Param=256, got %+v", pos.Param)
	}
}

func TestTranslateError_MultiError_EmitsValidationDetails(t *testing.T) {
	f1 := oops.In("service.vacation").
		Code("vacation_start_required").
		With("field", "start").
		Errorf("empty")
	f2 := oops.In("service.vacation").
		Code("vacation_end_before_start").
		With("field", "end").
		Errorf("order")

	got := translateError(silentLogger(), "/svc/Schedule", errors.Join(f1, f2))
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected status error, got %T", got)
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}

	var code *errorv1.ErrorCode
	var vfd *errorv1.ValidationFailedDetails
	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errorv1.ErrorCode:
			code = d
		case *errorv1.ValidationFailedDetails:
			vfd = d
		}
	}

	if code == nil || code.Code != "validation_failed" {
		t.Fatalf("expected ErrorCode{validation_failed}, got %+v", code)
	}
	if vfd == nil || len(vfd.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %+v", vfd)
	}

	fields := map[string]string{}
	for _, v := range vfd.Violations {
		fields[v.Field] = v.Rule
	}
	if fields["start"] != "vacation_start_required" {
		t.Errorf("missing start violation: %v", fields)
	}
	if fields["end"] != "vacation_end_before_start" {
		t.Errorf("missing end violation: %v", fields)
	}
}

func TestTranslateError_MultiError_FallsBackToCodeWhenNoField(t *testing.T) {
	f1 := oops.In("service.vacation").
		Code("vacation_start_required").
		Errorf("empty")
	f2 := oops.In("service.vacation").
		Code("vacation_end_before_start").
		Errorf("order")

	got := translateError(silentLogger(), "/svc/Schedule", errors.Join(f1, f2))
	st, _ := status.FromError(got)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}

	var vfd *errorv1.ValidationFailedDetails
	for _, detail := range st.Details() {
		if d, ok := detail.(*errorv1.ValidationFailedDetails); ok {
			vfd = d
		}
	}
	if vfd == nil || len(vfd.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %+v", vfd)
	}
	// When no "field" context key, Field falls back to the error code.
	if vfd.Violations[0].Field != "vacation_start_required" {
		t.Errorf("expected fallback field=code, got %q", vfd.Violations[0].Field)
	}
}

func TestTranslateError_InternalError_MasksMessage(t *testing.T) {
	err := oops.In("service.employee").
		Code("employee_save_failed").
		With("key_path", "/etc/secret/key.pem").
		Wrap(errors.New("db exploded"))

	got := translateError(silentLogger(), "/svc/Hire", err)
	st, _ := status.FromError(got)
	if st.Code() != codes.Internal {
		t.Errorf("expected Internal, got %v", st.Code())
	}
	if st.Message() != "internal error" {
		t.Errorf("expected masked message, got %q", st.Message())
	}

	// Server faults must not carry any ErrorCode detail.
	for _, detail := range st.Details() {
		if _, ok := detail.(*errorv1.ErrorCode); ok {
			t.Error("server-fault response must not carry ErrorCode detail")
		}
	}
}

func TestTranslateError_ZitadelVerifyFailed_Unavailable(t *testing.T) {
	err := oops.In("service.employee").
		Code("zitadel_verify_failed").
		Wrap(errors.New("upstream down"))

	got := translateError(silentLogger(), "/svc/Hire", err)
	st, _ := status.FromError(got)
	if st.Code() != codes.Unavailable {
		t.Errorf("expected Unavailable, got %v", st.Code())
	}
	if st.Message() != "internal error" {
		t.Errorf("expected masked message, got %q", st.Message())
	}
	for _, detail := range st.Details() {
		if _, ok := detail.(*errorv1.ErrorCode); ok {
			t.Error("server-fault response must not carry ErrorCode detail")
		}
	}
}

func TestTranslateError_PermissionDenied_EmitsCodeOnly(t *testing.T) {
	err := oops.In("service.authz").
		Code("permission_denied").
		Public("Permission denied.").
		With("caller_id", "zitadel-user-abc123").
		With("policy", "AdminOf.Organization(org-uuid)").
		With("organization_id", "org-uuid").
		Errorf("access denied")

	got := translateError(silentLogger(), "/svc/Method", err)
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected status error, got %T", got)
	}
	if st.Code() != codes.PermissionDenied {
		t.Errorf("expected PermissionDenied, got %v", st.Code())
	}

	var code *errorv1.ErrorCode
	for _, detail := range st.Details() {
		if ec, ok := detail.(*errorv1.ErrorCode); ok {
			code = ec
			break
		}
	}
	if code == nil {
		t.Fatal("expected ErrorCode detail for PermissionDenied")
	}
	if code.Code != "permission_denied" {
		t.Errorf("expected code=permission_denied, got %q", code.Code)
	}
	// No ValidationFailedDetails and no other details that expose metadata.
	if len(st.Details()) != 1 {
		t.Errorf("expected exactly 1 detail (ErrorCode only), got %d", len(st.Details()))
	}
}
