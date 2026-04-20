package grpcmw

import (
	"errors"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func silentLogger() *zerolog.Logger {
	l := zerolog.New(io.Discard)
	return &l
}

func TestGRPCCodeForError_Overrides(t *testing.T) {
	cases := map[string]codes.Code{
		"zitadel_verify_failed":       codes.Unavailable,
		"zitadel_client_build_failed": codes.Internal,
		"validate_failed":             codes.InvalidArgument,
		"cleanup_failed":              codes.Internal,
	}
	for code, want := range cases {
		if got := gRPCCodeForError(code); got != want {
			t.Errorf("gRPCCodeForError(%q) = %v, want %v", code, got, want)
		}
	}
}

func TestGRPCCodeForError_SuffixRules(t *testing.T) {
	cases := map[string]codes.Code{
		// Internal
		"employee_save_failed":           codes.Internal,
		"employee_load_failed":           codes.Internal,
		"organization_projection_failed": codes.Internal,
		"employee_projection_failed":     codes.Internal,
		"vacation_id_generation_failed":  codes.Internal,
		"department_lookup_failed":       codes.Internal,
		"postgres_open_failed":           codes.Internal,
		// InvalidArgument
		"employee_id_empty":              codes.InvalidArgument,
		"department_id_invalid":          codes.InvalidArgument,
		"clinic_name_empty":              codes.InvalidArgument,
		"employee_position_too_long":     codes.InvalidArgument,
		"employee_position_too_short":    codes.InvalidArgument,
		"address_latitude_out_of_range":  codes.InvalidArgument,
		"vacation_start_required":        codes.InvalidArgument,
		"vacation_end_before_start":      codes.InvalidArgument,
		"vacation_end_in_past":           codes.InvalidArgument,
		"vacation_start_in_past":         codes.InvalidArgument,
		"employee_zitadel_user_id_empty": codes.InvalidArgument,
		// NotFound
		"employee_not_found":          codes.NotFound,
		"department_not_found":        codes.NotFound,
		"zitadel_user_not_found":      codes.NotFound,
		"incident_category_not_found": codes.NotFound,
		// AlreadyExists
		"employee_already_hired":          codes.AlreadyExists,
		"clinic_head_already_assigned":    codes.AlreadyExists,
		"system_admin_already_granted":    codes.AlreadyExists,
		"vacation_already_started":        codes.AlreadyExists,
		"vacation_already_ended":          codes.AlreadyExists,
		"incident_category_name_conflict": codes.AlreadyExists,
		"vacation_overlap":                codes.AlreadyExists,
		// FailedPrecondition
		"employee_not_in_department":                     codes.FailedPrecondition,
		"department_not_in_same_organization":            codes.FailedPrecondition,
		"vacation_not_started":                           codes.FailedPrecondition,
		"deputy_not_assigned":                            codes.FailedPrecondition,
		"deputy_is_holder":                               codes.FailedPrecondition,
		"incident_category_max_depth_exceeded":           codes.FailedPrecondition,
		"incident_category_move_would_create_cycle":      codes.FailedPrecondition,
		"incident_category_move_would_exceed_depth":      codes.FailedPrecondition,
		"incident_category_reactivate_inactive_ancestor": codes.FailedPrecondition,
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

	var info *errdetails.ErrorInfo
	for _, detail := range st.Details() {
		if ei, ok := detail.(*errdetails.ErrorInfo); ok {
			info = ei
			break
		}
	}
	if info == nil {
		t.Fatal("expected ErrorInfo detail")
	}
	if info.GetReason() != "employee_not_found" {
		t.Errorf("expected reason=employee_not_found, got %q", info.GetReason())
	}
	if info.GetDomain() != "service.employee" {
		t.Errorf("expected domain=service.employee, got %q", info.GetDomain())
	}
	if info.GetMetadata()["employee_id"] != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("expected metadata to carry employee_id, got %v", info.GetMetadata())
	}
}

func TestTranslateError_MultiError_EmitsBadRequest(t *testing.T) {
	f1 := oops.In("service.employee").
		Code("employee_zitadel_user_id_empty").
		Public("Zitadel user ID is required.").
		With("field", "zitadel_user_id").
		Errorf("empty")
	f2 := oops.In("service.employee").
		Code("employee_position_too_long").
		Public("Position is too long.").
		With("field", "position").
		Errorf("too long")

	joined := errors.Join(f1, f2)
	got := translateError(silentLogger(), "/svc/Hire", joined)
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected status error, got %T", got)
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}

	var bad *errdetails.BadRequest
	for _, detail := range st.Details() {
		if br, ok := detail.(*errdetails.BadRequest); ok {
			bad = br
			break
		}
	}
	if bad == nil {
		t.Fatal("expected BadRequest detail")
	}
	if len(bad.GetFieldViolations()) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(bad.GetFieldViolations()))
	}

	fields := map[string]string{}
	for _, v := range bad.GetFieldViolations() {
		fields[v.GetField()] = v.GetReason()
	}
	if fields["zitadel_user_id"] != "employee_zitadel_user_id_empty" {
		t.Errorf("missing zitadel_user_id violation: %v", fields)
	}
	if fields["position"] != "employee_position_too_long" {
		t.Errorf("missing position violation: %v", fields)
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

	var bad *errdetails.BadRequest
	for _, detail := range st.Details() {
		if br, ok := detail.(*errdetails.BadRequest); ok {
			bad = br
			break
		}
	}
	if bad == nil || len(bad.GetFieldViolations()) != 2 {
		t.Fatalf("expected 2 violations, got %+v", bad)
	}
	if bad.GetFieldViolations()[0].GetField() != "vacation_start_required" {
		t.Errorf("expected fallback field=code, got %q", bad.GetFieldViolations()[0].GetField())
	}
}

func TestTranslateError_InternalError_MasksMessage(t *testing.T) {
	err := oops.In("service.employee").
		Code("employee_save_failed").
		Wrap(errors.New("db exploded"))

	got := translateError(silentLogger(), "/svc/Hire", err)
	st, _ := status.FromError(got)
	if st.Code() != codes.Internal {
		t.Errorf("expected Internal, got %v", st.Code())
	}
	if st.Message() != "internal error" {
		t.Errorf("expected masked message, got %q", st.Message())
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
}
