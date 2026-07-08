package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
)

const CommandType command.Type = "command.backoffice.user.hard_delete"

type Command struct {
	userId string
}

func NewCommand(userId string) *Command {
	return &Command{userId: userId}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.HardDelete)
}
