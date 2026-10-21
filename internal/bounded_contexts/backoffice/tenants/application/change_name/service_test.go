package change_name_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/change_name"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
)

func TestChangeName(t *testing.T) {
	t.Parallel()

	t.Run("saves tenant with updated name and publishes UserNameChanged event", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		userId := ids2.NewUserId()
		tenant := aTenant(t, userId, nil, "Old Name")

		var saved *tenants.Tenant
		repo.EXPECT().SearchById(mock.Anything, userId).Return(tenant, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *tenants.Tenant) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.UserNameChangedType
			}))

		svc := change_name.NewService(domain_services.NewTenantFinder(repo), repo, bus)
		require.NoError(t, svc.ChangeName(context.Background(), aCommand(t, withUserId(userId), withName("New Name"))))

		assert.Equal(t, "New Name", saved.ToPrimitives().Name)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *change_name.Command
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrUserNameNotChanged for empty name",
			cmd:      func(t *testing.T) *change_name.Command { return aCommand(t, withName("")) },
			wantCode: tenants.ErrUserNameNotChanged,
		},
		{
			desc: "returns ErrUserNotFound when tenant does not exist",
			cmd:  func(t *testing.T) *change_name.Command { return aCommand(t, withName("Valid Name")) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), nil)
			},
			wantCode: tenants.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotFound when repository find returns error",
			cmd:  func(t *testing.T) *change_name.Command { return aCommand(t, withName("Valid Name")) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrUserNotFound,
		},
		{
			desc: "returns internal error when save fails",
			cmd: func(t *testing.T) *change_name.Command {
				userId := ids2.NewUserId()
				return aCommand(t, withUserId(userId), withName("Valid Name"))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, id ids2.UserId) (*tenants.Tenant, error) {
						return aTenant(t, id, nil, "Old Name"), nil
					})
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*tenants.Tenant")).
					Return(fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrUserNotSaved,
		},
		{
			desc: "returns conflict when repository reports concurrent modification",
			cmd: func(t *testing.T) *change_name.Command {
				userId := ids2.NewUserId()
				return aCommand(t, withUserId(userId), withName("Valid Name"))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, id ids2.UserId) (*tenants.Tenant, error) {
						return aTenant(t, id, nil, "Old Name"), nil
					})
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*tenants.Tenant")).
					Return(persistance.ErrConcurrentModification)
			},
			wantCode: tenants.ErrUserConcurrentModification,
		},
		{
			desc: "returns forbidden when caller is not self or owning admin",
			cmd: func(t *testing.T) *change_name.Command {
				return aCommand(t, withName("Valid Name"),
					withAuthor(&authors.UserAuthor{Id: ids2.NewUserId(), Role: vo.TenantRole}))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aTenant(t, ids2.NewUserId(), nil, "Old Name"), nil)
			},
			wantCode: tenants.ErrInvalidOwnership,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := tenantmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := change_name.NewService(domain_services.NewTenantFinder(repo), repo, bus)
			err := svc.ChangeName(context.Background(), tc.cmd(t))

			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// --- helpers ---

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
