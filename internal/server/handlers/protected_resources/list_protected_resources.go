package protected_resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	prList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/list"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

type ListProtectedResourcesHandler struct {
	queryBus query.Bus
}

func NewListProtectedResourcesHandler(queryBus query.Bus) *ListProtectedResourcesHandler {
	return &ListProtectedResourcesHandler{queryBus: queryBus}
}

func (h *ListProtectedResourcesHandler) Execute(c *gin.Context) (any, error) {
	author, ok := vo.AuthorFromContext(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	result, err := query.DispatchTyped[[]*protected_resources.ProtectedResource](c.Request.Context(), h.queryBus,
		prList.NewQuery(author))
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
