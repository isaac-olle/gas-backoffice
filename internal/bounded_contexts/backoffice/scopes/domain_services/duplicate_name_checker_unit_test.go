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

	tests := []struct {
		name        string
		repoResult  []*scopes.Scope
		repoErr     error
		wantErr     bool
		wantErrCode app_errors.ErrorCode
	}{
		{
			name: "returns nil when no scope matches",
		},
		{
			name:        "returns conflict error when a scope already matches",
			repoResult:  []*scopes.Scope{aScope(t)},
			wantErr:     true,
			wantErrCode: scopes.ErrScopeDuplicatedName,
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
			repo := scopemocks.NewMockRepository(t)
			repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(tt.repoResult, tt.repoErr)

			checker := domain_services.NewDuplicateNameChecker(repo)
			err := checker.Check(context.Background(), "read:users", "pr-id")

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
