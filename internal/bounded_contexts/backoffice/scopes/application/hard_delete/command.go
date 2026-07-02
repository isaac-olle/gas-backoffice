package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.scope.hard_delete"

type Command struct {
	scopeId ids.ScopeId
}

func NewCommand(scopeId ids.ScopeId) *Command {
	return &Command{scopeId: scopeId}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.HardDelete)
}
