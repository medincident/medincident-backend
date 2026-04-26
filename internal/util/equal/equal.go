// Package equal provides generic equality helpers for comparable and
// Equaler types, as well as null.Value and pointer wrappers.
package equal

import "github.com/guregu/null/v6"

// Equaler constrains types that provide a value-receiver Equal(T) bool
// method. Used by EqualBy / ValueBy / PtrBy for types where == is not
// meaningful (e.g. structs with custom equality semantics like Point).
type Equaler[T any] interface {
	Equal(T) bool
}

// Equal reports whether a == b. Useful as a function value where a
// generic equality predicate is required.
func Equal[T comparable](a, b T) bool {
	return a == b
}

// EqualBy reports whether a.Equal(b). Useful as a function value where
// a generic equality predicate is required for types with custom
// equality semantics.
func EqualBy[T Equaler[T]](a, b T) bool {
	return a.Equal(b)
}

// Value reports whether two null.Value[T] wrappers hold the same state:
// both invalid, or both valid with V equal under ==. Use when T is
// comparable and == matches the intended equality.
func Value[T comparable](a, b null.Value[T]) bool {
	if a.Valid != b.Valid {
		return false
	}
	if !a.Valid {
		return true
	}
	return a.V == b.V
}

// ValueBy reports whether two null.Value[T] wrappers hold the same
// state under T's Equal method: both invalid, or both valid with
// a.V.Equal(b.V). Use when T carries custom equality semantics.
func ValueBy[T Equaler[T]](a, b null.Value[T]) bool {
	if a.Valid != b.Valid {
		return false
	}
	if !a.Valid {
		return true
	}
	return a.V.Equal(b.V)
}

// Ptr reports whether two *T pointers are equal: both nil, or both
// non-nil and dereferencing to equal T under ==.
func Ptr[T comparable](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return *a == *b
}

// PtrBy reports whether two *T pointers are equal under T's Equal
// method: both nil, or both non-nil with (*a).Equal(*b).
func PtrBy[T Equaler[T]](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return (*a).Equal(*b)
}
