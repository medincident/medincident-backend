package tx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/tx"
)

type fakeTx struct{ tx.Tx }

func TestFromContextEmpty(t *testing.T) {
	_, ok := tx.FromContext(context.Background())
	require.False(t, ok)
}

func TestWithContextRoundTrip(t *testing.T) {
	f := &fakeTx{}
	ctx := tx.WithContext(context.Background(), f)
	got, ok := tx.FromContext(ctx)
	require.True(t, ok)
	require.Same(t, f, got.(*fakeTx))
}
