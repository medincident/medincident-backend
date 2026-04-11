package postgres

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"sync/atomic"

	"github.com/jackc/pgx/v5"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

// Error codes emitted by the savepoint machinery and the unwrap helper.
const (
	CodeSavepointFailed         = "postgres_savepoint_failed"
	CodeSavepointRollbackFailed = "postgres_savepoint_rollback_failed"
	CodeSavepointReleaseFailed  = "postgres_savepoint_release_failed"
	CodeSavepointInvalidName    = "postgres_savepoint_invalid_name"

	CodeNilTx     = "postgres_nil_tx"
	CodeForeignTx = "postgres_foreign_tx"
)

// pgxTx wraps pgx.Tx and implements tx.Tx. Repositories in this package
// recover the underlying pgx.Tx via the package-private unwrap helper.
type pgxTx struct {
	raw      pgx.Tx
	finished atomic.Bool
}

// Commit runs the underlying pgx commit. Idempotent — a second call is
// a no-op so the deferred Close() from tx.Within cannot double-commit.
func (t *pgxTx) Commit(ctx context.Context) error {
	if !t.finished.CompareAndSwap(false, true) {
		return nil
	}
	if err := t.raw.Commit(ctx); err != nil {
		return oops.In("storage.postgres").Code(tx.CodeCommitFailed).Wrap(err)
	}
	return nil
}

// Rollback runs the underlying pgx rollback. Idempotent — same reason.
func (t *pgxTx) Rollback(ctx context.Context) error {
	if !t.finished.CompareAndSwap(false, true) {
		return nil
	}
	if err := t.raw.Rollback(ctx); err != nil {
		return oops.In("storage.postgres").Code(tx.CodeRollbackFailed).Wrap(err)
	}
	return nil
}

// Close is an idempotent rollback. Safe to defer.
//
// io.Closer's Close() has no context parameter, so we pass a background
// context. In practice Close is only called from deferred cleanup after
// the main operation has either committed (no-op) or returned an error
// (already-finished, also no-op) — the background context is only used
// if someone forgets to call Commit/Rollback, which should not happen
// via tx.Within.
func (t *pgxTx) Close() error {
	if t.finished.Load() {
		return nil
	}
	return t.Rollback(context.Background())
}

// Savepoint opens a SQL savepoint. See quoteSavepointName for the name
// validation rules.
func (t *pgxTx) Savepoint(ctx context.Context, name string) error {
	ident, err := quoteSavepointName(name)
	if err != nil {
		return err
	}
	if _, err := t.raw.Exec(ctx, "SAVEPOINT "+ident); err != nil {
		return oops.In("storage.postgres").Code(CodeSavepointFailed).
			With("name", name).Wrap(err)
	}
	return nil
}

// RollbackSavepoint rolls back to a previously opened savepoint.
func (t *pgxTx) RollbackSavepoint(ctx context.Context, name string) error {
	ident, err := quoteSavepointName(name)
	if err != nil {
		return err
	}
	if _, err := t.raw.Exec(ctx, "ROLLBACK TO SAVEPOINT "+ident); err != nil {
		return oops.In("storage.postgres").Code(CodeSavepointRollbackFailed).
			With("name", name).Wrap(err)
	}
	return nil
}

// ReleaseSavepoint releases (commits) a previously opened savepoint.
func (t *pgxTx) ReleaseSavepoint(ctx context.Context, name string) error {
	ident, err := quoteSavepointName(name)
	if err != nil {
		return err
	}
	if _, err := t.raw.Exec(ctx, "RELEASE SAVEPOINT "+ident); err != nil {
		return oops.In("storage.postgres").Code(CodeSavepointReleaseFailed).
			With("name", name).Wrap(err)
	}
	return nil
}

// compile-time assertion that pgxTx satisfies the abstract interface
var _ tx.Tx = (*pgxTx)(nil)

// savepointNamePattern allows alphanumerics and underscore only, starting
// with a letter or underscore. This prevents any form of SQL injection
// via savepoint names without needing to understand every quoting nuance
// in Postgres identifier rules.
var savepointNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,62}$`)

// quoteSavepointName validates a savepoint name against a strict allow-list
// and returns it in quoted form (so case is preserved). Rejecting malformed
// names at this layer means no SQL-composition pathway can be exploited:
// the concatenated string either matches the regex or we return an error
// before touching the database.
func quoteSavepointName(name string) (string, error) {
	if !savepointNamePattern.MatchString(name) {
		return "", oops.In("storage.postgres").
			Code(CodeSavepointInvalidName).
			Public("Savepoint name is invalid.").
			With("name", name).
			Hint("savepoint names must match ^[A-Za-z_][A-Za-z0-9_]{0,62}$").
			Errorf("invalid savepoint name")
	}
	return strconv.Quote(name), nil
}

// unwrap casts an abstract tx.Tx back to *pgxTx. Used by repository
// implementations in this package to reach the underlying pgx.Tx for
// query execution.
func unwrap(t tx.Tx) (*pgxTx, error) {
	if t == nil {
		return nil, oops.In("storage.postgres").
			Code(CodeNilTx).
			Errorf("nil transaction")
	}
	pg, ok := t.(*pgxTx)
	if !ok {
		return nil, oops.In("storage.postgres").
			Code(CodeForeignTx).
			With("got_type", fmt.Sprintf("%T", t)).
			Hint("tx.Tx was minted by a different beginner; check DI wiring").
			Errorf("foreign transaction implementation")
	}
	return pg, nil
}
