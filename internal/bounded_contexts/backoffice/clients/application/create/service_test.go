package create_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves client with correct data and publishes ClientCreated event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		cmd := aCommand(t, withName("my-client"))

		var saved *clients.Client
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(c *clients.Client) bool {
				saved = c
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ClientCreatedType
			}))

		require.NoError(t, newService(t, repo, bus).Create(context.Background(), cmd))

		p := saved.ToPrimitives()
		assert.Equal(t, "my-client", p.Name)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *create.Command
		setupRepo func(repo *clientmocks.MockRepository) // nil = Save no s'hauria de cridar
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects empty name",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withName("")) },
			wantCode: clients.ErrClientNotCreated,
		},
		{
			desc:     "rejects invalid redirect uri",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withRedirectUris("ftp://bad.com")) },
			wantCode: clients.ErrClientNotCreated,
		},
		{
			desc:     "rejects invalid flow",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withGrantTypes("implicit")) },
			wantCode: clients.ErrClientNotCreated,
		},
		{
			desc: "rejects public client with client_credentials",
			cmd: func(t *testing.T) *create.Command {
				return aCommand(t, withGrantTypes("client_credentials"), withPublic(true))
			},
			wantCode: clients.ErrClientNotCreated,
		},
		{
			desc: "rejects duplicate name",
			cmd:  func(t *testing.T) *create.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*clients.Client{{}}, nil)
			},
			wantCode: clients.ErrClientDuplicateName,
		},
		{
			desc: "wraps repository save failure",
			cmd:  func(t *testing.T) *create.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*clients.Client")).
					Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo, bus := newMocks(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			assertErrCode(t, newService(t, repo, bus).Create(context.Background(), tc.cmd(t)), tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newMocks(t *testing.T) (*clientmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return clientmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func newService(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus) *create.Service {
	t.Helper()
	return create.NewService(domain_services.NewDuplicateNameChecker(repo), repo, bus)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
