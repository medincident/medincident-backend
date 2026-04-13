package membership

// NATS subjects for Membership events.
const (
	SubjectEmployeeHired             = "medincident.event.employee.v1.hired"
	SubjectEmployeePositionChanged   = "medincident.event.employee.v1.position_changed"
	SubjectEmployeeDepartmentChanged = "medincident.event.employee.v1.department_changed"
	SubjectEmployeeTerminated        = "medincident.event.employee.v1.terminated"
	SubjectVacationStarted           = "medincident.event.employee.v1.vacation_started"
	SubjectVacationScheduled         = "medincident.event.employee.v1.vacation_scheduled"
	SubjectVacationEndDateChanged    = "medincident.event.employee.v1.vacation_end_date_changed"
	SubjectVacationEnded             = "medincident.event.employee.v1.vacation_ended"
	SubjectVacationCancelled         = "medincident.event.employee.v1.vacation_cancelled"
)

// AggregateTypeEmployee is the aggregate type string used in event
// envelopes. Both Employee and Vacation events use the same
// aggregate_type because Vacation has no independent lifecycle.
const AggregateTypeEmployee = "employee"
