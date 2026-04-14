package membership

// Error codes used throughout the membership service. Grouped by
// concern for readability.
const (
	// Input validation — 422-class.
	ErrCodeEmployeeZitadelUserIDEmpty = "employee_zitadel_user_id_empty"
	ErrCodeEmployeeDepartmentIDEmpty  = "employee_department_id_empty"
	ErrCodeEmployeeIDEmpty            = "employee_id_empty"
	ErrCodeEmployeePositionTooShort   = "employee_position_too_short"
	ErrCodeEmployeePositionTooLong    = "employee_position_too_long"
	ErrCodeVacationIDEmpty            = "vacation_id_empty"
	ErrCodeVacationEndBeforeStart     = "vacation_end_before_start"
	ErrCodeVacationEndInPast          = "vacation_end_in_past"
	ErrCodeVacationStartRequired      = "vacation_start_required"
	ErrCodeVacationStartInPast        = "vacation_start_in_past"

	// Business preconditions — 409/422-class.
	ErrCodeEmployeeAlreadyHired            = "employee_already_hired"
	ErrCodeEmployeeNotFound                = "employee_not_found"
	ErrCodeDepartmentNotFound              = "department_not_found"
	ErrCodeDepartmentNotInSameOrganization = "department_not_in_same_organization"
	ErrCodeVacationNotFound                = "vacation_not_found"
	ErrCodeVacationOverlap                 = "vacation_overlap"
	ErrCodeVacationAlreadyEnded            = "vacation_already_ended"
	ErrCodeVacationNotStarted              = "vacation_not_started"
	ErrCodeVacationAlreadyStarted          = "vacation_already_started"
	ErrCodeZitadelUserNotFound             = "zitadel_user_not_found"

	// Infrastructure — 500-class.
	ErrCodeEmployeeIDGenerationFailed = "employee_id_generation_failed"
	ErrCodeEmployeeSaveFailed         = "employee_save_failed"
	ErrCodeEmployeeLoadFailed         = "employee_load_failed"
	ErrCodeEmployeeDeleteFailed       = "employee_delete_failed"
	ErrCodeEmployeeEventBuildFailed   = "employee_event_build_failed"
	ErrCodeVacationIDGenerationFailed = "vacation_id_generation_failed"
	ErrCodeVacationSaveFailed         = "vacation_save_failed"
	ErrCodeVacationLoadFailed         = "vacation_load_failed"
	ErrCodeVacationDeleteFailed       = "vacation_delete_failed"
	ErrCodeVacationEventBuildFailed   = "vacation_event_build_failed"
	ErrCodeDepartmentLookupFailed     = "department_lookup_failed"
	ErrCodeZitadelVerifyFailed        = "zitadel_verify_failed"
)

// Postgres SQLSTATE codes we match against.
const (
	pgErrCodeUniqueViolation = "23505"
)

// Additional SQLSTATE codes used by vacation commands.
const (
	pgErrCodeForeignKeyViolation = "23503"
	pgErrCodeExclusionViolation  = "23P01"
)
