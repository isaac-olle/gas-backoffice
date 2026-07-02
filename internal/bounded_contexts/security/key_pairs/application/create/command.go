package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
)

const CommandType command.Type = "command.security.key_pair.create"

type Command struct {
	id                  string
	protectedResourceId string
	algorithm           string
}

func NewCommand(id, protectedResourceId, algorithm string) *Command {
	return &Command{id: id, protectedResourceId: protectedResourceId, algorithm: algorithm}
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
