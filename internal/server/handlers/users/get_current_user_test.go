package users_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adminAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	userAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func adminAuthorExtractor(adminId ids.UserId) func(context.Context) (*authors.UserAuthor, bool) {
	return func(_ context.Context) (*authors.UserAuthor, bool) {
		return &authors.UserAuthor{Id: adminId, AdminId: adminId, Role: vo.AdminRole}, true
	}
}

func TestGetCurrentUserHandler(t *testing.T) {
	t.Parallel()
	user, err := userAggregate.FromPrimitives(userAggregate.TenantPrimitives{
		Id: handlerstest.UserId, AdminId: handlerstest.AdminID, Name: "Jane Doe", Email: "jane@example.com",
	})
	require.NoError(t, err)

	adminId := ids.NewUserId()
	admin, err := adminAggregate.NewAdmin(adminId, "Jane Admin", "jane.admin@example.com", false)
	require.NoError(t, err)
	admin.PullDomainEvents()

	tests := []struct {
		name      string
		authorFn  func(context.Context) (*authors.UserAuthor, bool)
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
				Role: "tenant", AdminId: handlerstest.AdminID,
			},
		},
		{
			name: "success_returns_admin_profile", authorFn: adminAuthorExtractor(adminId), busResult: admin,
			want: users.CurrentUserResponse{
				Id: adminId.String(), Email: "jane.admin@example.com", DisplayName: "Jane Admin",
				Role: "admin", AdminId: adminId.String(),
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
