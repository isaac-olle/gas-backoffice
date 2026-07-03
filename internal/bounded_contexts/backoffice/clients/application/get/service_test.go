package get_test

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns client primitives when found", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		c, clientId := aClientWithId(t, "my-get-client")
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)

		result, err := newService(t, repo).Get(context.Background(), get.NewQuery(clientId.String()))
		require.NoError(t, err)
		p, ok := result.(clients.ClientPrimitives)
		require.True(t, ok)
		assert.Equal(t, "my-get-client", p.Name)
		assert.Equal(t, clientId.String(), p.Id)
	})

	errCases := []struct {
		desc      string
		query     func(t *testing.T) *get.Query
		setupRepo func(repo *clientmocks.MockRepository) // nil = SearchById no s'hauria de cridar
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid client id format",
			query:    func(t *testing.T) *get.Query { return aQuery(t, withClientId("bad-id")) },
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc:  "returns not found when client does not exist",
			query: func(t *testing.T) *get.Query { return aQuery(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
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
			_, err := newService(t, repo).Get(context.Background(), tc.query(t))
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newService(t *testing.T, repo *clientmocks.MockRepository) *get.Service {
	t.Helper()
	return get.NewService(domain_services.NewClientFinder(repo))
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
