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

// Role subject constants. See spec §6.5.

// DepartmentResponsible — aggregate_type = "department".
const (
	SubjectDepartmentResponsibleAssigned       = "medincident.event.department.v1.responsible_assigned"
	SubjectDepartmentResponsibleRevoked        = "medincident.event.department.v1.responsible_revoked"
	SubjectDepartmentResponsibleDeputyAssigned = "medincident.event.department.v1.responsible_deputy_assigned"
	SubjectDepartmentResponsibleDeputyRemoved  = "medincident.event.department.v1.responsible_deputy_removed"
)

// ClinicHead — aggregate_type = "clinic".
const (
	SubjectClinicHeadAssigned       = "medincident.event.clinic.v1.head_assigned"
	SubjectClinicHeadRevoked        = "medincident.event.clinic.v1.head_revoked"
	SubjectClinicHeadDeputyAssigned = "medincident.event.clinic.v1.head_deputy_assigned"
	SubjectClinicHeadDeputyRemoved  = "medincident.event.clinic.v1.head_deputy_removed"
)

// OrganizationAdmin / Head / Dispatcher — aggregate_type = "organization".
const (
	SubjectOrganizationAdminAssigned       = "medincident.event.organization.v1.admin_assigned"
	SubjectOrganizationAdminRevoked        = "medincident.event.organization.v1.admin_revoked"
	SubjectOrganizationAdminDeputyAssigned = "medincident.event.organization.v1.admin_deputy_assigned"
	SubjectOrganizationAdminDeputyRemoved  = "medincident.event.organization.v1.admin_deputy_removed"

	SubjectOrganizationHeadAssigned       = "medincident.event.organization.v1.head_assigned"
	SubjectOrganizationHeadRevoked        = "medincident.event.organization.v1.head_revoked"
	SubjectOrganizationHeadDeputyAssigned = "medincident.event.organization.v1.head_deputy_assigned"
	SubjectOrganizationHeadDeputyRemoved  = "medincident.event.organization.v1.head_deputy_removed"

	SubjectOrganizationDispatcherAssigned       = "medincident.event.organization.v1.dispatcher_assigned"
	SubjectOrganizationDispatcherRevoked        = "medincident.event.organization.v1.dispatcher_revoked"
	SubjectOrganizationDispatcherDeputyAssigned = "medincident.event.organization.v1.dispatcher_deputy_assigned"
	SubjectOrganizationDispatcherDeputyRemoved  = "medincident.event.organization.v1.dispatcher_deputy_removed"
)

// SystemAdmin — aggregate_type = "system_admin".
const (
	SubjectSystemAdminGranted = "medincident.event.system_admin.v1.granted"
	SubjectSystemAdminRevoked = "medincident.event.system_admin.v1.revoked"
)

// Aggregate type strings for role events. The existing
// AggregateTypeEmployee remains valid for Employee/Vacation events.
const (
	AggregateTypeDepartment   = "department"
	AggregateTypeClinic       = "clinic"
	AggregateTypeOrganization = "organization"
	AggregateTypeSystemAdmin  = "system_admin"
)
