// Package projector writes the command-side denormalized read models
// under the projections.* schema. Each exported function applies one
// command-side change atomically with the originating domain mutation:
// callers invoke the projector inside the same db.Transaction(...)
// callback as the domain Create/Update.
//
// Functions are stateless and take the active *gorm.DB (the tx) plus
// whatever domain model(s) carry the new state. Counter bumps that
// accompany an event live inline in the per-event function — there
// are no IncrementXxx helpers exported.
//
// Errors are wrapped with oops.In("projector.<aggregate>") and one of
// the package-level Code constants below. Wrapping happens here so
// command services do not double-wrap and the gRPC ErrorInterceptor
// in internal/middleware can map projector failures the same way it
// maps service failures (the _failed suffix routes them to
// codes.Internal by the existing default).
package projector

// Error codes emitted by every function in this package. One code per
// projection family — the family name appears as a suffix so the gRPC
// suffix-based mapping in internal/middleware/error.go routes them to
// codes.Internal by the existing _failed rule without per-code overrides.
const (
	ErrCodeOrganizationProjectionFailed = "organization_projection_failed"
	ErrCodeClinicProjectionFailed       = "clinic_projection_failed"
	ErrCodeDepartmentProjectionFailed   = "department_projection_failed"
	ErrCodeEmployeeProjectionFailed     = "employee_projection_failed"
	ErrCodeVacationProjectionFailed     = "vacation_projection_failed"
	ErrCodeRoleProjectionFailed         = "role_projection_failed"
	ErrCodeIncidentProjectionFailed     = "incident_projection_failed"
)
