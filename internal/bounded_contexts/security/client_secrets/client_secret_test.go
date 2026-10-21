package client_secrets_test

import (
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

// stubHasher stubs client_secrets.SecretHasher.
type stubHasher struct {
	hash string
	err  error
}

func (h stubHasher) Hash(_ string) (string, error) { return h.hash, h.err }

type fixture struct {
	id       ids.ClientSecretId
	clientId ids.ClientId
	adminId  ids.UserId
	userId   ids.UserId
}

func newFixture(t *testing.T) fixture {
	t.Helper()
	id, err := ids.NewClientSecretIdFromString(uuid.New().String())
	require.NoError(t, err)
	clientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	return fixture{id: id, clientId: clientId, adminId: ids.NewUserId(), userId: ids.NewUserId()}
}

func newSecret(t *testing.T, f fixture, hash string) *client_secrets.ClientSecret {
	t.Helper()
	cs, err := client_secrets.NewClientSecret(f.id, f.clientId, f.adminId, f.userId, "plain", stubHasher{hash: hash})
	require.NoError(t, err)
	return cs
}

func TestNewClientSecret(t *testing.T) {
	t.Parallel()

	f := newFixture(t)

	tests := []struct {
		name        string
		plainSecret string
		hasher      stubHasher
		wantErr     bool
		wantEvents  int
	}{
		{
			name:        "creates client secret and records domain event",
			plainSecret: "plain",
			hasher:      stubHasher{hash: "hashed-secret-value"},
			wantEvents:  1,
		},
		{
			name:        "rejects empty plain secret",
			plainSecret: "",
			hasher:      stubHasher{hash: "hash"},
			wantErr:     true,
		},
		{
			name:        "fails when hasher fails",
			plainSecret: "plain",
			hasher:      stubHasher{err: errors.New("hash error")},
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cs, err := client_secrets.NewClientSecret(f.id, f.clientId, f.adminId, f.userId, tc.plainSecret, tc.hasher)
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
		userId := uuid.New().String()
		adminId := uuid.New().String()

		p := client_secrets.ClientSecretPrimitives{
			Id:         id,
			ClientId:   clientId,
			UserId:     userId,
			AdminId:    adminId,
			SecretHash: "somehash",
		}

		cs, err := client_secrets.FromPrimitives(p)
		require.NoError(t, err)
		require.NotNil(t, cs)
		assert.Empty(t, cs.PullDomainEvents())

		prim := cs.ToPrimitives()
		assert.Equal(t, id, prim.Id)
		assert.Equal(t, clientId, prim.ClientId)
		assert.Equal(t, userId, prim.UserId)
		assert.Equal(t, adminId, prim.AdminId)
		assert.Equal(t, "somehash", prim.SecretHash)
	})

	t.Run("rejects empty secret hash in primitives", func(t *testing.T) {
		t.Parallel()
		p := client_secrets.ClientSecretPrimitives{
			Id:         uuid.New().String(),
			ClientId:   uuid.New().String(),
			UserId:     uuid.New().String(),
			AdminId:    uuid.New().String(),
			SecretHash: "",
		}
		_, err := client_secrets.FromPrimitives(p)
		assert.Error(t, err)
	})

	t.Run("rejects invalid ownership ids in primitives", func(t *testing.T) {
		t.Parallel()
		p := client_secrets.ClientSecretPrimitives{
			Id:         uuid.New().String(),
			ClientId:   uuid.New().String(),
			UserId:     uuid.New().String(),
			AdminId:    "bad-uuid",
			SecretHash: "hash",
		}
		_, err := client_secrets.FromPrimitives(p)
		assert.Error(t, err)
	})
}

func TestClientSecretEnsureOwnership(t *testing.T) {
	t.Parallel()

	f := newFixture(t)
	cs := newSecret(t, f, "hash")

	otherClientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	siblingTenantId := ids.NewUserId()
	otherAdminId := ids.NewUserId()

	tests := []struct {
		name      string
		ownership client_secrets.Ownership
		wantErr   bool
	}{
		{
			name:      "allows the owning user",
			ownership: client_secrets.NewOwnershipFromVos(f.clientId, f.adminId, f.userId),
		},
		{
			name:      "allows the owning admin acting as itself",
			ownership: client_secrets.NewOwnershipFromVos(f.clientId, f.adminId, f.adminId),
		},
		{
			name:      "rejects a sibling tenant under the same admin",
			ownership: client_secrets.NewOwnershipFromVos(f.clientId, f.adminId, siblingTenantId),
			wantErr:   true,
		},
		{
			name:      "rejects a different admin",
			ownership: client_secrets.NewOwnershipFromVos(f.clientId, otherAdminId, otherAdminId),
			wantErr:   true,
		},
		{
			name:      "rejects a different client even for the owning user",
			ownership: client_secrets.NewOwnershipFromVos(otherClientId, f.adminId, f.userId),
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := cs.EnsureOwnership(tc.ownership)
			if tc.wantErr {
				assertErrCode(t, err, client_secrets.ErrInvalidOwnership)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestClientSecretNewOwnership_InvalidUUID(t *testing.T) {
	t.Parallel()
	valid := uuid.New().String()
	tests := []struct {
		name                      string
		clientId, adminId, userId string
	}{
		{name: "invalid client id", clientId: "bad-uuid", adminId: valid, userId: valid},
		{name: "invalid admin id", clientId: valid, adminId: "bad-uuid", userId: valid},
		{name: "invalid user id", clientId: valid, adminId: valid, userId: "bad-uuid"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := client_secrets.NewOwnership(tc.clientId, tc.adminId, tc.userId)
			assertErrCode(t, err, client_secrets.ErrInvalidOwnership)
		})
	}
}

func TestClientSecretOwnership_ToPrimitives(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	o := client_secrets.NewOwnershipFromVos(f.clientId, f.adminId, f.userId)
	p := o.ToPrimitives()
	assert.Equal(t, f.clientId.String(), p.ClientId)
	assert.Equal(t, f.adminId.String(), p.AdminId)
	assert.Equal(t, f.userId.String(), p.UserId)
}

func TestClientSecretToPrimitives(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	cs := newSecret(t, f, "hashed")
	p := cs.ToPrimitives()
	assert.Equal(t, f.clientId.String(), p.ClientId)
	assert.Equal(t, f.adminId.String(), p.AdminId)
	assert.Equal(t, f.userId.String(), p.UserId)
	assert.Equal(t, "hashed", p.SecretHash)
}
