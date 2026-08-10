package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

// OIDCCallbackHandler receives the IDP redirect (?code&state), validates state
// against the oidc_state cookie, exchanges the code for tokens (with the client
// secret), verifies the id_token, resolves/provisions the user, creates the
// session, mints our own access token and hands it back to the frontend
// (302 + Set-Cookie HttpOnly).
//
// TODO: implementation pending design decisions (provisioning, expires_at,
// cookie delivery). The login side (state/nonce cookies) is already in place.
type OIDCCallbackHandler struct {
	commandExecutor command.CommandExecutor
	queryBus        query.Bus
}

func NewOIDCCallbackHandler(commandExecutor command.CommandExecutor, queryBus query.Bus) *OIDCCallbackHandler {
	return &OIDCCallbackHandler{commandExecutor: commandExecutor, queryBus: queryBus}
}

func (h *OIDCCallbackHandler) Execute(c *gin.Context) (any, error) {
	_ = c.Param("idp_id")
	_ = c.Query("code")
	_ = c.Query("state")
	return nil, nil
}

func (h *OIDCCallbackHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNotImplemented)
	return nil
}

func (h *OIDCCallbackHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
