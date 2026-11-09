package users_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
)

func TestSoftDeleteUserHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		userIdFn   func(context.Context) (ids.UserId, bool)
		authorFn   func(context.Context) (*vo.UserAuthor, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:     "missing_user_id_returns_error",
			userIdFn: handlerstest.MissingUserIdExtractor(),
			authorFn: handlerstest.ValidAuthorExtractor(),
			wantErr:  true,
		},
		{
			name:     "missing_author_returns_error",
			userIdFn: handlerstest.ValidUserIdExtractor(),
			authorFn: handlerstest.MissingAuthorExtractor(),
			wantErr:  true,
		},
		{
			name:     "bus_error_returns_error",
			userIdFn: handlerstest.ValidUserIdExtractor(),
			authorFn: handlerstest.ValidAuthorExtractor(),
			busErr:   errors.New("bus error"),
			wantErr:  true,
		},
		{
			name:       "success_returns_204",
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			authorFn:   handlerstest.ValidAuthorExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := users.NewSoftDeleteUserHandler(bus, tt.userIdFn, tt.authorFn)

			_, c := handlerstest.NewCtx(http.MethodDelete, nil)
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

func TestSoftDeleteUserHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := users.NewSoftDeleteUserHandler(bus, handlerstest.ValidUserIdExtractor(), handlerstest.ValidAuthorExtractor())
	_, c := handlerstest.NewCtx(http.MethodDelete, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
