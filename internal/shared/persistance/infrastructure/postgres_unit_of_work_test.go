package infrastructure

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type uowFakeTx struct {
	pgx.Tx
	committed bool
}

func (f *uowFakeTx) Commit(context.Context) error { f.committed = true; return nil }

// When a transaction is already in the context, Do must join it (run fn with the
// same context) and NOT commit — the owner does. The pool is nil to prove it is
// never touched in this path.
func TestPostgresUnitOfWork_ReusesExistingTx(t *testing.T) {
	t.Parallel()
	existing := &uowFakeTx{}
	uow := NewPostgresUnitOfWork(nil)

	ranWithTx := false
	err := uow.Do(ContextWithTx(context.Background(), existing), func(ctx context.Context) error {
		_, ranWithTx = TxFromContext(ctx)
		return nil
	})

	require.NoError(t, err)
	assert.True(t, ranWithTx, "fn must run with the ambient transaction in context")
	assert.False(t, existing.committed, "the owner of the tx commits it, not the reused unit of work")
}

func TestPostgresUnitOfWork_ReusePropagatesError(t *testing.T) {
	t.Parallel()
	uow := NewPostgresUnitOfWork(nil)
	boom := errors.New("boom")

	err := uow.Do(ContextWithTx(context.Background(), &uowFakeTx{}), func(context.Context) error {
		return boom
	})

	assert.ErrorIs(t, err, boom)
}
