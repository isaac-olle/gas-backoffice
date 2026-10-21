package client_assertion_keys_test

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

// stubPubKeyHex is a valid hex-encoded RSA 2048-bit PKIX public key for tests.
const stubPubKeyHex = "30820122300d06092a864886f70d01010105000382010f003082010a0282010100ad18fdd1191ba86551729d71622eb72fe64a6aa322dc8545866e9a75233e48d26c76212470c8eb5eef222a5c4dd7d91f24537f7c394ac19b8381d28bab6fd742a2a940722d17a92035bbb37e43566cb93379b5a7e6ca240bd06a812b47e089e86f93e583c4f8bf63e8ec9227f34f7c6c6b7134f2ce6c6ffe85e5fb2503590f2a1c5a94160f81b806bbe08a67ca187e1acef0662c747875ac82ae28bb191d885b5918e0c9a15fce1334b0aa167cd975dd06b30c86247567caba54bc42456879abfbdec8f2f0c1bc53c9db0d9d117375f810b97491ea98b8ebeaa57ee83aaae7a1a7f0371eb2d471b624b6d27a7867670c48282ab60a75975cc74a330efeb197590203010001"

func newValidIds(t *testing.T) (ids2.ClientAssertionKeyId, ids2.ClientId) {
	t.Helper()
	cakId, err := ids2.NewClientAssertionKeyIdFromString(uuid.New().String())
	require.NoError(t, err)
	clientId, err := ids2.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	return cakId, clientId
}

// newOwnerIds returns a tenant user and its owning admin.
func newOwnerIds(t *testing.T) (adminId, userId ids2.UserId) {
	t.Helper()
	return ids2.NewUserId(), ids2.NewUserId()
}

func TestNewClientAssertionKey(t *testing.T) {
	t.Parallel()

	cakId, clientId := newValidIds(t)
	adminId, userId := newOwnerIds(t)

	tests := []struct {
		name       string
		cakId      ids2.ClientAssertionKeyId
		clientId   ids2.ClientId
		pubKeyHex  string
		wantErr    bool
		wantEvents int
	}{
		{
			name:       "creates client assertion key and records domain event",
			cakId:      cakId,
			clientId:   clientId,
			pubKeyHex:  stubPubKeyHex,
			wantEvents: 1,
		},
		{
			name:      "rejects empty public key",
			cakId:     cakId,
			clientId:  clientId,
			pubKeyHex: "",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cak, err := client_assertion_keys.NewClientAssertionKey(tc.cakId, tc.clientId, adminId, userId, "RS256", tc.pubKeyHex)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, cak)
			assert.Len(t, cak.PullDomainEvents(), tc.wantEvents)
		})
	}
}

func TestClientAssertionKeyFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs from primitives without domain events", func(t *testing.T) {
		t.Parallel()
		id := uuid.New().String()
		clientId := uuid.New().String()
		userId := uuid.New().String()
		adminId := uuid.New().String()

		p := client_assertion_keys.ClientAssertionKeyPrimitives{
			Id:        id,
			ClientId:  clientId,
			UserId:    userId,
			AdminId:   adminId,
			Algorithm: "RS256",
			PublicKey: []byte("some-raw-key-bytes"),
		}

		cak, err := client_assertion_keys.FromPrimitives(p)
		require.NoError(t, err)
		require.NotNil(t, cak)
		assert.Empty(t, cak.PullDomainEvents())

		prim := cak.ToPrimitives()
		assert.Equal(t, id, prim.Id)
		assert.Equal(t, clientId, prim.ClientId)
		assert.Equal(t, userId, prim.UserId)
		assert.Equal(t, adminId, prim.AdminId)
	})

	t.Run("rejects invalid ownership ids in primitives", func(t *testing.T) {
		t.Parallel()
		p := client_assertion_keys.ClientAssertionKeyPrimitives{
			Id:        uuid.New().String(),
			ClientId:  uuid.New().String(),
			UserId:    "not-a-uuid",
			AdminId:   uuid.New().String(),
			Algorithm: "RS256",
			PublicKey: []byte("some-raw-key-bytes"),
		}
		_, err := client_assertion_keys.FromPrimitives(p)
		assert.Error(t, err)
	})
}

