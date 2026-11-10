package domain_services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func aClient(t *testing.T) *clients.Client {
	t.Helper()
	c, err := clients.FromPrimitives(clients.ClientPrimitives{
		Id:             ids.NewClientId().String(),
		Name:           "Test Client",
		SubscriptionId: ids.NewSubscriptionId().String(),
	})
	require.NoError(t, err)
	return c
}

func TestClientFinder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		repoClient *clients.Client
		repoErr    error
		wantErr    bool
	}{
		{
			name:       "returns client when repo finds it",
			repoClient: aClient(t),
		},
		{
			name:    "propagates error from repository",
			repoErr: fmt.Errorf("db error"),
			wantErr: true,
		},
		{
			name:       "returns not found error when repo returns nil",
			repoClient: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := clientmocks.NewMockRepository(t)
			repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(tt.repoClient, tt.repoErr)

			finder := domain_services.NewClientFinder(repo)
			got, err := finder.Find(context.Background(), ids.NewClientId())

			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				return
			}
			assert.Equal(t, tt.repoClient, got)
		})
	}
}
