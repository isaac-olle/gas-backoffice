package get_test

import (
	"context"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"

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
		subscriptionId := ids2.NewSubscriptionId()
		c, clientId := aClientWithId(t, "my-get-client", subscriptionId)
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)
		author, err := vo.NewUserAuthorWithMemberships(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole,
			[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
		require.NoError(t, err)

		result, err := newService(t, repo).Get(context.Background(), get.NewQuery(clientId.String(), author))
		require.NoError(t, err)
		c, ok := result.(*clients.Client)
		require.True(t, ok)
		p := c.ToPrimitives()
		assert.Equal(t, "my-get-client", p.Name)
		assert.Equal(t, clientId.String(), p.Id)
	})

	t.Run("returns the client when the caller is a system administrator", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		clientId := ids2.NewClientId()
		c, err := clients.NewClient(clientId, "tenant-client", clients.NewOwnership(ids2.NewSubscriptionId()), nil, nil, false, false)
		require.NoError(t, err)
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)
		adminId := ids2.NewUserId()
		author, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
		require.NoError(t, err)

		_, err = newService(t, repo).Get(context.Background(), get.NewQuery(clientId.String(), author))
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
			desc:  "returns not found when client does not exist",
			query: func(t *testing.T) *get.Query { return aQuery(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns forbidden when author lacks subscription access",
			query: func(t *testing.T) *get.Query {
				unauthorized, err := vo.NewUserAuthorFromString(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole)
				require.NoError(t, err)
				return aQuery(t, withAuthor(unauthorized))
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClientWithId(t, "other-client", ids2.NewSubscriptionId())
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
	return get.NewService(finder, sharedDomainServices.NewSubscriptionAuthorizer())
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
