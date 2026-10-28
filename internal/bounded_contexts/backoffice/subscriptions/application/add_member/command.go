package add_member

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.subscription.add_member"

type Command struct {
	id     ids.SubscriptionId
	userId ids.UserId
	role   vo.SubscriptionRole
}

func NewCommand(id ids.SubscriptionId, userId ids.UserId, role vo.SubscriptionRole) *Command {
	return &Command{id: id, userId: userId, role: role}
}
func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(s *Service) *CommandHandler { return &CommandHandler{service: s} }
func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Add)
}
