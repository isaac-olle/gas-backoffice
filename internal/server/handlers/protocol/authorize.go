package protocol

import (
	"net/http"

	"github.com/gin-gonic/gin"
	buildAuthResponse "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/build_authorization_response"
	sessionCreate "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/create_authorization_session"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type authorizeRequest struct {
	ResponseType        string `form:"response_type" binding:"required"`
	ClientID            string `form:"client_id" binding:"required"`
	RedirectURI         string `form:"redirect_uri" binding:"required"`
	Scope               string `form:"scope" binding:"required"`
	State               string `form:"state"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
}

type AuthorizeHandler struct {
	commandBus command.Bus
	queryBus   query.Bus
}

func NewAuthorizeHandler(cmdBus command.Bus, qryBus query.Bus) *AuthorizeHandler {
	return &AuthorizeHandler{commandBus: cmdBus, queryBus: qryBus}
}

func (h *AuthorizeHandler) Execute(c *gin.Context) (any, error) {
	var req authorizeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	sessionID := ids.NewSessionId()
	cmd := sessionCreate.NewCommand(
		sessionID,
		req.ResponseType,
		req.ClientID,
		req.RedirectURI,
		req.Scope,
		req.State,
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err := h.commandBus.Dispatch(c.Request.Context(), cmd); err != nil {
		return nil, err
	}

	return query.DispatchTyped[string](c.Request.Context(), h.queryBus, buildAuthResponse.NewQuery(sessionID))
}

func (h *AuthorizeHandler) Success(c *gin.Context, result any) error {
	c.Redirect(http.StatusFound, result.(string))
	return nil
}

func (h *AuthorizeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
