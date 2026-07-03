package infrastructure

import (
	"context"
	"fmt"
	query2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"sync"
)

var ErrHandlerNotFound = fmt.Errorf("handler not found")

type QueryBus struct {
	mu       sync.RWMutex
	handlers map[query2.Type]query2.Handler
}

func NewQueryBus() *QueryBus {
	return &QueryBus{handlers: make(map[query2.Type]query2.Handler)}
}

func (b *QueryBus) Dispatch(ctx context.Context, qry query2.Query) (any, error) {
	b.mu.RLock()
	handler, ok := b.handlers[qry.Type()]
	b.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no handler registered for query type %s: %w", qry.Type(), ErrHandlerNotFound)
	}
	return handler.Handle(ctx, qry)
}

func (b *QueryBus) Register(qryType query2.Type, handler query2.Handler) {
	b.handlers[qryType] = handler
}
