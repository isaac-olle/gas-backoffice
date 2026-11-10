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

	tests := []struct {
		name    string
		busLogs []errorLogAggregate.ErrorLog
		busErr  error
		wantErr bool
	}{
		{
			name:    "bus_error_returns_error",
			busErr:  errors.New("db"),
			wantErr: true,
		},
		{
			name:    "success_returns_200",
			busLogs: []errorLogAggregate.ErrorLog{{ErrorId: "e1", Code: "X"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := error_logs.NewListErrorLogsHandler(handlerstest.NewQueryBus(tt.busLogs, tt.busErr))
			_, c := handlerstest.NewCtx(http.MethodGet, nil)
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.busLogs, result)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusOK, c.Writer.Status())
		})
	}
}

func TestListErrorLogsHandler_Default(t *testing.T) {
	t.Parallel()
	h := error_logs.NewListErrorLogsHandler(handlerstest.NewQueryBus(nil, nil))
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
