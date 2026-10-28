package create_user

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.user.create"

type Command struct {
	userId   ids.UserId
	authorId ids.UserId
	name     string
	isSystem bool
	email    string
}

func NewCommand(userId ids.UserId, authorId ids.UserId, name string, isSystem bool, email ...string) *Command {
	var e string
	if len(email) > 0 {
		e = email[0]
	}
	return &Command{userId: userId, authorId: authorId, name: name, isSystem: isSystem, email: e}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler { return &CommandHandler{service: service} }

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
