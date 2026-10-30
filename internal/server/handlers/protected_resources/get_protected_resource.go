package protected_resources

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	prGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type GetProtectedResourceHandler struct {
	queryBus query.Bus
	prId     func(context.Context) (ids.ProtectedResourceId, bool)
	author   func(context.Context) (*vo.UserAuthor, bool)
}

func NewGetProtectedResourceHandler(
	queryBus query.Bus,
	prId func(context.Context) (ids.ProtectedResourceId, bool),
	author func(context.Context) (*vo.UserAuthor, bool),
) *GetProtectedResourceHandler {
	return &GetProtectedResourceHandler{queryBus: queryBus, prId: prId, author: author}
}

func (h *GetProtectedResourceHandler) Execute(c *gin.Context) (any, error) {
	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
	}

	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	pr, err := query.DispatchTyped[*protected_resources.ProtectedResource](
		c.Request.Context(), h.queryBus, prGet.NewQuery(prId, author),
	)
	if err != nil {
		return nil, err
	}
	return pr.ToPrimitives(), nil
}

func (h *GetProtectedResourceHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetProtectedResourceHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
