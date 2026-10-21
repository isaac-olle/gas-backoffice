package change_name

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.admin.change_name"

type Command struct {
	adminId  ids.UserId
	name     string
	callerId ids.UserId
}

func NewCommand(adminId ids.UserId, name string, callerId ids.UserId) *Command {
	return &Command{adminId: adminId, name: name, callerId: callerId}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.ChangeName)
}
