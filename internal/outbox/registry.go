package outbox

import (
	"reflect"
	"sync"

	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
)

// Error codes emitted by the registry's startup-time duplicate checks.
const (
	CodeDuplicateMapper   = "outbox_duplicate_mapper"
	CodeDuplicateTypeName = "outbox_duplicate_type_name"
)

// EventInfo describes how the outbox machinery handles a given domain
// event type.
//
// TypeName is the stable identifier stored in outbox.events.event_type.
// It survives Go type renames because feature infra owns this string
// explicitly — changing a Go type name does not change the TypeName
// unless the infra package is updated deliberately.
//
// Zero returns a pointer to a fresh zero-valued domain event struct ready
// for json.Unmarshal. Used by the publisher when re-materialising events
// from the outbox row's JSONB payload.
//
// ToProto converts the materialised domain event into a proto.Message.
// Called by the publisher at publish time. This is the ONLY place proto
// enters the outbox subsystem at runtime; everywhere else the events are
// plain Go structs.
type EventInfo struct {
	TypeName string
	Zero     func() any
	ToProto  func(event any) (proto.Message, error)
}

// Registry is the registry of event mappers, keyed by both reflect.Type
// (used by Publish at write time) and TypeName (used by the publisher at
// read time).
type Registry interface {
	Register(goType reflect.Type, info EventInfo)
	ByGoType(goType reflect.Type) (EventInfo, bool)
	ByTypeName(typeName string) (EventInfo, bool)
}

// NewRegistry returns an empty in-memory registry. Intended to be built
// once at startup (DI providers register each feature's mappers) and
// thereafter read-only.
func NewRegistry() Registry {
	return &registry{
		byType: make(map[reflect.Type]EventInfo),
		byName: make(map[string]EventInfo),
	}
}

type registry struct {
	mu     sync.RWMutex
	byType map[reflect.Type]EventInfo
	byName map[string]EventInfo
}

// Register adds an EventInfo to the registry under both its reflect.Type
// and its TypeName. Panics on duplicate registration — duplicates are a
// programmer error (two infra packages trying to own the same event) and
// must fail loudly at startup rather than quietly overwriting.
func (r *registry) Register(t reflect.Type, info EventInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, dup := r.byType[t]; dup {
		panic(oops.In("outbox").
			Code(CodeDuplicateMapper).
			With("go_type", t.String()).
			Errorf("duplicate registration").
			Error())
	}
	if _, dup := r.byName[info.TypeName]; dup {
		panic(oops.In("outbox").
			Code(CodeDuplicateTypeName).
			With("type_name", info.TypeName).
			Errorf("duplicate TypeName").
			Error())
	}
	r.byType[t] = info
	r.byName[info.TypeName] = info
}

// ByGoType looks up an EventInfo by the reflect.Type of a live domain
// event value. Used at write time when we have the struct in hand.
func (r *registry) ByGoType(t reflect.Type) (EventInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.byType[t]
	return info, ok
}

// ByTypeName looks up an EventInfo by the stable string identifier
// stored in the outbox row. Used at read time (publisher) when we
// have only the column value and need to reconstruct the Go type.
func (r *registry) ByTypeName(name string) (EventInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.byName[name]
	return info, ok
}
