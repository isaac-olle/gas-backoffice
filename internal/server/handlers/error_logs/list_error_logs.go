package error_logs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	errorLogList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/application/list"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type ListErrorLogsHandler struct {
	queryBus query.Bus
}

func NewListErrorLogsHandler(queryBus query.Bus) *ListErrorLogsHandler {
	return &ListErrorLogsHandler{queryBus: queryBus}
}

func (h *ListErrorLogsHandler) Execute(c *gin.Context) (any, error) {
	return query.DispatchTyped[[]error_logs.ErrorLog](c.Request.Context(), h.queryBus, errorLogList.NewQuery())
}

func (h *ListErrorLogsHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *ListErrorLogsHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
