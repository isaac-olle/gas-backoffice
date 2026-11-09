package domain_services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDuplicateNameChecker(t *testing.T) {
	t.Parallel()

	t.Run("returns nil when no scope matches", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "read:users", "pr-id")
		require.NoError(t, err)
	})

	t.Run("returns conflict error when a scope already matches", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*scopes.Scope{aScope(t)}, nil)

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "read:users", "pr-id")

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, scopes.ErrScopeDuplicatedName, appErr.Code)
	})

	t.Run("returns internal error when repository query fails", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "read:users", "pr-id")

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, app_errors.ErrQueryingDatabase, appErr.Code)
	})
}
