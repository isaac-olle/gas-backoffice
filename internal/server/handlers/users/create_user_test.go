package users_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestCreateUserHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       io.Reader
		userIdFn   func(context.Context) (ids2.UserId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:     "invalid_body_returns_error",
			body:     bytes.NewBufferString(`{}`),
			userIdFn: handlerstest.ValidUserIdExtractor(),
			wantErr:  true,
		},
		{
			name:     "missing_user_id_returns_error",
			body:     handlerstest.JsonBody(map[string]string{"name": "My User"}),
			userIdFn: handlerstest.MissingUserIdExtractor(),
			wantErr:  true,
		},
		{
			name:     "bus_error_returns_error",
			body:     handlerstest.JsonBody(map[string]string{"name": "My User"}),
			userIdFn: handlerstest.ValidUserIdExtractor(),
			busErr:   errors.New("bus error"),
			wantErr:  true,
		},
		{
			name:       "success_returns_201",
			body:       handlerstest.JsonBody(map[string]string{"name": "My User"}),
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := users.NewCreateUserHandler(bus, tt.userIdFn, handlerstest.ValidUserIdExtractor())

			_, c := handlerstest.NewCtx(http.MethodPost, tt.body)
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

func TestCreateUserHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := users.NewCreateUserHandler(bus, handlerstest.ValidUserIdExtractor(), handlerstest.ValidUserIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
