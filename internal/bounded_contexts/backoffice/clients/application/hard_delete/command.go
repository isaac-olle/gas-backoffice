package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

const CommandType command.Type = "command.backoffice.client.hard_delete"

type Command struct {
	clientId string
	author   *vo.UserAuthor
}

func NewCommand(clientId string, author ...*vo.UserAuthor) *Command {
	var a *vo.UserAuthor
	if len(author) > 0 {
		a = author[0]
	}
	return &Command{clientId: clientId, author: a}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	if c, ok := cmd.(*Command); ok {
		c.author, _ = vo.AuthorFromContext(ctx)
	}
	return command.Process[*Command, context.Context](ctx, cmd, h.service.HardDelete)
}
