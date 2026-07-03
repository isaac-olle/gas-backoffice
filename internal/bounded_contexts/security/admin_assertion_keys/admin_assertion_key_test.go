package admin_assertion_keys_test

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

// stubPubKeyHex is a valid hex-encoded RSA 2048-bit PKIX public key for tests.
const stubPubKeyHex = "30820122300d06092a864886f70d01010105000382010f003082010a0282010100ad18fdd1191ba86551729d71622eb72fe64a6aa322dc8545866e9a75233e48d26c76212470c8eb5eef222a5c4dd7d91f24537f7c394ac19b8381d28bab6fd742a2a940722d17a92035bbb37e43566cb93379b5a7e6ca240bd06a812b47e089e86f93e583c4f8bf63e8ec9227f34f7c6c6b7134f2ce6c6ffe85e5fb2503590f2a1c5a94160f81b806bbe08a67ca187e1acef0662c747875ac82ae28bb191d885b5918e0c9a15fce1334b0aa167cd975dd06b30c86247567caba54bc42456879abfbdec8f2f0c1bc53c9db0d9d117375f810b97491ea98b8ebeaa57ee83aaae7a1a7f0371eb2d471b624b6d27a7867670c48282ab60a75975cc74a330efeb197590203010001"

func TestNewAdminAssertionKey(t *testing.T) {
	t.Parallel()

	adminId, err := ids.NewAdminIdFromString(uuid.New().String())
	require.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		adminId    ids.AdminId
		pubKeyHex  string
		wantErr    bool
		wantEvents int
	}{
		{
			name:       "creates admin assertion key and records domain event",
			id:         uuid.New().String(),
			adminId:    adminId,
			pubKeyHex:  stubPubKeyHex,
			wantEvents: 1,
		},
		{
			name:      "rejects invalid id",
			id:        "bad-id",
			adminId:   adminId,
			pubKeyHex: stubPubKeyHex,
			wantErr:   true,
		},
		{
			name:      "rejects empty public key",
			id:        uuid.New().String(),
			adminId:   adminId,
			pubKeyHex: "",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			aak, err := admin_assertion_keys.NewAdminAssertionKey(tc.id, tc.adminId, "RS256", tc.pubKeyHex)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, aak)
			assert.Len(t, aak.PullDomainEvents(), tc.wantEvents)
		})
	}
}

func TestAdminAssertionKeyFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs from primitives without domain events", func(t *testing.T) {
		t.Parallel()
		id := uuid.New().String()
		adminId := uuid.New().String()
		p := admin_assertion_keys.AdminAssertionKeyPrimitives{
			Id:        id,
			AdminId:   adminId,
			Algorithm: "RS256",
			PublicKey: []byte("some-raw-key-bytes"),
		}

		aak, err := admin_assertion_keys.FromPrimitives(p)
		require.NoError(t, err)
		require.NotNil(t, aak)
		assert.Empty(t, aak.PullDomainEvents())

		prim := aak.ToPrimitives()
		assert.Equal(t, id, prim.Id)
		assert.Equal(t, adminId, prim.AdminId)
	})
}

func TestAdminAssertionKeyCheckOwnership(t *testing.T) {
	t.Parallel()

	adminId, err := ids.NewAdminIdFromString(uuid.New().String())
	require.NoError(t, err)

	aak, err := admin_assertion_keys.NewAdminAssertionKey(
		uuid.New().String(), adminId, "RS256", stubPubKeyHex,
	)
	require.NoError(t, err)

	t.Run("returns true for matching admin id", func(t *testing.T) {
		t.Parallel()
		ownership := admin_assertion_keys.NewOwnershipFromVos(adminId)
		assert.True(t, aak.CheckOwnership(ownership))
	})

	t.Run("returns false for different admin id", func(t *testing.T) {
		t.Parallel()
		otherId, err := ids.NewAdminIdFromString(uuid.New().String())
		require.NoError(t, err)
		assert.False(t, aak.CheckOwnership(admin_assertion_keys.NewOwnershipFromVos(otherId)))
	})
}

func TestAdminAssertionKeyNewOwnership_InvalidUUID(t *testing.T) {
	t.Parallel()
	_, err := admin_assertion_keys.NewOwnership("not-a-uuid")
	assertErrCode(t, err, admin_assertion_keys.ErrInvalidOwnership)
}

func TestAdminAssertionKeyOwnership_ToPrimitives(t *testing.T) {
	t.Parallel()
	adminId, err := ids.NewAdminIdFromString(uuid.New().String())
	require.NoError(t, err)
	o := admin_assertion_keys.NewOwnershipFromVos(adminId)
	p := o.ToPrimitives()
	assert.Equal(t, adminId.String(), p.AdminId)
}

func TestAdminAssertionKeyToPrimitives(t *testing.T) {
	t.Parallel()
	adminId, err := ids.NewAdminIdFromString(uuid.New().String())
	require.NoError(t, err)
	aak, err := admin_assertion_keys.NewAdminAssertionKey(uuid.New().String(), adminId, "RS256", stubPubKeyHex)
	require.NoError(t, err)
	p := aak.ToPrimitives()
	assert.Equal(t, adminId.String(), p.AdminId)
	assert.Equal(t, "RS256", p.Algorithm)
	assert.NotEmpty(t, p.PublicKey)
}
