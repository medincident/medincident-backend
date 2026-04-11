package tx

import (
	"context"
	"errors"

	"github.com/samber/oops"
)

// Error codes emitted by Within.
const (
	CodeCommitFailed   = "tx_commit_failed"
	CodeRollbackFailed = "tx_rollback_failed"
	CodeCloseFailed    = "tx_close_failed"
)

// Within runs fn inside the provided transaction t. It stores t in ctx,
// commits on success, rolls back on error or panic, and always closes t.
// fn receives both the enriched ctx and the transaction t explicitly so
// callers can use either access style.
//
// Error semantics (all paths use samber/oops with codes from this package):
//
//   - fn returns nil, Commit succeeds → nil.
//   - fn returns nil, Commit fails → oops.Code(CodeCommitFailed)
//     wrapping the underlying driver error.
//   - fn returns err, Rollback succeeds → the original err is returned
//     unchanged so the caller's error chain (oops / domain codes) is preserved.
//   - fn returns err, Rollback ALSO fails → oops.Code(CodeRollbackFailed)
//     wrapping errors.Join(originalErr, rollbackErr). Both errors are still
//     discoverable via errors.Is / errors.As because errors.Join exposes
//     Unwrap() []error.
//   - fn panics → Rollback is called best-effort and the original panic
//     is re-raised. Any rollback or Close error on that path is unavoidably
//     lost — the panic payload is the caller's primary signal and there is
//     no return channel for secondary errors.
//   - Close errors in the deferred cleanup → wrapped as
//     oops.Code(CodeCloseFailed) IFF no earlier error was already set;
//     never overrides a real failure (fn, commit, or rollback error).
//
// Commit and Rollback mark the tx as finished; Close is idempotent and
// therefore a no-op in normal paths.
func Within(ctx context.Context, t Tx, fn func(ctx context.Context, t Tx) error) (err error) {
	ctx = WithContext(ctx, t)

	defer func() {
		if cerr := t.Close(); cerr != nil && err == nil {
			err = oops.In("tx").Code(CodeCloseFailed).Wrap(cerr)
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			// Best-effort rollback. Any rollback error is unavoidably lost
			// here: we cannot return from a panicking function, and the
			// panic payload is the caller's primary signal.
			_ = t.Rollback(ctx)
			panic(r)
		}
	}()

	if err = fn(ctx, t); err != nil {
		if rerr := t.Rollback(ctx); rerr != nil {
			// NOTE: this is the ONE place in the codebase where it is
			// OK to wrap an errors.Join with oops.Wrap. The combined
			// error is only ever inspected by generic handlers (gRPC
			// interceptor fallback logging, operator dashboards); it
			// never flows through validation multi-error flattening,
			// because tx errors never carry Code("invalid_argument").
			// At the VO/aggregate layer this pattern is forbidden —
			// see AGENTS.md rule #14. errors.Is / errors.As still find
			// both leaves because errors.Join exposes Unwrap() []error.
			return oops.
				In("tx").
				Code(CodeRollbackFailed).
				Hint("fn returned an error and the subsequent rollback also failed; both errors are preserved in the join").
				Wrap(errors.Join(err, rerr))
		}
		return err
	}

	if cerr := t.Commit(ctx); cerr != nil {
		return oops.In("tx").Code(CodeCommitFailed).Wrap(cerr)
	}
	return nil
}
