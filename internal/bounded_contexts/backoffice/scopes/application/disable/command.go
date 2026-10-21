package disable

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
)

const CommandType command.Type = "command.backoffice.scope.disable"

type Command struct {
	scopeId             string
	protectedResourceId string
}

func NewCommand(scopeId, protectedResourceId string) *Command {
	return &Command{scopeId: scopeId, protectedResourceId: protectedResourceId}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Disable)
}
