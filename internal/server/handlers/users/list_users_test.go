package users_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestListUsersHandler(t *testing.T) {
	t.Parallel()
	item, err := aggregate.FromPrimitives(aggregate.UserPrimitives{Id: handlerstest.UserId, Name: "Jane", Email: "jane@example.com", SystemRole: "subscription_user"})
	require.NoError(t, err)

	tests := []struct {
		name      string
		userIdFn  func(context.Context) (ids.UserId, bool)
		busResult any
		busErr    error
		wantErr   bool
		wantLen   int
	}{
		{name: "missing_user_id_returns_error", userIdFn: handlerstest.MissingUserIdExtractor(), wantErr: true},
		{name: "bus_error_returns_error", userIdFn: handlerstest.ValidUserIdExtractor(), busErr: errors.New("db"), wantErr: true},
		{name: "empty_list_returns_200", userIdFn: handlerstest.ValidUserIdExtractor(), busResult: []*aggregate.User{}, wantLen: 0},
		{name: "success_returns_200", userIdFn: handlerstest.ValidUserIdExtractor(), busResult: []*aggregate.User{item}, wantLen: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := users.NewListUsersHandler(bus, tt.userIdFn)

			_, c := handlerstest.NewCtx(http.MethodGet, nil)
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			primitives, ok := result.([]aggregate.UserPrimitives)
			require.True(t, ok)
			assert.Len(t, primitives, tt.wantLen)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusOK, c.Writer.Status())
		})
	}
}

func TestListUsersHandler_Default(t *testing.T) {
	t.Parallel()
	h := users.NewListUsersHandler(handlerstest.NewQueryBus(nil, nil), handlerstest.ValidUserIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
