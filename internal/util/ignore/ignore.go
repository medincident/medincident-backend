// Package ignore discards the second return value of common two-value
// Go functions so they can be called inline where only the value is
// wanted.
//
// These helpers should be used sparingly — they exist for cases where
// the discarded value is deliberately handled by a downstream check
// (e.g. uuid.Parse into a command field whose Validate rejects
// uuid.Nil), not as a shortcut for ignoring errors that should be
// handled.
package ignore

// Any returns value and discards the second return value regardless
// of its type. Error and Exists delegate here so every "drop the
// second return value" site funnels through a single primitive.
func Any[T any, U any](value T, _ U) T { return value }

// Error returns value and discards the error. Alias of Any with the
// second argument constrained to error so the call site documents
// intent.
func Error[T any](value T, err error) T { return Any(value, err) }

// Exists returns value and discards the presence flag. Alias of Any
// for the two-value map-lookup / type-assertion shape.
func Exists[T any](value T, ok bool) T { return Any(value, ok) }
