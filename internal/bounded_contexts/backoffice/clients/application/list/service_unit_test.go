package list_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/list"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()

	subscriptionId := ids.NewSubscriptionId()
	subscriptionAuthor, err := vo.NewUserAuthorWithMemberships(ids.NewUserId().String(), ids.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	adminId := ids.NewUserId()
	systemAuthor, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)

	tests := []struct {
		desc       string
		author     *vo.UserAuthor
		wantFilter *criteria.Criteria
		repoErr    error
		wantCode   app_errors.ErrorCode
	}{
		{
			desc:   "filters a subscription caller by allowed subscription ids",
			author: subscriptionAuthor,
			wantFilter: &criteria.Criteria{Filters: []criteria.Filter{
				{Field: "subscription_id", Operator: criteria.AnyOf, Value: []string{subscriptionId.String()}},
			}},
		},
		{
			desc:       "returns every client for a system administrator without filtering",
			author:     systemAuthor,
			wantFilter: &criteria.Criteria{},
		},
		{
			desc:     "rejects a missing author",
			author:   nil,
			wantCode: "MISSING_AUTHOR",
		},
		{
			desc:       "wraps repository errors",
			author:     systemAuthor,
			wantFilter: &criteria.Criteria{},
			repoErr:    fmt.Errorf("db error"),
			wantCode:   clients.ErrClientsNotListed,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := clientmocks.NewMockRepository(t)
			expected := []*clients.Client{}
			if tc.wantFilter != nil {
				repo.EXPECT().Matching(mock.Anything, *tc.wantFilter).Return(expected, tc.repoErr)
			}

			result, err := list.NewService(repo, sharedDomainServices.NewSubscriptionAuthorizer()).List(context.Background(), list.NewQuery(tc.author))

			if tc.wantCode != "" {
				if appErr, ok := err.(*app_errors.AppError); ok {
					assert.Equal(t, tc.wantCode, appErr.Code)
				} else {
					require.Error(t, err)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	}
}
