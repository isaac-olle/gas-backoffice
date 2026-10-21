package create

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.backoffice.identity_provider.create"

type Command struct {
	id                    ids.IdentityProviderId
	referenceName         string
	name                  string
	issuer                string
	clientID              string
	clientSecret          string
	authorizationEndpoint string
	tokenEndpoint         string
	userinfoEndpoint      string
	jwksURI               string
	scopes                []string
	backofficeRedirectURI string
}

func NewCommand(
	id ids.IdentityProviderId,
	referenceName, name, issuer,
	clientID, clientSecret,
	authorizationEndpoint, tokenEndpoint, userinfoEndpoint, jwksURI string,
	scopes []string,
	backofficeRedirectURI string,
) *Command {
	return &Command{
		id:                    id,
		referenceName:         referenceName,
		name:                  name,
		issuer:                issuer,
		clientID:              clientID,
		clientSecret:          clientSecret,
		authorizationEndpoint: authorizationEndpoint,
		tokenEndpoint:         tokenEndpoint,
		userinfoEndpoint:      userinfoEndpoint,
		jwksURI:               jwksURI,
		scopes:                scopes,
		backofficeRedirectURI: backofficeRedirectURI,
	}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler { return &CommandHandler{service: service} }

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
