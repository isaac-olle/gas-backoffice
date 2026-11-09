package middlewares_test

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeTx struct {
	pgx.Tx
	commitErr  error
	committed  bool
	rolledBack bool
}

func (f *fakeTx) Commit(_ context.Context) error {
	f.committed = true
	return f.commitErr
}

func (f *fakeTx) Rollback(_ context.Context) error {
	f.rolledBack = true
	return nil
}

func (f *fakeTx) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

type fakeTxPool struct {
	tx       *fakeTx
	beginErr error
}

func newFakeTxPool(beginErr, commitErr error) *fakeTxPool {
	return &fakeTxPool{
		tx:       &fakeTx{commitErr: commitErr},
		beginErr: beginErr,
	}
}

func (f *fakeTxPool) Begin(_ context.Context) (pgx.Tx, error) {
	if f.beginErr != nil {
		return nil, f.beginErr
	}
	return f.tx, nil
}

func (f *fakeTxPool) committed_() bool  { return f.tx.committed }
func (f *fakeTxPool) rolledBack_() bool { return f.tx.rolledBack }
