package tx

import (
	"context"
	"io"
)

// Tx is a driver-agnostic transaction handle. Implementations live in the
// storage adapter packages (e.g. internal/storage/postgres).
type Tx interface {
	io.Closer
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	Savepoint(ctx context.Context, name string) error
	RollbackSavepoint(ctx context.Context, name string) error
	ReleaseSavepoint(ctx context.Context, name string) error
}

// Beginner opens new transactions. Implementations return an error
// wrapped with oops.Code(CodeBeginFailed) when the underlying driver
// cannot start a transaction (connection lost, pool exhausted, context
// cancelled, etc.). A successful Begin returns a live Tx that the
// caller MUST eventually Commit, Rollback, or Close — the
// recommended path is to pass the returned Tx to tx.Within, which
// handles the full lifecycle.
type Beginner interface {
	Begin(ctx context.Context) (Tx, error)
}

// CodeBeginFailed is the public vocabulary for a Begin failure emitted
// by Beginner implementations (e.g. storage/postgres.pgxBeginner). It
// lives here — alongside the interface that defines the contract —
// because the implementation in a sibling package imports this symbol
// as part of the tx vocabulary, not as a package-private detail.
const CodeBeginFailed = "tx_begin_failed"
