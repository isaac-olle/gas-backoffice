package identity_providers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	idpCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateIdentityProviderHandler struct {
	commandExecutor command.CommandExecutor
}

func NewCreateIdentityProviderHandler(commandExecutor command.CommandExecutor) *CreateIdentityProviderHandler {
	return &CreateIdentityProviderHandler{commandExecutor: commandExecutor}
}

type createIdentityProviderRequest struct {
	ReferenceName         string   `json:"reference_name" binding:"required"`
	Name                  string   `json:"name" binding:"required"`
	Issuer                string   `json:"issuer" binding:"required"`
	ClientID              string   `json:"client_id" binding:"required"`
	ClientSecret          string   `json:"client_secret" binding:"required"`
	AuthorizationEndpoint string   `json:"authorization_endpoint" binding:"required"`
	TokenEndpoint         string   `json:"token_endpoint" binding:"required"`
	UserinfoEndpoint      string   `json:"userinfo_endpoint"`
	JwksURI               string   `json:"jwks_uri" binding:"required"`
	Scopes                []string `json:"scopes" binding:"required"`
	BackofficeRedirectURI string   `json:"backoffice_redirect_uri"`
}

func (h *CreateIdentityProviderHandler) Execute(c *gin.Context) (any, error) {
	var req createIdentityProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(),
		idpCreate.NewCommand(
			ids.NewIdentityProviderId(),
			req.ReferenceName, req.Name, req.Issuer,
			req.ClientID, req.ClientSecret,
			req.AuthorizationEndpoint, req.TokenEndpoint, req.UserinfoEndpoint, req.JwksURI,
			req.Scopes,
			req.BackofficeRedirectURI,
		),
	)
}

func (h *CreateIdentityProviderHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateIdentityProviderHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
