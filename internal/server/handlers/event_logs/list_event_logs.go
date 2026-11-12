package event_logs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/event_logs"
	eventLogList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/event_logs/application/list"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type ListEventLogsHandler struct {
	queryBus query.Bus
}

func NewListEventLogsHandler(queryBus query.Bus) *ListEventLogsHandler {
	return &ListEventLogsHandler{queryBus: queryBus}
}

func (h *ListEventLogsHandler) Execute(c *gin.Context) (any, error) {
	return query.DispatchTyped[[]event_logs.EventLog](c.Request.Context(), h.queryBus, eventLogList.NewQuery())
}

func (h *ListEventLogsHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *ListEventLogsHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
