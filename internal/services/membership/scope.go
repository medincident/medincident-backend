package membership

// Scope strings for oops errors. The gRPC error-mapping interceptor
// carries them as the "domain" field of ErrorInfo so clients can see
// which subsystem produced the error.
const (
	scopeEmployee              = "services.membership.employee"
	scopeVacation              = "services.membership.vacation"
	scopeSystemAdmin           = "services.membership.system_admin"
	scopeClinicHead            = "services.membership.clinic_head"
	scopeDepartmentResponsible = "services.membership.department_responsible"
	scopeOrgAdmin              = "services.membership.org_admin"
	scopeOrgDispatcher         = "services.membership.org_dispatcher"
	scopeOrgHead               = "services.membership.org_head"
)
