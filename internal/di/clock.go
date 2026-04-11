package di

import (
	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/shared/clock"
)

// ProvideClock provides the process-wide clock.Clock. Production uses
// clock.System; tests can override by registering a fake clock after
// the container is built.
func ProvideClock(_ do.Injector) (clock.Clock, error) {
	return clock.System{}, nil
}
