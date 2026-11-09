package domain_services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDuplicateNameChecker(t *testing.T) {
	t.Parallel()

	t.Run("returns nil when no client matches", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "My Client", "user-id")
		require.NoError(t, err)
	})

	t.Run("returns conflict error when a client already matches", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*clients.Client{aClient(t)}, nil)

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "My Client", "user-id")

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, clients.ErrClientDuplicateName, appErr.Code)
	})

	t.Run("returns internal error when repository query fails", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		checker := domain_services.NewDuplicateNameChecker(repo)
		err := checker.Check(context.Background(), "My Client", "user-id")

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, app_errors.ErrQueryingDatabase, appErr.Code)
	})
}
