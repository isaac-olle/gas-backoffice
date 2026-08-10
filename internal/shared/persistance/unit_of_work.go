package persistance

import "context"

// UnitOfWork runs fn inside a single transaction. The context passed to fn
// carries the transaction, which repositories pick up transparently via
// TxFromContext. If fn returns an error the work is rolled back; otherwise it is
// committed. The signature is storage-agnostic (no driver types leak), so it is
// a genuine port.
type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