func TestClientAssertionKeyEnsureOwnership(t *testing.T) {
	t.Parallel()

	cakId, clientId := newValidIds(t)
	adminId, userId := newOwnerIds(t)
	cak, err := client_assertion_keys.NewClientAssertionKey(cakId, clientId, adminId, userId, "RS256", stubPubKeyHex)
	require.NoError(t, err)

	_, otherClientId := newValidIds(t)
	siblingTenantId := ids2.NewUserId()
	otherAdminId := ids2.NewUserId()

	tests := []struct {
		name      string
		ownership client_assertion_keys.Ownership
		wantErr   bool
	}{
		{
			name:      "allows the owning user",
			ownership: client_assertion_keys.NewOwnershipFromVos(clientId, adminId, userId),
		},
		{
			name:      "allows the owning admin acting as itself",
			ownership: client_assertion_keys.NewOwnershipFromVos(clientId, adminId, adminId),
		},
		{
			name:      "rejects a sibling tenant under the same admin",
			ownership: client_assertion_keys.NewOwnershipFromVos(clientId, adminId, siblingTenantId),
			wantErr:   true,
		},
		{
			name:      "rejects a different admin",
			ownership: client_assertion_keys.NewOwnershipFromVos(clientId, otherAdminId, otherAdminId),
			wantErr:   true,
		},
		{
			name:      "rejects a different client even for the owning user",
			ownership: client_assertion_keys.NewOwnershipFromVos(otherClientId, adminId, userId),
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := cak.EnsureOwnership(tc.ownership)
			if tc.wantErr {
				assertErrCode(t, err, client_assertion_keys.ErrInvalidOwnership)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestClientAssertionKeyNewOwnership_InvalidUUID(t *testing.T) {
	t.Parallel()
	valid := uuid.New().String()
	tests := []struct {
		name                      string
		clientId, adminId, userId string
	}{
		{name: "invalid client id", clientId: "not-a-uuid", adminId: valid, userId: valid},
		{name: "invalid admin id", clientId: valid, adminId: "not-a-uuid", userId: valid},
		{name: "invalid user id", clientId: valid, adminId: valid, userId: "not-a-uuid"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := client_assertion_keys.NewOwnership(tc.clientId, tc.adminId, tc.userId)
			assertErrCode(t, err, client_assertion_keys.ErrInvalidOwnership)
		})
	}
}

func TestClientAssertionKeyOwnership_ToPrimitives(t *testing.T) {
	t.Parallel()
	_, clientId := newValidIds(t)
	adminId, userId := newOwnerIds(t)
	o := client_assertion_keys.NewOwnershipFromVos(clientId, adminId, userId)
	p := o.ToPrimitives()
	assert.Equal(t, clientId.String(), p.ClientId)
	assert.Equal(t, adminId.String(), p.AdminId)
	assert.Equal(t, userId.String(), p.UserId)
}

func TestClientAssertionKeyToPrimitives(t *testing.T) {
	t.Parallel()
	cakId, clientId := newValidIds(t)
	adminId, userId := newOwnerIds(t)
	cak, err := client_assertion_keys.NewClientAssertionKey(cakId, clientId, adminId, userId, "RS256", stubPubKeyHex)
	require.NoError(t, err)
	p := cak.ToPrimitives()
	assert.Equal(t, clientId.String(), p.ClientId)
	assert.Equal(t, adminId.String(), p.AdminId)
	assert.Equal(t, userId.String(), p.UserId)
	assert.Equal(t, "RS256", p.Algorithm)
	assert.NotEmpty(t, p.PublicKey)
}
