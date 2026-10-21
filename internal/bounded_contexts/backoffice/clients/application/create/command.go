package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.client.create"

type Command struct {
	clientID          ids.ClientId
	name              string
	ownership         clients.Ownership
	allowedGrantTypes []string
	redirectUris      []string
	public            bool
	isSystem          bool
}

func NewCommand(id ids.ClientId, name string, ownership clients.Ownership, allowedGrantTypes, redirectUris []string, public bool, isSystem bool) *Command {
	return &Command{clientID: id, name: name, ownership: ownership, allowedGrantTypes: allowedGrantTypes, redirectUris: redirectUris, public: public, isSystem: isSystem}
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
