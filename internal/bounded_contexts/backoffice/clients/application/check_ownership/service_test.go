package check_ownership_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func aClient(t *testing.T, userId ids2.UserId) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(clientId, "my-client", clients.NewOwnership(ids2.NewUserId(), userId), nil, nil, false, false)
	require.NoError(t, err)
	c.PullDomainEvents()
	return c, clientId
}

func newService(repo *clientmocks.MockRepository) *check_ownership.Service {
	return check_ownership.NewService(domain_services.NewClientFinder(repo))
}

func TestCheck(t *testing.T) {
	t.Parallel()

	t.Run("returns ownership primitives when the client belongs to the user", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		userId := ids2.NewUserId()
		c, clientId := aClient(t, userId)
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)

		result, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(clientId.String(), ids2.NewUserId().String(), userId.String()))
		require.NoError(t, err)
		p, ok := result.(clients.OwnershipPrimitives)
		require.True(t, ok)
		assert.Equal(t, userId.String(), p.UserId)
	})

	t.Run("returns ownership primitives when the caller is the owning admin", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		adminId := ids2.NewUserId()
		clientId := ids2.NewClientId()
		c, err := clients.NewClient(clientId, "tenant-client", clients.NewOwnership(adminId, ids2.NewUserId()), nil, nil, false, false)
		require.NoError(t, err)
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)

		result, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(clientId.String(), adminId.String(), adminId.String()))
		require.NoError(t, err)
		p, ok := result.(clients.OwnershipPrimitives)
		require.True(t, ok)
		assert.Equal(t, adminId.String(), p.AdminId)
	})

	t.Run("rejects invalid admin id", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(ids2.NewClientId().String(), "bad-id", ids2.NewUserId().String()))
		assertErrCode(t, err, clients.ErrClientNotFound)
	})

	t.Run("rejects invalid client id", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery("bad-id", ids2.NewUserId().String(), ids2.NewUserId().String()))
		assertErrCode(t, err, clients.ErrClientNotFound)
	})

	t.Run("rejects invalid user id", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(ids2.NewClientId().String(), ids2.NewUserId().String(), "bad-id"))
		assertErrCode(t, err, clients.ErrClientNotFound)
	})

	t.Run("returns not found when client does not exist", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)

		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(ids2.NewClientId().String(), ids2.NewUserId().String(), ids2.NewUserId().String()))
		assertErrCode(t, err, clients.ErrClientNotFound)
	})

	t.Run("returns forbidden when client belongs to another user", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		c, clientId := aClient(t, ids2.NewUserId())
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)

		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(clientId.String(), ids2.NewUserId().String(), ids2.NewUserId().String()))
		assertErrCode(t, err, clients.ErrInvalidOwnership)
	})
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
