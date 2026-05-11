// Package projector writes the query-side denormalized read models under
// the projections.* schema. Each exported function applies one event to
// the projection inside a gorm transaction supplied by the domain consumer.
//
// Functions receive the decoded event proto plus metadata from the Envelope
// (aggregateID string, occurredAt time.Time) instead of domain model structs.
// SQL shapes are identical to the former command-side projector; only the
// data source (event proto fields vs. model struct fields) differs.
//
// Cross-projection lookups (lookupClinicID, lookupOrgName, etc.) are safe
// because NATS JetStream delivers events in strict publication order
// (seq-ordered outbox → single active publisher → ordered stream):
// parent projections always exist before child projections.
package projector

// Error codes — _failed suffix maps to codes.Internal via the gRPC
// error interceptor's existing suffix rules.
const (
	ErrCodeOrganizationProjectionFailed       = "organization_projection_failed"
	ErrCodeClinicProjectionFailed             = "clinic_projection_failed"
	ErrCodeDepartmentProjectionFailed         = "department_projection_failed"
	ErrCodeEmployeeProjectionFailed           = "employee_projection_failed"
	ErrCodeVacationProjectionFailed           = "vacation_projection_failed"
	ErrCodeRoleProjectionFailed               = "role_projection_failed"
	ErrCodeIncidentProjectionFailed           = "incident_projection_failed"
	ErrCodeIncidentClassifierProjectionFailed = "incident_classifier_projection_failed"
	ErrCodeRequestTypeProjectionFailed        = "request_type_projection_failed"
	ErrCodeServiceRequestProjectionFailed     = "service_request_projection_failed"
)

// Projectors is a value-receiver grouping of all projection functions,
// used to pass the full projector set to the domain consumer dispatcher
// without listing every function as a separate dependency.
// All methods are simple forwarders to the package-level functions.
type Projectors struct{}

// NewProjectors returns a Projectors instance. No state — projector
// functions are stateless.
func NewProjectors() *Projectors { return &Projectors{} }
