package create_tenant

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.tenant.create"

type Command struct {
	id       ids.UserId
	adminId  ids.UserId
	name     string
	isSystem bool
}

func NewCommand(id ids.UserId, adminId ids.UserId, name string, isSystem bool) *Command {
	return &Command{id: id, adminId: adminId, name: name, isSystem: isSystem}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler { return &CommandHandler{service: service} }

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
