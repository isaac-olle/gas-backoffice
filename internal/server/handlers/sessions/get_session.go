package sessions

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	sessionGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type GetSessionHandler struct {
	queryBus  query.Bus
	sessionId func(context.Context) (ids.SessionId, bool)
}

func NewGetSessionHandler(queryBus query.Bus, sessionId func(context.Context) (ids.SessionId, bool)) *GetSessionHandler {
	return &GetSessionHandler{queryBus: queryBus, sessionId: sessionId}
}

func (h *GetSessionHandler) Execute(c *gin.Context) (any, error) {
	sessionId, ok := h.sessionId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SESSION_ID")
	}
	session, err := query.DispatchTyped[*auth.Session](c.Request.Context(), h.queryBus, sessionGet.NewQuery(sessionId))
	if err != nil {
		return nil, err
	}
	return session.ToPrimitives(), nil
}

func (h *GetSessionHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetSessionHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
