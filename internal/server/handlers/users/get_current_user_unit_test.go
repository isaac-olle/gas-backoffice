package users_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func adminAuthorExtractor(adminId ids.UserId) func(context.Context) (*vo.UserAuthor, bool) {
	return func(_ context.Context) (*vo.UserAuthor, bool) {
		return &vo.UserAuthor{Id: adminId, AdminId: adminId, Role: vo.SystemAdministratorRole}, true
	}
}

func TestGetCurrentUserHandler(t *testing.T) {
	t.Parallel()
	user, err := userAggregate.FromPrimitives(userAggregate.UserPrimitives{
		Id: handlerstest.UserId, Name: "Jane Doe", Email: "jane@example.com", SystemRole: vo.SubscriptionUserRole.String(),
	})
	require.NoError(t, err)

	adminId := ids.NewUserId()
	admin, err := userAggregate.NewSystemAdministrator(adminId, "Jane Admin", "jane.admin@example.com")
	require.NoError(t, err)
	admin.PullDomainEvents()

	tests := []struct {
		name      string
		authorFn  func(context.Context) (*vo.UserAuthor, bool)
		busResult any
		busErr    error
		wantErr   bool
		want      users.CurrentUserResponse
	}{
		{name: "missing_author_returns_error", authorFn: handlerstest.MissingAuthorExtractor(), wantErr: true},
		{name: "bus_error_returns_error", authorFn: handlerstest.ValidAuthorExtractor(), busErr: errors.New("not found"), wantErr: true},
		{
			name: "success_returns_tenant_profile", authorFn: handlerstest.ValidAuthorExtractor(), busResult: user,
			want: users.CurrentUserResponse{
				Id: handlerstest.UserId, Email: "jane@example.com", DisplayName: "Jane Doe",
				Role: vo.SubscriptionUserRole.String(),
			},
		},
		{
			name: "success_returns_admin_profile", authorFn: adminAuthorExtractor(adminId), busResult: admin,
			want: users.CurrentUserResponse{
				Id: adminId.String(), Email: "jane.admin@example.com", DisplayName: "Jane Admin",
				Role: vo.SystemAdministratorRole.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := users.NewGetCurrentUserHandler(bus, tt.authorFn)

			_, c := handlerstest.NewCtx(http.MethodGet, nil)
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusOK, c.Writer.Status())
		})
	}
}

func TestGetCurrentUserHandler_Default(t *testing.T) {
	t.Parallel()
	h := users.NewGetCurrentUserHandler(handlerstest.NewQueryBus(nil, nil), handlerstest.ValidAuthorExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
