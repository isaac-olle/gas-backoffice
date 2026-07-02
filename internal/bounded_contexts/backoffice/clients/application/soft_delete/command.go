package soft_delete

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
)

const CommandType command.Type = "command.backoffice.client.soft_delete"

type Command struct {
	clientId string
}

func NewCommand(clientId string) *Command {
	return &Command{clientId: clientId}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.SoftDelete)
}
