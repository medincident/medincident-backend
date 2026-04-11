package outbox

import (
	"reflect"
	"sync"

	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
)

// Error codes emitted by the registry's startup-time duplicate check.
const (
	ErrCodeDuplicateMapper = "duplicate_mapper"
)

// EventInfo describes how the outbox machinery handles a given domain
// event type. Subject is the NATS subject the publisher will send the
// event to; ToProto converts the in-memory domain event to a proto
// message that the publisher will wrap in google.protobuf.Any.
type EventInfo struct {
	Subject string
	ToProto func(event any) (proto.Message, error)
}

// Registry maps domain-event Go types to their EventInfo. Built once at
// startup, read-only thereafter.
type Registry interface {
	Register(goType reflect.Type, info EventInfo)
	Lookup(goType reflect.Type) (EventInfo, bool)
}

// NewRegistry returns an empty in-memory registry.
func NewRegistry() Registry {
	return &registry{byType: make(map[reflect.Type]EventInfo)}
}

type registry struct {
	mu     sync.RWMutex
	byType map[reflect.Type]EventInfo
}

// Register adds an EventInfo. Panics on duplicate Go-type registration —
// duplicates are a programmer error and must fail loudly at startup.
func (r *registry) Register(t reflect.Type, info EventInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, dup := r.byType[t]; dup {
		panic(oops.In("outbox").
			Code(ErrCodeDuplicateMapper).
			With("go_type", t.String()).
			Errorf("duplicate registration").
			Error())
	}
	r.byType[t] = info
}

// Lookup returns the EventInfo for a given Go event type.
func (r *registry) Lookup(t reflect.Type) (EventInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.byType[t]
	return info, ok
}
