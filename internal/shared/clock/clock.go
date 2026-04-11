// Package clock provides a pluggable source of wall-clock time. The
// Clock interface lets a single request see a single "now" and lets
// tests run deterministically with a fixed time.
package clock

import "time"

// Clock is the injection point for time.Now.
type Clock interface {
	Now() time.Time
}

// System is the production Clock implementation backed by time.Now.
type System struct{}

// Now returns the current wall-clock time.
func (System) Now() time.Time { return time.Now() }
