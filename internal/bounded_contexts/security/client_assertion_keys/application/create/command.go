package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.security.client_assertion_key.create"

type Command struct {
	id        ids2.ClientAssertionKeyId
	clientId  ids2.ClientId
	algorithm string
	publicKey string
	author    *vo.UserAuthor
}

func NewCommand(id ids2.ClientAssertionKeyId, clientId ids2.ClientId, algorithm, publicKey string, author ...*vo.UserAuthor) *Command {
	var a *vo.UserAuthor
	if len(author) > 0 {
		a = author[0]
	}
	return &Command{id: id, clientId: clientId, algorithm: algorithm, publicKey: publicKey, author: a}
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
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
