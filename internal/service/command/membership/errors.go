package membership

// Error codes used throughout the membership service. Generic
// request-validation codes (string_required, uuid_required, …) live
// in internal/validation; this file only declares aggregate-specific
// codes (business preconditions, concurrency, infrastructure).
const (
	// Domain-specific input validation — codes that do not map to a
	// generic validation primitive because they express aggregate-level
	// constraints on time, not structural field rules.
	ErrCodeVacationEndBeforeStart = "vacation_end_before_start"
	ErrCodeVacationEndInPast      = "vacation_end_in_past"
	ErrCodeVacationEndRequired    = "vacation_end_required"
	ErrCodeVacationStartRequired  = "vacation_start_required"
	ErrCodeVacationStartInPast    = "vacation_start_in_past"

	// Business preconditions — 409/422-class.
	ErrCodeEmployeeAlreadyHired            = "employee_already_hired"
	ErrCodeEmployeeNotFound                = "employee_not_found"
	ErrCodeDepartmentNotFound              = "department_not_found"
	ErrCodeClinicNotFound                  = "clinic_not_found"
	ErrCodeOrganizationNotFound            = "organization_not_found"
	ErrCodeDepartmentNotInSameOrganization = "department_not_in_same_organization"
	ErrCodeVacationNotFound                = "vacation_not_found"
	ErrCodeVacationOverlap                 = "vacation_overlap"
	ErrCodeVacationAlreadyEnded            = "vacation_already_ended"
	ErrCodeVacationNotStarted              = "vacation_not_started"
	ErrCodeVacationAlreadyStarted          = "vacation_already_started"
	ErrCodeZitadelUserNotFound             = "zitadel_user_not_found"

	ErrCodeEmployeeNotInDepartment               = "employee_not_in_department"
	ErrCodeEmployeeNotInClinic                   = "employee_not_in_clinic"
	ErrCodeEmployeeNotInOrganization             = "employee_not_in_organization"
	ErrCodeDepartmentResponsibleAlreadyAssigned  = "department_responsible_already_assigned"
	ErrCodeDepartmentResponsibleNotFound         = "department_responsible_not_found"
	ErrCodeClinicHeadAlreadyAssigned             = "clinic_head_already_assigned"
	ErrCodeClinicHeadNotFound                    = "clinic_head_not_found"
	ErrCodeOrganizationAdminAlreadyAssigned      = "organization_admin_already_assigned"
	ErrCodeOrganizationAdminNotFound             = "organization_admin_not_found"
	ErrCodeOrganizationHeadAlreadyAssigned       = "organization_head_already_assigned"
	ErrCodeOrganizationHeadNotFound              = "organization_head_not_found"
	ErrCodeOrganizationDispatcherAlreadyAssigned = "organization_dispatcher_already_assigned"
	ErrCodeOrganizationDispatcherNotFound        = "organization_dispatcher_not_found"
	ErrCodeSystemAdminAlreadyGranted             = "system_admin_already_granted"
	ErrCodeSystemAdminNotFound                   = "system_admin_not_found"
	ErrCodeDeputyNotFound                        = "deputy_not_found"
	ErrCodeDeputyNotInDepartment                 = "deputy_not_in_department"
	ErrCodeDeputyNotInClinic                     = "deputy_not_in_clinic"
	ErrCodeDeputyNotInOrganization               = "deputy_not_in_organization"
	ErrCodeDeputyIsHolder                        = "deputy_is_holder"
	ErrCodeDeputyAlreadyAssigned                 = "deputy_already_assigned"
	ErrCodeDeputyNotAssigned                     = "deputy_not_assigned"

	// Infrastructure — 500-class.
	ErrCodeEmployeeIDGenerationFailed = "employee_id_generation_failed"
	ErrCodeEmployeeSaveFailed         = "employee_save_failed"
	ErrCodeEmployeeLoadFailed         = "employee_load_failed"
	ErrCodeEmployeeDeleteFailed       = "employee_delete_failed"
	ErrCodeEmployeeProjectionFailed   = "employee_projection_failed"
	ErrCodeVacationIDGenerationFailed = "vacation_id_generation_failed"
	ErrCodeVacationSaveFailed         = "vacation_save_failed"
	ErrCodeVacationLoadFailed         = "vacation_load_failed"
	ErrCodeVacationDeleteFailed       = "vacation_delete_failed"
	ErrCodeVacationProjectionFailed   = "vacation_projection_failed"
	ErrCodeDepartmentLookupFailed     = "department_lookup_failed"
	ErrCodeClinicLookupFailed         = "clinic_lookup_failed"
	ErrCodeOrganizationLookupFailed   = "organization_lookup_failed"
	// Zitadel verify failures are emitted via
	// services/zitadel.ErrCodeZitadelVerifyFailed — do not re-declare
	// the string here so the two constants can't drift apart.

	ErrCodeDepartmentResponsibleSaveFailed    = "department_responsible_save_failed"
	ErrCodeDepartmentResponsibleLoadFailed    = "department_responsible_load_failed"
	ErrCodeDepartmentResponsibleDeleteFailed  = "department_responsible_delete_failed"
	ErrCodeClinicHeadSaveFailed               = "clinic_head_save_failed"
	ErrCodeClinicHeadLoadFailed               = "clinic_head_load_failed"
	ErrCodeClinicHeadDeleteFailed             = "clinic_head_delete_failed"
	ErrCodeOrganizationAdminSaveFailed        = "organization_admin_save_failed"
	ErrCodeOrganizationAdminLoadFailed        = "organization_admin_load_failed"
	ErrCodeOrganizationAdminDeleteFailed      = "organization_admin_delete_failed"
	ErrCodeOrganizationHeadSaveFailed         = "organization_head_save_failed"
	ErrCodeOrganizationHeadLoadFailed         = "organization_head_load_failed"
	ErrCodeOrganizationHeadDeleteFailed       = "organization_head_delete_failed"
	ErrCodeOrganizationDispatcherSaveFailed   = "organization_dispatcher_save_failed"
	ErrCodeOrganizationDispatcherLoadFailed   = "organization_dispatcher_load_failed"
	ErrCodeOrganizationDispatcherDeleteFailed = "organization_dispatcher_delete_failed"
	ErrCodeSystemAdminSaveFailed              = "system_admin_save_failed"
	ErrCodeSystemAdminDeleteFailed            = "system_admin_delete_failed"
)
