package domain_services_test

import (
	"context"
	"fmt"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func TestScopeFinder_Find(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		repoScope     *scopes.Scope
		repoErr       error
		wantErr       bool
		wantNilResult bool
	}{
		{
			name:      "returns scope when found",
			repoScope: aScope(t),
		},
		{
			name:    "propagates repository error",
			repoErr: fmt.Errorf("db error"),
			wantErr: true,
		},
		{
			name:          "returns not found error when repo returns nil",
			repoScope:     nil,
			wantErr:       true,
			wantNilResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(tt.repoScope, tt.repoErr)

			finder := domain_services.NewScopeFinder(repo)
			got, err := finder.Find(context.Background(), ids2.NewScopeId())

			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				if tt.wantNilResult {
					assert.Nil(t, got)
				}
				return
			}
			assert.Equal(t, tt.repoScope, got)
		})
	}
}

func aScope(t *testing.T) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(ids2.NewProtectedResourceId(), ids2.NewSubscriptionId())
	scope, err := scopes.NewScope(ids2.NewScopeId(), "read:users", ownership, false)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}
