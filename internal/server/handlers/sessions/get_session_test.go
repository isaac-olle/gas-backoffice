package sessions_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/sessions"
)

func TestGetSessionHandler(t *testing.T) {
	t.Parallel()
	validSession, _ := authAggregate.FromPrimitives(authAggregate.SessionPrimitives{
		Id:     handlerstest.SessionID,
		UserId: handlerstest.UserId,
	})

	tests := []struct {
		name        string
		sessionIdFn func(context.Context) (ids2.SessionId, bool)
		busResult   any
		busErr      error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "missing_session_id_returns_error",
			sessionIdFn: handlerstest.MissingSessionIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "bus_error_returns_error",
			sessionIdFn: handlerstest.ValidSessionIdExtractor(),
			busErr:      errors.New("not found"),
			wantErr:     true,
		},
		{
			name:        "success_returns_200",
			sessionIdFn: handlerstest.ValidSessionIdExtractor(),
			busResult:   validSession,
			wantErr:     false,
			wantStatus:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := sessions.NewGetSessionHandler(bus, tt.sessionIdFn)

			_, c := handlerstest.NewCtx(http.MethodGet, nil)
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, tt.wantStatus, c.Writer.Status())
		})
	}
}

func TestGetSessionHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := sessions.NewGetSessionHandler(bus, handlerstest.ValidSessionIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
