package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	userList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/list"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ListUsersHandler struct {
	queryBus query.Bus
	userId   func(context.Context) (ids.UserId, bool)
}

func NewListUsersHandler(queryBus query.Bus, userId func(context.Context) (ids.UserId, bool)) *ListUsersHandler {
	return &ListUsersHandler{queryBus: queryBus, userId: userId}
}

func (h *ListUsersHandler) Execute(c *gin.Context) (any, error) {
	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}
	result, err := query.DispatchTyped[[]*tenants.Tenant](c.Request.Context(), h.queryBus, userList.NewQuery(userId))
	if err != nil {
		return nil, err
	}
	primitives := make([]tenants.TenantPrimitives, 0, len(result))
	for _, item := range result {
		primitives = append(primitives, item.ToPrimitives())
	}
	return primitives, nil
}

func (h *ListUsersHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *ListUsersHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
