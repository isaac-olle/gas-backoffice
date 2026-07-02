package clients

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	clientGet "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/application/get"
)

type GetClientHandler struct {
	queryBus query.Bus
	clientId func(context.Context) (ids.ClientId, bool)
}

func NewGetClientHandler(queryBus query.Bus, clientId func(context.Context) (ids.ClientId, bool)) *GetClientHandler {
	return &GetClientHandler{queryBus: queryBus, clientId: clientId}
}

func (h *GetClientHandler) Execute(c *gin.Context) (any, error) {
	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}
	return h.queryBus.Dispatch(c.Request.Context(), clientGet.NewQuery(clientId.String()))
}

func (h *GetClientHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
