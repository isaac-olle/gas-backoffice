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

func aPR(t *testing.T, subscriptionId ids2.SubscriptionId) (*protected_resources.ProtectedResource, ids2.ProtectedResourceId) {
	t.Helper()
	prId := ids2.NewProtectedResourceId()
	pr, err := protected_resources.NewProtectedResource(prId, "my-resource", protected_resources.NewOwnership(subscriptionId), false)
	require.NoError(t, err)
	pr.PullDomainEvents()
	return pr, prId
}

func newService(repo *prmocks.MockRepository) *check_ownership.Service {
	return check_ownership.NewService(domain_services.NewProtectedResourceFinder(repo))
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(t *testing.T, repo *prmocks.MockRepository) (ctx context.Context, qry *check_ownership.Query, wantResult string)
		wantErr     bool
		wantErrCode app_errors.ErrorCode
	}{
		{
			name: "returns the subscription id when it is in the allowed set",
			setup: func(t *testing.T, repo *prmocks.MockRepository) (context.Context, *check_ownership.Query, string) {
				subscriptionId := ids2.NewSubscriptionId()
				pr, prId := aPR(t, subscriptionId)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				return context.Background(),
					check_ownership.NewQueryForSubscriptions(prId.String(), []string{subscriptionId.String()}, false),
					subscriptionId.String()
			},
		},
		{
			name: "allows any subscription when allowAll is set",
			setup: func(t *testing.T, repo *prmocks.MockRepository) (context.Context, *check_ownership.Query, string) {
				pr, prId := aPR(t, ids2.NewSubscriptionId())
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				return context.Background(), check_ownership.NewQueryForResource(prId.String()), pr.ToPrimitives().SubscriptionId
			},
		},
		{
			name: "falls back to the subscription id in the context when none is given",
			setup: func(t *testing.T, repo *prmocks.MockRepository) (context.Context, *check_ownership.Query, string) {
				subscriptionId := ids2.NewSubscriptionId()
				pr, prId := aPR(t, subscriptionId)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				ctx := ids2.WithSubscriptionId(context.Background(), subscriptionId)
				return ctx, check_ownership.NewQuery(prId.String()), subscriptionId.String()
			},
		},
		{
			name: "rejects invalid protected resource id",
			setup: func(t *testing.T, repo *prmocks.MockRepository) (context.Context, *check_ownership.Query, string) {
				return context.Background(), check_ownership.NewQuery("bad-id"), ""
			},
			wantErr:     true,
			wantErrCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			name: "rejects invalid allowed subscription id format",
			setup: func(t *testing.T, repo *prmocks.MockRepository) (context.Context, *check_ownership.Query, string) {
				pr, prId := aPR(t, ids2.NewSubscriptionId())
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				return context.Background(), check_ownership.NewQueryForSubscriptions(prId.String(), []string{"bad-id"}, false), ""
			},
			wantErr:     true,
			wantErrCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			name: "returns not found when protected resource does not exist",
			setup: func(t *testing.T, repo *prmocks.MockRepository) (context.Context, *check_ownership.Query, string) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
				return context.Background(),
					check_ownership.NewQueryForSubscriptions(ids2.NewProtectedResourceId().String(), []string{ids2.NewSubscriptionId().String()}, false),
					""
			},
			wantErr:     true,
			wantErrCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			name: "returns forbidden when protected resource belongs to another subscription",
			setup: func(t *testing.T, repo *prmocks.MockRepository) (context.Context, *check_ownership.Query, string) {
				pr, prId := aPR(t, ids2.NewSubscriptionId())
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				return context.Background(), check_ownership.NewQueryForSubscriptions(prId.String(), []string{ids2.NewSubscriptionId().String()}, false), ""
			},
			wantErr:     true,
			wantErrCode: protected_resources.ErrInvalidOwnership,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := prmocks.NewMockRepository(t)
			ctx, qry, wantResult := tt.setup(t, repo)

			result, err := newService(repo).Check(ctx, qry)

			if tt.wantErr {
				assertErrCode(t, err, tt.wantErrCode)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, wantResult, result)
		})
	}
}
