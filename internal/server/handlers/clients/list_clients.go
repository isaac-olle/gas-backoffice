package clients

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	clientList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/list"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ListClientsHandler struct {
	queryBus       query.Bus
	subscriptionId func(context.Context) (ids.SubscriptionId, bool)
}

func NewListClientsHandler(queryBus query.Bus, subscriptionId func(context.Context) (ids.SubscriptionId, bool)) *ListClientsHandler {
	return &ListClientsHandler{queryBus: queryBus, subscriptionId: subscriptionId}
}

func (h *ListClientsHandler) Execute(c *gin.Context) (any, error) {
	subscriptionId, ok := h.subscriptionId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SUBSCRIPTION_ID")
	}

	result, err := query.DispatchTyped[[]*clients.Client](c.Request.Context(), h.queryBus, clientList.NewQuery(subscriptionId))
	if err != nil {
		return nil, err
	}
	primitives := make([]clients.ClientPrimitives, 0, len(result))
	for _, item := range result {
		primitives = append(primitives, item.ToPrimitives())
	}
	return primitives, nil
}

func (h *ListClientsHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *ListClientsHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
