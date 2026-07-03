package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	command2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"sync"
)

var ErrHandlerNotFound app_errors.ErrorCode = "HANDLER_NOT_FOUND"

type CommandBus struct {
	mu       sync.RWMutex
	handlers map[command2.Type]command2.Handler
}

func NewCommandBus() *CommandBus {
	return &CommandBus{handlers: make(map[command2.Type]command2.Handler)}
}

func (b *CommandBus) Dispatch(ctx context.Context, cmd command2.Command) error {
	b.mu.RLock()
	handler, ok := b.handlers[cmd.Type()]
	b.mu.RUnlock()

	if !ok {
		return app_errors.InternalError(ErrHandlerNotFound).
			WithDescription(fmt.Sprintf("handler not found for command type %s", cmd.Type().String()))
	}

	return handler.Handle(ctx, cmd)
}

func (b *CommandBus) Register(cmdType command2.Type, handler command2.Handler) {
	b.handlers[cmdType] = handler
}
