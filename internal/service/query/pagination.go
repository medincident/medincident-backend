// Package query holds cross-cutting constants and helpers shared by all
// query-side reader packages (membership, orgstructure,
// incident/classifier, …).
package query

// Pagination bounds applied to all List endpoints.
const (
	MinLimit     = 1
	MaxLimit     = 500
	DefaultLimit = 50
)
