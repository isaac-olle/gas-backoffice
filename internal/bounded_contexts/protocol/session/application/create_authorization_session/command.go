package create_authorization_session

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.protocol.session.authorize"

type Command struct {
	sessionID           ids.SessionId
	responseType        string
	clientId            string
	redirectUri         string
	scope               string
	state               string
	codeChallenge       string
	codeChallengeMethod string
}

func NewCommand(sessionID ids.SessionId, responseType, clientId, redirectUri, scope, state, codeChallenge, codeChallengeMethod string) *Command {
	return &Command{
		sessionID:           sessionID,
		responseType:        responseType,
		clientId:            clientId,
		redirectUri:         redirectUri,
		scope:               scope,
		state:               state,
		codeChallenge:       codeChallenge,
		codeChallengeMethod: codeChallengeMethod,
	}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler { return &CommandHandler{service: service} }

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
