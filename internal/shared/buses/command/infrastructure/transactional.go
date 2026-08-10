package infrastructure

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

// TransactionalCommandBus decorates a command.Bus so every command runs inside a
// unit of work: the aggregate write and the outbox event insert commit
// atomically (or roll back together). It delegates the whole transaction
// protocol to the UnitOfWork — the single owner of begin / reuse / commit — so
// this decorator knows nothing about the DB driver.
type TransactionalCommandBus struct {
	inner command.Bus
	uow   persistance.UnitOfWork
}

func NewTransactionalCommandBus(inner command.Bus, uow persistance.UnitOfWork) *TransactionalCommandBus {
	return &TransactionalCommandBus{inner: inner, uow: uow}
}

func (b *TransactionalCommandBus) Register(cmdType command.Type, handler command.Handler) {
	b.inner.Register(cmdType, handler)
}

func (b *TransactionalCommandBus) Dispatch(ctx context.Context, cmd command.Command) error {
	return b.uow.Do(ctx, func(ctx context.Context) error {
		return b.inner.Dispatch(ctx, cmd)
	})
}
