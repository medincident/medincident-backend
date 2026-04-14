package membership

// Error codes emitted by the membership gRPC handler transport layer.
// These cover request-parsing failures (invalid UUIDs, malformed
// fields) before the service layer takes over.
const (
	ErrCodeHandlerInvalidEmployeeID       = "employee_id_invalid"
	ErrCodeHandlerInvalidDepartmentID     = "department_id_invalid"
	ErrCodeHandlerInvalidClinicID         = "clinic_id_invalid"
	ErrCodeHandlerInvalidVacationID       = "vacation_id_invalid"
	ErrCodeHandlerInvalidDeputyEmployeeID = "deputy_employee_id_invalid"
)
