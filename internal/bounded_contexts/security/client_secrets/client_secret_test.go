package client_secrets_test

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestNewClientSecret(t *testing.T) {
	t.Parallel()

	clientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		secretHash string
		wantErr    bool
		wantEvents int
	}{
		{
			name:       "creates client secret and records domain event",
			id:         uuid.New().String(),
			secretHash: "hashed-secret-value",
			wantEvents: 1,
		},
		{
			name:       "rejects invalid id",
			id:         "bad-id",
			secretHash: "hash",
			wantErr:    true,
		},
		{
			name:       "rejects empty secret hash",
			id:         uuid.New().String(),
			secretHash: "",
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cs, err := client_secrets.NewClientSecret(tc.id, clientId, tc.secretHash)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, cs)
			assert.Len(t, cs.PullDomainEvents(), tc.wantEvents)
		})
	}
}

func TestClientSecretFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs from primitives without domain events", func(t *testing.T) {
		t.Parallel()
		id := uuid.New().String()
		clientId := uuid.New().String()

		p := client_secrets.ClientSecretPrimitives{
			Id:         id,
			ClientId:   clientId,
			SecretHash: "somehash",
		}

		cs, err := client_secrets.FromPrimitives(p)
		require.NoError(t, err)
		require.NotNil(t, cs)
		assert.Empty(t, cs.PullDomainEvents())

		prim := cs.ToPrimitives()
		assert.Equal(t, id, prim.Id)
		assert.Equal(t, clientId, prim.ClientId)
		assert.Equal(t, "somehash", prim.SecretHash)
	})

	t.Run("rejects empty secret hash in primitives", func(t *testing.T) {
		t.Parallel()
		p := client_secrets.ClientSecretPrimitives{
			Id:         uuid.New().String(),
			ClientId:   uuid.New().String(),
			SecretHash: "",
		}
		_, err := client_secrets.FromPrimitives(p)
		assert.Error(t, err)
	})
}

func TestClientSecretCheckOwnership(t *testing.T) {
	t.Parallel()

	clientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	cs, err := client_secrets.NewClientSecret(uuid.New().String(), clientId, "hash")
	require.NoError(t, err)

	t.Run("returns true for matching client id", func(t *testing.T) {
		t.Parallel()
		assert.True(t, cs.CheckOwnership(client_secrets.NewOwnershipFromVos(clientId)))
	})

	t.Run("returns false for different client id", func(t *testing.T) {
		t.Parallel()
		otherId, err := ids.NewClientIdFromString(uuid.New().String())
		require.NoError(t, err)
		assert.False(t, cs.CheckOwnership(client_secrets.NewOwnershipFromVos(otherId)))
	})
}

func TestClientSecretNewOwnership_InvalidUUID(t *testing.T) {
	t.Parallel()
	_, err := client_secrets.NewOwnership("bad-uuid")
	assertErrCode(t, err, client_secrets.ErrInvalidOwnership)
}

func TestClientSecretOwnership_ToPrimitives(t *testing.T) {
	t.Parallel()
	clientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	o := client_secrets.NewOwnershipFromVos(clientId)
	p := o.ToPrimitives()
	assert.Equal(t, clientId.String(), p.ClientId)
}

func TestClientSecretToPrimitives(t *testing.T) {
	t.Parallel()
	clientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	cs, err := client_secrets.NewClientSecret(uuid.New().String(), clientId, "hashed")
	require.NoError(t, err)
	p := cs.ToPrimitives()
	assert.Equal(t, clientId.String(), p.ClientId)
	assert.Equal(t, "hashed", p.SecretHash)
}
