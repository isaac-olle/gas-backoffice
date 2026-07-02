package command

import (
	"context"
	"fmt"
	errorDomain "github.com/isaac-olle/gas/internal/shared/app_errors"
)

var (
	ErrInvalidContext errorDomain.ErrorCode = "INVALID_CONTEXT"
	ErrInvalidCommand errorDomain.ErrorCode = "INVALID_COMMAND"
)

type Type string

func (t Type) String() string { return string(t) }

type Command interface {
	Type() Type
}

func Process[T Command, C context.Context](ctx context.Context, cmd Command, handler func(C, T) error) error {
	cCtx, ok := ctx.(C)
	if !ok {
		return errorDomain.InternalError(ErrInvalidContext).WithDescription("unexpected context type")
	}

	parsedCommand, ok := cmd.(T)
	if !ok {
		return errorDomain.InternalError(ErrInvalidCommand).
			WithDescription(fmt.Sprintf("unexpected command type %s", cmd.Type().String()))
	}

	return handler(cCtx, parsedCommand)
}
