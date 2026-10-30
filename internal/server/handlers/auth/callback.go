package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create_session"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type CallbackHandler struct {
	queryBus query.Bus
}

func NewCallbackHandler(queryBus query.Bus) *CallbackHandler {
	return &CallbackHandler{queryBus: queryBus}
}

type callbackRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

func (h *CallbackHandler) Execute(c *gin.Context) (any, error) {
	var req callbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}
	return h.queryBus.Dispatch(c.Request.Context(), create_session.NewQuery(req.Code, req.State))
}

func (h *CallbackHandler) Success(c *gin.Context, result any) error {
	res := result.(create_session.Result)
	c.SetCookie(authenticators.AccessTokenCookieName, res.AccessToken, 3600, "/", "", false, true)
	c.SetCookie(authenticators.RefreshTokenCookieName, res.RefreshToken, 30*24*60*60, "/", "", false, true)
	return nil
}

func (h *CallbackHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
