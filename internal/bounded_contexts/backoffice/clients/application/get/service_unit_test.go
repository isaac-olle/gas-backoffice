package get_test

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
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
		userId := ids2.NewUserId()
		c, clientId := aClientWithId(t, "my-get-client", userId)
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)

		result, err := newService(t, repo).Get(context.Background(), get.NewQuery(clientId.String(), ids2.NewUserId().String(), userId.String()))
		require.NoError(t, err)
		c, ok := result.(*clients.Client)
		require.True(t, ok)
		p := c.ToPrimitives()
		assert.Equal(t, "my-get-client", p.Name)
		assert.Equal(t, clientId.String(), p.Id)
	})

	t.Run("returns the client when the caller is its owning admin", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		adminId := ids2.NewUserId()
		clientId := ids2.NewClientId()
		c, err := clients.NewClient(clientId, "tenant-client", clients.NewOwnership(adminId, ids2.NewUserId()), nil, nil, false, false)
		require.NoError(t, err)
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)

		_, err = newService(t, repo).Get(context.Background(), get.NewQuery(clientId.String(), adminId.String(), adminId.String()))
		require.NoError(t, err)
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
			desc:     "rejects invalid user id format",
			query:    func(t *testing.T) *get.Query { return aQuery(t, withUserId("bad-id")) },
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
		{
			desc: "returns forbidden when ownership does not match",
			query: func(t *testing.T) *get.Query {
				return aQuery(t)
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClientWithId(t, "other-client", ids2.NewUserId())
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrInvalidOwnership,
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
	finder := domain_services.NewClientFinder(repo)
	return get.NewService(finder)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
