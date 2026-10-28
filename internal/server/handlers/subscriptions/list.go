package subscriptions

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	backofficeSubscriptions "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	listSubscriptions "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/list"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

type ListHandler struct {
	queryBus query.Bus
	author   func(context.Context) (*authors.UserAuthor, bool)
}

func NewListHandler(queryBus query.Bus, author func(context.Context) (*authors.UserAuthor, bool)) *ListHandler {
	return &ListHandler{queryBus: queryBus, author: author}
}

func (h *ListHandler) Execute(c *gin.Context) (any, error) {
	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	return query.DispatchTyped[[]backofficeSubscriptions.Summary](c.Request.Context(), h.queryBus,
		listSubscriptions.NewQuery(author.Role.IsSystemAdministrator(), author.SubscriptionIds))
}

func (h *ListHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *ListHandler) Default(c *gin.Context) { c.Status(http.StatusInternalServerError) }
