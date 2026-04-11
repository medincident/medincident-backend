package outbox

import (
	"regexp"
	"strings"

	"google.golang.org/protobuf/proto"
)

// SubjectFromProto derives a NATS subject from a proto message's full
// name by convention:
//
//	medincident.orgstructure.v1.OrganizationCreated
//	→ medincident.orgstructure.v1.organization_created
//
// The package path is preserved; the message name is converted to
// snake_case. The publisher uses this to decide which JetStream subject
// to publish each envelope to.
func SubjectFromProto(m proto.Message) string {
	full := string(m.ProtoReflect().Descriptor().FullName())
	idx := strings.LastIndex(full, ".")
	if idx < 0 {
		return toSnake(full)
	}
	pkg := full[:idx]
	name := full[idx+1:]
	return pkg + "." + toSnake(name)
}

var camelBoundary = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func toSnake(s string) string {
	s = camelBoundary.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(s)
}
