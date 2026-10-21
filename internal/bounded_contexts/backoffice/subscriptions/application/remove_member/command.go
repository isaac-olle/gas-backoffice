package remove_member

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.subscription.remove_member"

type Command struct {
	id     ids.SubscriptionId
	userId ids.UserId
}

func NewCommand(id ids.SubscriptionId, userId ids.UserId) *Command {
	return &Command{id: id, userId: userId}
}
func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(s *Service) *CommandHandler { return &CommandHandler{service: s} }
func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Remove)
}
