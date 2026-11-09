package error_logs_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errorLogAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestListErrorLogsHandler(t *testing.T) {
	t.Parallel()

	t.Run("bus_error_returns_error", func(t *testing.T) {
		t.Parallel()
		h := error_logs.NewListErrorLogsHandler(handlerstest.NewQueryBus(nil, errors.New("db")))
		_, c := handlerstest.NewCtx(http.MethodGet, nil)
		_, err := h.Execute(c)
		assert.Error(t, err)
	})

	t.Run("success_returns_200", func(t *testing.T) {
		t.Parallel()
		logs := []errorLogAggregate.ErrorLog{{ErrorId: "e1", Code: "X"}}
		h := error_logs.NewListErrorLogsHandler(handlerstest.NewQueryBus(logs, nil))
		_, c := handlerstest.NewCtx(http.MethodGet, nil)
		result, err := h.Execute(c)
		require.NoError(t, err)
		assert.Equal(t, logs, result)
		require.NoError(t, h.Success(c, result))
		assert.Equal(t, http.StatusOK, c.Writer.Status())
	})
}

func TestListErrorLogsHandler_Default(t *testing.T) {
	t.Parallel()
	h := error_logs.NewListErrorLogsHandler(handlerstest.NewQueryBus(nil, nil))
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
