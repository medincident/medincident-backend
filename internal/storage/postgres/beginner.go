package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

// Beginner implements tx.Beginner on top of a Pool. DI provides it as
// the tx.Beginner interface — consumers never see the concrete type.
type Beginner struct {
	pool *Pool
}

// NewBeginner returns a Beginner bound to the given pool.
func NewBeginner(pool *Pool) *Beginner { return &Beginner{pool: pool} }

// Begin opens a new pgx transaction with default isolation and wraps it
// as an abstract tx.Tx.
func (b *Beginner) Begin(ctx context.Context) (tx.Tx, error) {
	raw, err := b.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, oops.
			In("storage.postgres").
			Code(tx.CodeBeginFailed).
			Wrap(err)
	}
	return &pgxTx{raw: raw}, nil
}

var _ tx.Beginner = (*Beginner)(nil)
