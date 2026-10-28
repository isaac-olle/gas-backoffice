package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.scope.hard_delete"

type Command struct {
	scopeId ids.ScopeId
	author  *authors.UserAuthor
}

func NewCommand(scopeId ids.ScopeId, author ...*authors.UserAuthor) *Command {
	var a *authors.UserAuthor
	if len(author) > 0 {
		a = author[0]
	}
	return &Command{scopeId: scopeId, author: a}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	if c, ok := cmd.(*Command); ok {
		c.author, _ = vo.AuthorFromContext[*authors.UserAuthor](ctx)
	}
	return command.Process[*Command, context.Context](ctx, cmd, h.service.HardDelete)
}
