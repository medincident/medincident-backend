// Package pgerr holds the Postgres SQLSTATE strings that the command
// service matches against when classifying pgconn.PgError values.
// These are shared across domain packages so the string literals
// live in exactly one place.
package pgerr

// SQLSTATE codes — see
// https://www.postgresql.org/docs/current/errcodes-appendix.html.
const (
	CodeUniqueViolation     = "23505"
	CodeForeignKeyViolation = "23503"
	CodeExclusionViolation  = "23P01"
)
