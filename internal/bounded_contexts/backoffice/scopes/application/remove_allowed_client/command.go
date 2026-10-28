package remove_allowed_client

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

const CommandType command.Type = "command.backoffice.scope.remove_allowed_client"

type Command struct {
	scopeId  string
	clientId string
	author   *authors.UserAuthor
}

func NewCommand(scopeId, clientId string, author ...*authors.UserAuthor) *Command {
	var a *authors.UserAuthor
	if len(author) > 0 {
		a = author[0]
	}
	return &Command{scopeId: scopeId, clientId: clientId, author: a}
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
		c.author, _ = vo.AuthorFromContext[*authors.UserAuthor](ctx)
	}
	return command.Process[*Command, context.Context](ctx, cmd, h.service.RemoveAllowedClient)
}
