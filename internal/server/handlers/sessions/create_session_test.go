package sessions_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/sessions"
)

func TestCreateSessionHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		sessionIdFn func(context.Context) (ids2.SessionId, bool)
		userIdFn    func(context.Context) (ids2.UserId, bool)
		busErr      error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "missing_session_id_returns_error",
			sessionIdFn: handlerstest.MissingSessionIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_user_id_returns_error",
			sessionIdFn: handlerstest.ValidSessionIdExtractor(),
			userIdFn:    handlerstest.MissingUserIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "bus_error_returns_error",
			sessionIdFn: handlerstest.ValidSessionIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			busErr:      errors.New("bus error"),
			wantErr:     true,
		},
		{
			name:        "success_returns_201",
			sessionIdFn: handlerstest.ValidSessionIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			wantErr:     false,
			wantStatus:  http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := sessions.NewCreateSessionHandler(bus, tt.sessionIdFn, tt.userIdFn)

			_, c := handlerstest.NewCtx(http.MethodPut, nil)
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

func TestCreateSessionHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := sessions.NewCreateSessionHandler(bus, handlerstest.ValidSessionIdExtractor(), handlerstest.ValidUserIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPut, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
