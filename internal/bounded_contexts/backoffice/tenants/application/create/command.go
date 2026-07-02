package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.tenant.create"

type Command struct {
	id        ids.TenantId
	name      string
	ownership tenants.Ownership
}

func NewCommand(id ids.TenantId, name string, ownership tenants.Ownership) *Command {
	return &Command{id: id, name: name, ownership: ownership}
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
