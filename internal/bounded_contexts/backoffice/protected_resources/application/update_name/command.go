package update_name

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.protected_resource.update_name"

type Command struct {
	protectedResourceId ids.ProtectedResourceId
	name                string
}

func NewCommand(protectedResourceId ids.ProtectedResourceId, name string) *Command {
	return &Command{protectedResourceId: protectedResourceId, name: name}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct {
	service *Service
}

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.UpdateName)
}
