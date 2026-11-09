package check_ownership_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func aPR(t *testing.T, userId ids2.UserId) (*protected_resources.ProtectedResource, ids2.ProtectedResourceId) {
	t.Helper()
	prId := ids2.NewProtectedResourceId()
	pr, err := protected_resources.NewProtectedResource(prId, "my-resource", protected_resources.NewOwnership(ids2.NewUserId(), userId), false)
	require.NoError(t, err)
	pr.PullDomainEvents()
	return pr, prId
}

func newService(repo *prmocks.MockRepository) *check_ownership.Service {
	return check_ownership.NewService(domain_services.NewProtectedResourceFinder(repo))
}

func TestCheck(t *testing.T) {
	t.Parallel()

	t.Run("returns ownership primitives when the protected resource belongs to the user", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		userId := ids2.NewUserId()
		pr, prId := aPR(t, userId)
		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)

		result, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(prId.String(), ids2.NewUserId().String(), userId.String()))
		require.NoError(t, err)
		p, ok := result.(protected_resources.OwnershipPrimitives)
		require.True(t, ok)
		assert.Equal(t, userId.String(), p.UserId)
	})

	t.Run("rejects invalid protected resource id", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery("bad-id", ids2.NewUserId().String(), ids2.NewUserId().String()))
		assertErrCode(t, err, protected_resources.ErrProtectedResourceNotFound)
	})

	t.Run("rejects invalid user id", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(ids2.NewProtectedResourceId().String(), ids2.NewUserId().String(), "bad-id"))
		assertErrCode(t, err, protected_resources.ErrInvalidOwnership)
	})

	t.Run("returns not found when protected resource does not exist", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)

		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(ids2.NewProtectedResourceId().String(), ids2.NewUserId().String(), ids2.NewUserId().String()))
		assertErrCode(t, err, protected_resources.ErrProtectedResourceNotFound)
	})

	t.Run("returns forbidden when protected resource belongs to another user", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		pr, prId := aPR(t, ids2.NewUserId())
		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)

		_, err := newService(repo).Check(context.Background(), check_ownership.NewQuery(prId.String(), ids2.NewUserId().String(), ids2.NewUserId().String()))
		assertErrCode(t, err, protected_resources.ErrInvalidOwnership)
	})
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
