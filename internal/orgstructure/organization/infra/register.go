package organizationinfra

import "github.com/medincident/medincident-command-service/internal/outbox"

// RegisterOutboxMappers registers every Organization domain event with
// the outbox registry. Called once at startup from DI wiring.
func RegisterOutboxMappers(reg outbox.Registry) {
	for _, e := range eventInfos {
		reg.Register(e.Type, e.Info)
	}
}
