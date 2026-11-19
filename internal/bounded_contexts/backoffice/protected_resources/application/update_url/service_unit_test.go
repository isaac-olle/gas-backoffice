package update_url_test

import (
	"context"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/update_url"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedservices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateURL(t *testing.T) {
	for _, tc := range []struct {
		name, url                            string
		invalid, forbidden, system, conflict bool
		want                                 app_errors.ErrorCode
	}{
		{name: "update", url: "https://api.example.com"},
		{name: "clear", url: ""},
		{name: "reject fragment", url: "https://api.example.com#section", invalid: true, want: protected_resources.ErrProtectedResourceNotUpdated},
		{name: "reject relative", url: "/api", invalid: true, want: protected_resources.ErrProtectedResourceNotUpdated},
		{name: "reject other subscription", url: "https://api.example.com", forbidden: true, want: protected_resources.ErrInvalidOwnership},
		{name: "reject system resource", url: "https://api.example.com", system: true, want: protected_resources.ErrSystemResource},
		{name: "concurrent modification", url: "https://api.example.com", conflict: true, want: protected_resources.ErrProtectedResourceConcurrentModification},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := prmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			id, subscriptionID := ids.NewProtectedResourceId(), ids.NewSubscriptionId()
			pr, err := protected_resources.NewProtectedResourceWithURL(id, "Resource", protected_resources.NewOwnership(subscriptionID), tc.system, "http://localhost:8090")
			require.NoError(t, err)
			pr.PullDomainEvents()
			authorSubscription := subscriptionID
			if tc.forbidden {
				authorSubscription = ids.NewSubscriptionId()
			}
			author, err := vo.NewUserAuthorWithMemberships(ids.NewUserId().String(), ids.NewUserId().String(), vo.SubscriptionUserRole,
				[]vo.SubscriptionMembership{{SubscriptionId: authorSubscription, SubscriptionRole: vo.UserSubscriptionRole}})
			require.NoError(t, err)
			if !tc.invalid {
				repo.EXPECT().SearchById(mock.Anything, id).Return(pr, nil)
			}
			if !tc.invalid && !tc.forbidden && !tc.system {
				var saveErr error
				if tc.conflict {
					saveErr = persistance.ErrConcurrentModification
				}
				repo.EXPECT().Save(mock.Anything, pr).Return(saveErr)
				if !tc.conflict {
					bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(events []event.Event) bool {
						return len(events) == 1 && events[0].Type() == backoffice.ProtectedResourceURLUpdatedType && events[0].Primitives().Data["url"] == tc.url
					})).Return()
				}
			}
			svc := update_url.NewService(domain_services.NewProtectedResourceFinder(repo), repo, bus, sharedservices.NewSubscriptionAuthorizer())
			err = svc.UpdateURL(context.Background(), update_url.NewCommand(id, tc.url, author))
			if tc.want != "" {
				var appError *app_errors.AppError
				require.ErrorAs(t, err, &appError)
				assert.Equal(t, tc.want, appError.Code)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.url, pr.ToPrimitives().URL)
		})
	}
}
