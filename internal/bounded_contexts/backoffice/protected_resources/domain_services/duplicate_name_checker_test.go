package domain_services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDuplicateNameChecker(t *testing.T) {
	t.Parallel()

	t.Run("returns nil when no protected resource matches", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "My Resource", "user-id")
		require.NoError(t, err)
	})

	t.Run("returns conflict error when a protected resource already matches", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*protected_resources.ProtectedResource{aPR(t)}, nil)

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "My Resource", "user-id")

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, protected_resources.ErrProtectedResourceDuplicateName, appErr.Code)
	})

	t.Run("returns internal error when repository query fails", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "My Resource", "user-id")

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, app_errors.ErrQueryingDatabase, appErr.Code)
	})
}
