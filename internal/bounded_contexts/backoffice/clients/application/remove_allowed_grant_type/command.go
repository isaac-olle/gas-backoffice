package remove_allowed_grant_type

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.client.remove_allowed_grant_type"

type Command struct {
	clientId  ids.ClientId
	grantType string
	author    *vo.UserAuthor
}

func NewCommand(clientId ids.ClientId, grantType string, author ...*vo.UserAuthor) *Command {
	var a *vo.UserAuthor
	if len(author) > 0 {
		a = author[0]
	}
	return &Command{clientId: clientId, grantType: grantType, author: a}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct {
	service *Service
}

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	if c, ok := cmd.(*Command); ok {
		c.author, _ = vo.AuthorFromContext(ctx)
	}
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Remove)
}
