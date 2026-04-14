package membership

import (
	"time"

	"github.com/guregu/null/v6"
)

// nullTimeFromPtr converts an optional time pointer into a null.Time
// ready for persistence. Shared by every vacation command that
// accepts an optional ends_at parameter.
func nullTimeFromPtr(t *time.Time) null.Time {
	if t == nil {
		return null.Time{}
	}
	return null.TimeFrom(*t)
}
