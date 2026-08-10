package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authenticateWithIdp "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_with_idp"
	createLogin "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/create_login_authorization_session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// LoginHandler starts a federated login. The whole OIDC flow lives on the
// backend and is driven by two use cases: create_login_authorization_session
// persists the login session (state included), and authenticate_with_idp asks
// the identity provider adapter for the authorization URL to redirect to. Both
// the session id and the identity provider id come from the frontend.
type LoginHandler struct {
	commandExecutor command.CommandExecutor
	queryBus        query.Bus
}

func NewLoginHandler(commandExecutor command.CommandExecutor, queryBus query.Bus) *LoginHandler {
	return &LoginHandler{commandExecutor: commandExecutor, queryBus: queryBus}
}

func (h *LoginHandler) Execute(c *gin.Context) (any, error) {
	sessionID, err := ids.NewSessionIdFromString(c.Param("session_id"))
	if err != nil {
		return nil, app_errors.BadRequestError("INVALID_SESSION_ID").
			WithDescription("invalid session id").
			WithCause(err)
	}
	idpID := c.Param("idp_id")

	if err = h.commandExecutor.SingleExecution(
		c.Request.Context(), createLogin.NewCommand(sessionID)); err != nil {
		return nil, err
	}

	authorizationURL, err := query.DispatchTyped[string](
		c.Request.Context(), h.queryBus, authenticateWithIdp.NewQuery(sessionID, idpID),
	)
	if err != nil {
		return nil, err
	}

	return authorizationURL, nil
}

func (h *LoginHandler) Success(c *gin.Context, result any) error {
	c.Redirect(http.StatusFound, result.(string))
	return nil
}

func (h *LoginHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
