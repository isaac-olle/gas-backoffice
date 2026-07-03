package hard_delete_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/hard_delete"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes client with correct id from repository", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		clientId := ids.NewClientId()
		repo.EXPECT().HardDelete(mock.Anything, clientId).Return(nil)

		err := hard_delete.NewService(repo).HardDelete(context.Background(), hard_delete.NewCommand(clientId.String()))
		require.NoError(t, err)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *hard_delete.Command
		setupRepo func(repo *clientmocks.MockRepository) // nil = HardDelete no s'hauria de cridar
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid client id format",
			cmd:      func(t *testing.T) *hard_delete.Command { return aCommand(t, withClientId("bad-id")) },
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "wraps repository hard delete failure",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().
					HardDelete(mock.Anything, mock.AnythingOfType("ids.ClientId")).
					Return(fmt.Errorf("db error"))
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
			assertErrCode(t, hard_delete.NewService(repo).HardDelete(context.Background(), tc.cmd(t)), tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
