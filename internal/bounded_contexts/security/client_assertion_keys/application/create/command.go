package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

const CommandType command.Type = "command.security.client_assertion_key.create"

type Command struct {
	id        ids2.ClientAssertionKeyId
	clientId  ids2.ClientId
	algorithm string
	publicKey string
}

func NewCommand(id ids2.ClientAssertionKeyId, clientId ids2.ClientId, algorithm, publicKey string) *Command {
	return &Command{id: id, clientId: clientId, algorithm: algorithm, publicKey: publicKey}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct {
	service *Service
}

func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
