package add_redirect_uri

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.client.add_redirect_uri"

type Command struct {
	clientId    ids.ClientId
	redirectUri string
	author      *vo.UserAuthor
}

func NewCommand(clientId ids.ClientId, redirectUri string, author ...*vo.UserAuthor) *Command {
	var a *vo.UserAuthor
	if len(author) > 0 {
		a = author[0]
	}
	return &Command{clientId: clientId, redirectUri: redirectUri, author: a}
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
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Add)
}
