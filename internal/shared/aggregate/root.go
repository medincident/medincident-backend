package aggregate

import "time"

// Root is the embedded base of every aggregate root. It owns the
// lifecycle timestamps and the buffered domain events.
//
// CreatedAt/UpdatedAt are exported. Writing to them from outside is a
// contract violation, not a system bug — Go cannot fully enforce private
// fields and we accept the idiom. events stays unexported: its only
// legitimate access protocol is Raise/PullEvents.
type Root struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	events []any
}

// NewRoot creates a fresh Root with CreatedAt == UpdatedAt == now.
// UpdatedAt is never zero.
func NewRoot(now time.Time) Root {
	return Root{CreatedAt: now, UpdatedAt: now}
}

// HydrateRoot rebuilds a Root from persisted timestamps. It is the
// restoration-path counterpart of NewRoot and is called exclusively
// from repository loaders (postgres, tests) after reading a row. It
// does NOT validate the timestamps (the repository is trusted) and
// does NOT raise any events — the aggregate is brought back into
// memory with an empty event buffer, so PullEvents immediately after
// a Hydrate returns nil.
func HydrateRoot(createdAt, updatedAt time.Time) Root {
	return Root{CreatedAt: createdAt, UpdatedAt: updatedAt}
}

// Raise registers a domain event and atomically advances UpdatedAt.
// This is the ONLY legitimate way to move UpdatedAt.
func (r *Root) Raise(event any, now time.Time) {
	r.events = append(r.events, event)
	r.UpdatedAt = now
}

// PullEvents drains the buffered events and returns them.
// Satisfies outbox.EventSource.
func (r *Root) PullEvents() []any {
	out := r.events
	r.events = nil
	return out
}
