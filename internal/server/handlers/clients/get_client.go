package clients

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	clientGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
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

	client, err := query.DispatchTyped[*clients.Client](c.Request.Context(), h.queryBus, clientGet.NewQuery(clientId.String()))
	if err != nil {
		return nil, err
	}
	return client.ToPrimitives(), nil
}

func (h *GetClientHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
