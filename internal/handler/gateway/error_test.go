package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	errorv1 "github.com/medincident/medincident-backend/pkg/error/v1"
)

func callHandler(t *testing.T, err error) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	GatewayErrorHandler(context.Background(), nil, nil, rec, httptest.NewRequestWithContext(context.Background(), "GET", "/", http.NoBody), err)
	return rec
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v\nbody=%s", err, rec.Body.String())
	}
	return body
}

func TestGatewayErrorHandler_UnexpectedError_Returns500(t *testing.T) {
	// No ErrorCode in details → unexpected_error.
	grpcErr := grpcstatus.New(codes.Internal, "internal error").Err()
	rec := callHandler(t, grpcErr)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
	body := decodeBody(t, rec)
	if body["code"] != "unexpected_error" {
		t.Errorf("expected code=unexpected_error, got %v", body["code"])
	}
	if body["message"] != "internal error" {
		t.Errorf("expected message=internal error, got %v", body["message"])
	}
	if _, ok := body["details"]; ok {
		t.Error("details must be absent for unexpected_error")
	}
}

func TestGatewayErrorHandler_DomainError_MapsCodeToHTTP(t *testing.T) {
	st := grpcstatus.New(codes.NotFound, "Employee not found.")
	st2, err := st.WithDetails(&errorv1.ErrorCode{Code: "employee_not_found"})
	if err != nil {
		t.Fatalf("WithDetails: %v", err)
	}
	rec := callHandler(t, st2.Err())

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
	body := decodeBody(t, rec)
	if body["code"] != "employee_not_found" {
		t.Errorf("expected code=employee_not_found, got %v", body["code"])
	}
	if body["message"] != "Employee not found." {
		t.Errorf("expected public message, got %v", body["message"])
	}
	if _, ok := body["details"]; ok {
		t.Error("details must be absent for domain error without violations")
	}
}

func TestGatewayErrorHandler_ValidationFailed_Returns400WithViolations(t *testing.T) {
	param := "256"
	st := grpcstatus.New(codes.InvalidArgument, "request is invalid")
	st2, err := st.WithDetails(
		&errorv1.ErrorCode{Code: "validation_failed"},
		&errorv1.ValidationFailedDetails{
			Violations: []*errorv1.ValidationFailedDetails_FieldViolation{
				{Field: "payload.name", Rule: "required", Message: "name is required"},
				{Field: "payload.position", Rule: "max", Message: "too long", Param: &param},
			},
		},
	)
	if err != nil {
		t.Fatalf("WithDetails: %v", err)
	}
	rec := callHandler(t, st2.Err())

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	body := decodeBody(t, rec)
	if body["code"] != "validation_failed" {
		t.Errorf("expected code=validation_failed, got %v", body["code"])
	}

	details, ok := body["details"].(map[string]any)
	if !ok {
		t.Fatalf("expected details object, got %T: %v", body["details"], body["details"])
	}
	violations, ok := details["violations"].([]any)
	if !ok || len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %v", details["violations"])
	}
	v0 := violations[0].(map[string]any)
	if v0["field"] != "payload.name" || v0["rule"] != "required" {
		t.Errorf("unexpected first violation: %v", v0)
	}
	v1 := violations[1].(map[string]any)
	if v1["param"] != "256" {
		t.Errorf("expected param=256, got %v", v1["param"])
	}
}

func TestGatewayErrorHandler_PermissionDenied_Returns403(t *testing.T) {
	st := grpcstatus.New(codes.PermissionDenied, "Permission denied.")
	st2, err := st.WithDetails(&errorv1.ErrorCode{Code: "permission_denied"})
	if err != nil {
		t.Fatalf("WithDetails: %v", err)
	}
	rec := callHandler(t, st2.Err())

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestGatewayErrorHandler_Unauthenticated_Returns401(t *testing.T) {
	st := grpcstatus.New(codes.Unauthenticated, "Invalid or expired token.")
	st2, err := st.WithDetails(&errorv1.ErrorCode{Code: "unauthenticated"})
	if err != nil {
		t.Fatalf("WithDetails: %v", err)
	}
	rec := callHandler(t, st2.Err())

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestGatewayErrorHandler_ServerFaultCode_Returns500MaskedMessage(t *testing.T) {
	// A server-fault code like employee_save_failed reaches the gateway
	// with message "internal error" (masked by ErrorInterceptor) and no
	// ErrorCode detail (server faults skip detail attachment).
	// The gateway must return 500 "unexpected_error".
	st := grpcstatus.New(codes.Internal, "internal error")
	// No details — server fault path in ErrorInterceptor adds none.
	rec := callHandler(t, st.Err())

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
	body := decodeBody(t, rec)
	if body["code"] != "unexpected_error" {
		t.Errorf("expected code=unexpected_error, got %v", body["code"])
	}
}

func TestGatewayErrorHandler_ConflictCode_Returns409(t *testing.T) {
	st := grpcstatus.New(codes.AlreadyExists, "Employee already hired.")
	st2, err := st.WithDetails(&errorv1.ErrorCode{Code: "employee_already_hired"})
	if err != nil {
		t.Fatalf("WithDetails: %v", err)
	}
	rec := callHandler(t, st2.Err())
	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rec.Code)
	}
}

func TestGatewayErrorHandler_UnprocessableCode_Returns422(t *testing.T) {
	st := grpcstatus.New(codes.FailedPrecondition, "Invalid status transition.")
	st2, err := st.WithDetails(&errorv1.ErrorCode{Code: "incident_invalid_status_transition"})
	if err != nil {
		t.Fatalf("WithDetails: %v", err)
	}
	rec := callHandler(t, st2.Err())
	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rec.Code)
	}
}

func TestGatewayErrorHandler_ContentTypeHeader(t *testing.T) {
	rec := callHandler(t, grpcstatus.New(codes.Internal, "internal error").Err())
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type: application/json, got %q", ct)
	}
}
