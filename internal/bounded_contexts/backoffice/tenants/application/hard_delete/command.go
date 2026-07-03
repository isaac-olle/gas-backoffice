package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
)

const CommandType command.Type = "command.backoffice.tenant.hard_delete"

type Command struct {
	tenantId string
}

func NewCommand(tenantId string) *Command {
	return &Command{tenantId: tenantId}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.HardDelete)
}
