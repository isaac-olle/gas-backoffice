package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
)

const CommandType command.Type = "command.security.admin_assertion_key.create"

type Command struct {
	id        string
	adminId   string
	algorithm string
	publicKey string
}

func NewCommand(id, adminId, algorithm, publicKey string) *Command {
	return &Command{id: id, adminId: adminId, algorithm: algorithm, publicKey: publicKey}
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
