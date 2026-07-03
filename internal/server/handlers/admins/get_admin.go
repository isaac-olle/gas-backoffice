package admins

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	adminGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/get"
)

type GetAdminHandler struct {
	queryBus query.Bus
	adminId  func(context.Context) (ids.AdminId, bool)
}

func NewGetAdminHandler(queryBus query.Bus, adminId func(context.Context) (ids.AdminId, bool)) *GetAdminHandler {
	return &GetAdminHandler{queryBus: queryBus, adminId: adminId}
}

func (h *GetAdminHandler) Execute(c *gin.Context) (any, error) {
	adminId, ok := h.adminId(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_ADMIN_ID")
	}
	admin, err := query.DispatchTyped[*admins.Admin](c.Request.Context(), h.queryBus, adminGet.NewQuery(adminId))
	if err != nil {
		return nil, err
	}
	return admin.ToPrimitives(), nil
}

func (h *GetAdminHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetAdminHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
