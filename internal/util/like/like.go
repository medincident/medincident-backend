// Package like provides helpers for safe SQL LIKE/ILIKE pattern
// construction.
package like

import "strings"

// EscapePattern escapes the SQL LIKE metacharacters (\, %, _) in s so
// the resulting string can be safely embedded in a LIKE/ILIKE pattern
// without unintended wildcard matching.
func EscapePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
