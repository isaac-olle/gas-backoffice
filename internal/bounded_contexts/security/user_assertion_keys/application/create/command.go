package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
)

const CommandType command.Type = "command.security.user_assertion_key.create"

type Command struct {
	id        string
	userId    string
	algorithm string
	publicKey string
}

func NewCommand(id, userId, algorithm, publicKey string) *Command {
	return &Command{id: id, userId: userId, algorithm: algorithm, publicKey: publicKey}
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
