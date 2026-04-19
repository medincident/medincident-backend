package model

import "github.com/guregu/null/v6"

// Equaler constrains types that provide a value-receiver Equal(T) bool
// method. Used by EqualValueBy / EqualPtrBy for types where == is not
// meaningful (e.g. structs with custom equality semantics like Point).
type Equaler[T any] interface {
	Equal(T) bool
}

// EqualValue reports whether two null.Value[T] wrappers hold the same
// state: both invalid, or both valid with V equal under ==. Use when T
// is comparable and == matches the intended equality.
func EqualValue[T comparable](a, b null.Value[T]) bool {
	if a.Valid != b.Valid {
		return false
	}
	if !a.Valid {
		return true
	}
	return a.V == b.V
}

// EqualValueBy reports whether two null.Value[T] wrappers hold the same
// state under T's Equal method: both invalid, or both valid with
// a.V.Equal(b.V). Use when T carries custom equality semantics.
func EqualValueBy[T Equaler[T]](a, b null.Value[T]) bool {
	if a.Valid != b.Valid {
		return false
	}
	if !a.Valid {
		return true
	}
	return a.V.Equal(b.V)
}

// EqualPtr reports whether two *T pointers are equal: both nil, or both
// non-nil and dereferencing to equal T under ==.
func EqualPtr[T comparable](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return *a == *b
}

// EqualPtrBy reports whether two *T pointers are equal under T's Equal
// method: both nil, or both non-nil with (*a).Equal(*b).
func EqualPtrBy[T Equaler[T]](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return (*a).Equal(*b)
}
