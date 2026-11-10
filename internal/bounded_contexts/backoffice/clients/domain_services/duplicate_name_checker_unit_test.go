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

	tests := []struct {
		name        string
		repoResult  []*clients.Client
		repoErr     error
		wantErr     bool
		wantErrCode app_errors.ErrorCode
	}{
		{
			name: "returns nil when no client matches",
		},
		{
			name:        "returns conflict error when a client already matches",
			repoResult:  []*clients.Client{aClient(t)},
			wantErr:     true,
			wantErrCode: clients.ErrClientDuplicateName,
		},
		{
			name:        "returns internal error when repository query fails",
			repoErr:     fmt.Errorf("db error"),
			wantErr:     true,
			wantErrCode: app_errors.ErrQueryingDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := clientmocks.NewMockRepository(t)
			repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(tt.repoResult, tt.repoErr)

			checker := domain_services.NewDuplicateNameChecker(repo)
			err := checker.Check(context.Background(), "My Client", "user-id")

			if tt.wantErr {
				var appErr *app_errors.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tt.wantErrCode, appErr.Code)
				return
			}
			require.NoError(t, err)
		})
	}
}
