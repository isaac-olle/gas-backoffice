package domain_services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services"
	cakmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func aClientAssertionKey(t *testing.T) *client_assertion_keys.ClientAssertionKey {
	t.Helper()
	cak, err := client_assertion_keys.FromPrimitives(client_assertion_keys.ClientAssertionKeyPrimitives{
		Id:             ids.NewClientAssertionKeyId().String(),
		ClientId:       ids.NewClientId().String(),
		SubscriptionId: ids.NewSubscriptionId().String(),
		Algorithm:      "RS256",
		PublicKey:      []byte("stub-public-key"),
	})
	require.NoError(t, err)
	return cak
}

func TestClientAssertionKeyFinder_Find(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repoKey *client_assertion_keys.ClientAssertionKey
		repoErr error
		wantErr bool
	}{
		{
			name:    "returns key when repo finds it",
			repoKey: aClientAssertionKey(t),
		},
		{
			name:    "propagates error from repository",
			repoErr: fmt.Errorf("db error"),
			wantErr: true,
		},
		{
			name:    "returns not found error when repo returns nil",
			repoKey: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := cakmocks.NewMockClientAssertionKeyRepository(t)
			repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ClientAssertionKeyId")).Return(tt.repoKey, tt.repoErr)

			finder := domain_services.NewClientAssertionKeyFinder(repo)
			got, err := finder.Find(context.Background(), ids.NewClientAssertionKeyId())

			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				return
			}
			assert.Equal(t, tt.repoKey, got)
		})
	}
}

func TestClientAssertionKeyFinder_FindByClientId(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repoKey *client_assertion_keys.ClientAssertionKey
		repoErr error
		wantErr bool
	}{
		{
			name:    "returns key when repo finds it by client id",
			repoKey: aClientAssertionKey(t),
		},
		{
			name:    "propagates error from repository",
			repoErr: fmt.Errorf("db error"),
			wantErr: true,
		},
		{
			name:    "returns not found error when repo returns nil",
			repoKey: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := cakmocks.NewMockClientAssertionKeyRepository(t)
			repo.EXPECT().FindByClientId(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(tt.repoKey, tt.repoErr)

			finder := domain_services.NewClientAssertionKeyFinder(repo)
			got, err := finder.FindByClientId(context.Background(), ids.NewClientId())

			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				return
			}
			assert.Equal(t, tt.repoKey, got)
		})
	}
}
