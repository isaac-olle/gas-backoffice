package logout

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.auth.logout"

type Command struct{ sessionID ids.SessionId }

func NewCommand(sessionID ids.SessionId) *Command { return &Command{sessionID: sessionID} }
func (c *Command) Type() command.Type             { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler { return &CommandHandler{service: service} }

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Logout)
}
