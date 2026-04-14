package membership

// Scope strings for oops errors. The gRPC error-mapping interceptor
// matches on these to produce the right status code.
const (
	scopeEmployee    = "services.membership.employee"
	scopeVacation    = "services.membership.vacation"
	scopeSystemAdmin = "services.membership.system_admin"
)
