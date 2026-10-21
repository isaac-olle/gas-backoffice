package users_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
)

func TestGetUserHandler(t *testing.T) {
	t.Parallel()
	validUser, _ := userAggregate.FromPrimitives(userAggregate.TenantPrimitives{
		Id:      handlerstest.UserId,
		AdminId: handlerstest.UserId,
		Name:    "Acme",
	})

	tests := []struct {
		name       string
		userIdFn   func(context.Context) (ids.UserId, bool)
		authorFn   func(context.Context) (*authors.UserAuthor, bool)
		busResult  any
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
			busErr:   errors.New("not found"),
			wantErr:  true,
		},
		{
			name:       "success_returns_200",
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			authorFn:   handlerstest.ValidAuthorExtractor(),
			busResult:  validUser,
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := users.NewGetUserHandler(bus, tt.userIdFn, tt.authorFn)

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

func TestGetUserHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := users.NewGetUserHandler(bus, handlerstest.ValidUserIdExtractor(), handlerstest.ValidAuthorExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
