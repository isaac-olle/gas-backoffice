package add_allowed_client

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
)

const CommandType command.Type = "command.backoffice.scope.add_allowed_client"

type Command struct {
	scopeId             string
	clientId            string
	protectedResourceId string
}

func NewCommand(scopeId, clientId, protectedResourceId string) *Command {
	return &Command{scopeId: scopeId, clientId: clientId, protectedResourceId: protectedResourceId}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct {
	service *Service
}

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.AddAllowedClient)
}
