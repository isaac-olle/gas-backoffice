package protected_resources

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	prList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/list"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ListProtectedResourcesHandler struct {
	queryBus       query.Bus
	subscriptionId func(context.Context) (ids.SubscriptionId, bool)
}

func NewListProtectedResourcesHandler(queryBus query.Bus, subscriptionId func(context.Context) (ids.SubscriptionId, bool)) *ListProtectedResourcesHandler {
	return &ListProtectedResourcesHandler{queryBus: queryBus, subscriptionId: subscriptionId}
}

func (h *ListProtectedResourcesHandler) Execute(c *gin.Context) (any, error) {
	subscriptionId, ok := h.subscriptionId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SUBSCRIPTION_ID")
	}

	result, err := query.DispatchTyped[[]*protected_resources.ProtectedResource](c.Request.Context(), h.queryBus, prList.NewQuery(subscriptionId))
	if err != nil {
		return nil, err
	}
	primitives := make([]protected_resources.ProtectedResourcePrimitives, 0, len(result))
	for _, item := range result {
		primitives = append(primitives, item.ToPrimitives())
	}
	return primitives, nil
}

func (h *ListProtectedResourcesHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *ListProtectedResourcesHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
