package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
)

const CommandType command.Type = "command.security.client_secret.create"

type Command struct {
	id          string
	clientId    string
	plainSecret string
}

func NewCommand(id, clientId, plainSecret string) *Command {
	return &Command{id: id, clientId: clientId, plainSecret: plainSecret}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct {
	service *Service
}

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
