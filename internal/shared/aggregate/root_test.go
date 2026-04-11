package aggregate_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/shared/aggregate"
)

type stubEvent struct{ Label string }

func TestNewRootInitialisesTimestampsToNow(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	r := aggregate.NewRoot(now)
	require.Equal(t, now, r.CreatedAt)
	require.Equal(t, now, r.UpdatedAt)
	require.Empty(t, r.PullEvents())
}

func TestHydrateRootRestoresExactTimestamps(t *testing.T) {
	created := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	r := aggregate.HydrateRoot(created, updated)
	require.Equal(t, created, r.CreatedAt)
	require.Equal(t, updated, r.UpdatedAt)
	require.Empty(t, r.PullEvents())
}

func TestRaiseAppendsAndAdvancesUpdatedAt(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	later := start.Add(2 * time.Hour)

	r := aggregate.NewRoot(start)
	r.Raise(stubEvent{Label: "one"}, later)

	require.Equal(t, start, r.CreatedAt)
	require.Equal(t, later, r.UpdatedAt)

	events := r.PullEvents()
	require.Len(t, events, 1)
	require.Equal(t, stubEvent{Label: "one"}, events[0])
}

func TestPullEventsDrainsBuffer(t *testing.T) {
	r := aggregate.NewRoot(time.Now())
	r.Raise(stubEvent{Label: "a"}, time.Now())
	r.Raise(stubEvent{Label: "b"}, time.Now())

	first := r.PullEvents()
	require.Len(t, first, 2)

	second := r.PullEvents()
	require.Empty(t, second)
}
