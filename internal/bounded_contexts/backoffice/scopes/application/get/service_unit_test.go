package get_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(repo *scopemocks.MockRepository) *get.Service {
	return get.NewService(domain_services.NewScopeFinder(repo), sharedDomainServices.NewSubscriptionAuthorizer())
}

func authorFor(t *testing.T, subscriptionId ids.SubscriptionId) *vo.UserAuthor {
	t.Helper()
	author, err := vo.NewUserAuthorWithMemberships(ids.NewUserId().String(), ids.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	return author
}

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		sameSubscripion bool
		repoErr         error
		wantErr         bool
		wantErrCode     app_errors.ErrorCode
	}{
		{
			name:            "returns scope when found",
			sameSubscripion: true,
		},
		{
			name:            "returns forbidden error when author lacks subscription access",
			sameSubscripion: false,
			wantErr:         true,
			wantErrCode:     scopes.ErrInvalidOwnership,
		},
		{
			name:    "returns error when repo errors",
			repoErr: fmt.Errorf("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			scopeId := ids.NewScopeId()
			prId := ids.NewProtectedResourceId()
			scopeSubscriptionId := ids.NewSubscriptionId()
			var scope *scopes.Scope
			if tt.repoErr == nil {
				scope = aScope(t, scopeId, prId, scopeSubscriptionId)
			}

			repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, tt.repoErr)

			authorSubscriptionId := scopeSubscriptionId
			if !tt.sameSubscripion {
				authorSubscriptionId = ids.NewSubscriptionId()
			}

			svc := newService(repo)
			result, err := svc.Get(context.Background(), get.NewQuery(scopeId, authorFor(t, authorSubscriptionId)))

			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				if tt.wantErrCode != "" {
					var appErr *app_errors.AppError
					require.ErrorAs(t, err, &appErr)
					assert.Equal(t, tt.wantErrCode, appErr.Code)
				}
				return
			}
			s, ok := result.(*scopes.Scope)
			require.True(t, ok)
			assert.Equal(t, scopeId.String(), s.ToPrimitives().Id)
		})
	}
}
