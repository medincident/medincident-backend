package tx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

type recordingTx struct {
	committed   int
	rolledBack  int
	closed      int
	commitErr   error
	rollbackErr error
	closeErr    error
}

func (r *recordingTx) Close() error                                    { r.closed++; return r.closeErr }
func (r *recordingTx) Commit(context.Context) error                    { r.committed++; return r.commitErr }
func (r *recordingTx) Rollback(context.Context) error                  { r.rolledBack++; return r.rollbackErr }
func (r *recordingTx) Savepoint(context.Context, string) error         { return nil }
func (r *recordingTx) RollbackSavepoint(context.Context, string) error { return nil }
func (r *recordingTx) ReleaseSavepoint(context.Context, string) error  { return nil }

func TestWithinCommitsOnSuccess(t *testing.T) {
	rt := &recordingTx{}
	err := tx.Within(context.Background(), rt, func(ctx context.Context, got tx.Tx) error {
		require.Same(t, rt, got.(*recordingTx))
		// ensure tx is stored in ctx as well
		fromCtx, ok := tx.FromContext(ctx)
		require.True(t, ok)
		require.Same(t, rt, fromCtx.(*recordingTx))
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 1, rt.committed)
	require.Equal(t, 0, rt.rolledBack)
	require.Equal(t, 1, rt.closed)
}

func TestWithinRollbacksOnError(t *testing.T) {
	rt := &recordingTx{}
	sentinel := errors.New("boom")
	err := tx.Within(context.Background(), rt, func(context.Context, tx.Tx) error {
		return sentinel
	})
	require.ErrorIs(t, err, sentinel)
	require.Equal(t, 0, rt.committed)
	require.Equal(t, 1, rt.rolledBack)
	require.Equal(t, 1, rt.closed)
}

func TestWithinRollbacksOnPanic(t *testing.T) {
	rt := &recordingTx{}
	require.Panics(t, func() {
		_ = tx.Within(context.Background(), rt, func(context.Context, tx.Tx) error {
			panic("oops")
		})
	})
	require.Equal(t, 0, rt.committed)
	require.Equal(t, 1, rt.rolledBack)
	require.Equal(t, 1, rt.closed)
}

func TestWithinWrapsCommitError(t *testing.T) {
	commitErr := errors.New("commit boom")
	rt := &recordingTx{commitErr: commitErr}

	err := tx.Within(context.Background(), rt, func(context.Context, tx.Tx) error {
		return nil
	})
	require.Error(t, err)
	require.ErrorIs(t, err, commitErr)

	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, "tx", oe.Domain())
	require.Equal(t, tx.CodeCommitFailed, oe.Code())
}

func TestWithinWrapsBothErrorsWhenRollbackFails(t *testing.T) {
	fnErr := errors.New("fn boom")
	rollbackErr := errors.New("rollback boom")
	rt := &recordingTx{rollbackErr: rollbackErr}

	err := tx.Within(context.Background(), rt, func(context.Context, tx.Tx) error {
		return fnErr
	})
	require.Error(t, err)
	// Both underlying errors are still discoverable via errors.Is because
	// errors.Join exposes Unwrap() []error.
	require.ErrorIs(t, err, fnErr)
	require.ErrorIs(t, err, rollbackErr)

	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, "tx", oe.Domain())
	require.Equal(t, tx.CodeRollbackFailed, oe.Code())
}

func TestWithinWrapsCloseErrorWhenNothingElseFailed(t *testing.T) {
	closeErr := errors.New("close boom")
	rt := &recordingTx{closeErr: closeErr}

	err := tx.Within(context.Background(), rt, func(context.Context, tx.Tx) error {
		return nil
	})
	require.Error(t, err)
	require.ErrorIs(t, err, closeErr)

	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, tx.CodeCloseFailed, oe.Code())
}

func TestWithinPrefersFnErrorOverCloseError(t *testing.T) {
	fnErr := errors.New("fn boom")
	closeErr := errors.New("close boom")
	rt := &recordingTx{closeErr: closeErr}

	err := tx.Within(context.Background(), rt, func(context.Context, tx.Tx) error {
		return fnErr
	})
	require.Error(t, err)
	// fnErr is the primary: Within's deferred Close runs (all defers do),
	// but it intentionally ignores its Close error when a real error is
	// already set, so Close's error is silently dropped.
	require.ErrorIs(t, err, fnErr)
	require.NotErrorIs(t, err, closeErr)
}
