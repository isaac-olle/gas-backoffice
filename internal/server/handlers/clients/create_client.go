package clients

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	clientCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
	scopeAddAllowedClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	clientAssertionKeyCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/create"
	clientSecretCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateClientHandler struct {
	commandExecutor command.CommandExecutor
	clientId        func(context.Context) (ids.ClientId, bool)
	subscriptionId  func(context.Context) (ids.SubscriptionId, bool)
}

func NewCreateClientHandler(
	commandExecutor command.CommandExecutor,
	clientId func(context.Context) (ids.ClientId, bool),
	subscriptionId func(context.Context) (ids.SubscriptionId, bool),
) *CreateClientHandler {
	return &CreateClientHandler{
		commandExecutor: commandExecutor,
		clientId:        clientId,
		subscriptionId:  subscriptionId,
	}
}

type createClientRequest struct {
	Name              string                       `json:"name" binding:"required"`
	AllowedGrantTypes []string                     `json:"allowed_grant_types" binding:"required,min=1"`
	AllowedScopes     []allowedScopeRequest        `json:"allowed_scopes"`
	RedirectUris      []string                     `json:"redirect_uris" binding:"required,min=1"`
	Public            bool                         `json:"public"`
	Authentication    *clientAuthenticationRequest `json:"authentication"`
}

type allowedScopeRequest struct {
	ScopeId string `json:"scope_id" binding:"required"`
}

type clientAuthenticationRequest struct {
	ClientSecret       *clientSecretRequest       `json:"client_secret"`
	ClientAssertionKey *clientAssertionKeyRequest `json:"client_assertion_key"`
}

type clientSecretRequest struct {
	Id     string `json:"id" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type clientAssertionKeyRequest struct {
	Id     string `json:"id" binding:"required"`
	Key    string `json:"key" binding:"required"`
	Format string `json:"format" binding:"required"`
}

func (h *CreateClientHandler) Execute(c *gin.Context) (any, error) {
	var req createClientRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	subscriptionId, ok := h.subscriptionId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SUBSCRIPTION_ID")
	}

	commands, err := createClientCommands(req, clientId, subscriptionId)
	if err != nil {
		return nil, err
	}
	return nil, h.commandExecutor.MultipleExecution(c.Request.Context(), commands...)
}

func createClientCommands(req createClientRequest, clientId ids.ClientId, subscriptionId ids.SubscriptionId) ([]command.Command, error) {
	commands := []command.Command{
		clientCreate.NewCommand(clientId, req.Name, clients.NewOwnership(subscriptionId), req.AllowedGrantTypes, req.RedirectUris, req.Public, false),
	}

	scopeCommands, err := allowedScopeCommands(req.AllowedScopes, clientId)
	if err != nil {
		return nil, err
	}
	commands = append(commands, scopeCommands...)

	authenticationCommands, err := authenticationCommandsFor(req.Authentication, clientId)
	if err != nil {
		return nil, err
	}
	return append(commands, authenticationCommands...), nil
}

func allowedScopeCommands(allowedScopes []allowedScopeRequest, clientId ids.ClientId) ([]command.Command, error) {
	commands := make([]command.Command, 0, len(allowedScopes))
	for _, allowedScope := range allowedScopes {
		scopeId, err := ids.NewScopeIdFromString(allowedScope.ScopeId)
		if err != nil {
			return nil, app_errors.BadRequestError("INVALID_SCOPE_ID").WithCause(err, false)
		}
		commands = append(commands, scopeAddAllowedClient.NewCommand(scopeId.String(), clientId.String()))
	}
	return commands, nil
}

func authenticationCommandsFor(authentication *clientAuthenticationRequest, clientId ids.ClientId) ([]command.Command, error) {
	if authentication == nil || (authentication.ClientSecret == nil && authentication.ClientAssertionKey == nil) {
		return nil, app_errors.BadRequestError("CLIENT_AUTHENTICATION_REQUIRED").WithDescription("authentication requires a client_secret or client_assertion_key")
	}

	commands := make([]command.Command, 0, 2)
	secretCommand, err := clientSecretCommand(authentication.ClientSecret, clientId)
	if err != nil {
		return nil, err
	}
	if secretCommand != nil {
		commands = append(commands, secretCommand)
	}

	keyCommand, err := clientAssertionKeyCommand(authentication.ClientAssertionKey, clientId)
	if err != nil {
		return nil, err
	}
	if keyCommand != nil {
		commands = append(commands, keyCommand)
	}
	return commands, nil
}

func clientSecretCommand(request *clientSecretRequest, clientId ids.ClientId) (command.Command, error) {
	if request == nil {
		return nil, nil
	}
	secretId, err := ids.NewClientSecretIdFromString(request.Id)
	if err != nil {
		return nil, app_errors.BadRequestError("INVALID_CLIENT_SECRET_ID").WithCause(err, false)
	}
	return clientSecretCreate.NewCommand(secretId, clientId, request.Secret), nil
}

func clientAssertionKeyCommand(request *clientAssertionKeyRequest, clientId ids.ClientId) (command.Command, error) {
	if request == nil {
		return nil, nil
	}
	if !strings.EqualFold(request.Format, "PEM") {
		return nil, app_errors.BadRequestError("UNSUPPORTED_CLIENT_ASSERTION_KEY_FORMAT").WithDescription("only PEM client assertion keys are supported")
	}
	keyId, err := ids.NewClientAssertionKeyIdFromString(request.Id)
	if err != nil {
		return nil, app_errors.BadRequestError("INVALID_CLIENT_ASSERTION_KEY_ID").WithCause(err, false)
	}
	return clientAssertionKeyCreate.NewCommand(keyId, clientId, "RS256", request.Key), nil
}

func (h *CreateClientHandler) Success(c *gin.Context, _ any) error {
	c.JSON(http.StatusCreated, gin.H{})
	return nil
}

func (h *CreateClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
