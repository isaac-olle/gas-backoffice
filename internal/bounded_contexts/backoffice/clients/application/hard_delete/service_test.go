package hard_delete_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes client with correct id from repository", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		clientId := ids.NewClientId()
		userId := ids.NewUserId()
		c, _ := clients.NewClient(clientId, "valid-client", clients.NewOwnership(ids.NewUserId(), userId), nil, nil, false, false)
		c.PullDomainEvents()
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)
		repo.EXPECT().HardDelete(mock.Anything, clientId).Return(nil)

		err := newService(t, repo).HardDelete(context.Background(), hard_delete.NewCommand(clientId.String(), ids.NewUserId().String(), userId.String()))
		require.NoError(t, err)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *hard_delete.Command
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid client id format",
			cmd:      func(t *testing.T) *hard_delete.Command { return aCommand(t, withClientId("bad-id")) },
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc:     "rejects invalid user id format",
			cmd:      func(t *testing.T) *hard_delete.Command { return aCommand(t, withUserId("bad-id")) },
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns not found when client does not exist",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns forbidden when caller does not own the client",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids.NewClientId(), "valid-client", clients.NewOwnership(ids.NewUserId(), ids.NewUserId()), nil, nil, false, false)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrInvalidOwnership,
		},
		{
			desc: "returns forbidden when client is system resource",
			cmd: func(t *testing.T) *hard_delete.Command {
				return aCommand(t, withUserId(systemResourceUserId.String()))
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.FromPrimitives(clients.ClientPrimitives{
					Id:       ids.NewClientId().String(),
					Name:     "system-client",
					AdminId:  ids.NewUserId().String(),
					UserId:   systemResourceUserId.String(),
					IsSystem: true,
				})
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrSystemResource,
		},
		{
			desc: "wraps repository hard delete failure",
			cmd: func(t *testing.T) *hard_delete.Command {
				return aCommand(t, withUserId(saveFailureUserId.String()))
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids.NewClientId(), "valid-client", clients.NewOwnership(ids.NewUserId(), saveFailureUserId), nil, nil, false, false)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().HardDelete(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := clientmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			assertErrCode(t, newService(t, repo).HardDelete(context.Background(), tc.cmd(t)), tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

var (
	systemResourceUserId = ids.NewUserId()
	saveFailureUserId    = ids.NewUserId()
)

func newService(t *testing.T, repo *clientmocks.MockRepository) *hard_delete.Service {
	t.Helper()
	return hard_delete.NewService(domain_services.NewClientFinder(repo), repo)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, code, appErr.Code)
}
