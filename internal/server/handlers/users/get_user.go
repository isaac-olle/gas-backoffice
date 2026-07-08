package users

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	userGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/get"
)

type GetUserHandler struct {
	queryBus query.Bus
	userId   func(context.Context) (ids.UserId, bool)
}

func NewGetUserHandler(queryBus query.Bus, userId func(context.Context) (ids.UserId, bool)) *GetUserHandler {
	return &GetUserHandler{queryBus: queryBus, userId: userId}
}

func (h *GetUserHandler) Execute(c *gin.Context) (any, error) {
	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}
	user, err := query.DispatchTyped[*users.User](c.Request.Context(), h.queryBus, userGet.NewQuery(userId))
	if err != nil {
		return nil, err
	}
	return user.ToPrimitives(), nil
}

func (h *GetUserHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
