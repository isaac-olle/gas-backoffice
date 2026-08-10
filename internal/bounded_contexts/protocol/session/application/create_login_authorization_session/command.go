package create_login_authorization_session

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.protocol.session.login_authorize"

// Command only carries what the frontend provides: the session id. Everything
// else (client id, redirect uri, scope) is injected into the service, and the
// response type and state are computed there. There is no PKCE in this flow.
type Command struct {
	sessionID ids.SessionId
}

func NewCommand(sessionID ids.SessionId) *Command {
	return &Command{sessionID: sessionID}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler { return &CommandHandler{service: service} }

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
