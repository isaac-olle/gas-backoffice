package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.client.create"

type Command struct {
	clientID     ids.ClientId
	name         string
	ownership    clients.Ownership
	allowedFlows []string
	redirectUris []string
}

func NewCommand(id ids.ClientId, name string, ownership clients.Ownership, allowedFlows, redirectUris []string) *Command {
	return &Command{clientID: id, name: name, ownership: ownership, allowedFlows: allowedFlows, redirectUris: redirectUris}
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
