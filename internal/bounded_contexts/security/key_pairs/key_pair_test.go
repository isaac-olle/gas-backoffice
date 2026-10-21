package key_pairs_test

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	keyPairInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestNewKeyPair(t *testing.T) {
	t.Parallel()

	validId := uuid.New().String()
	validPrId := uuid.New().String()

	stubPrivKeyPEM, stubPubKeyPEM, err := keyPairInfra.NewPEMKeyGenerator().Generate("RS256")
	require.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		prId       string
		algorithm  string
		privKeyPEM string
		pubKeyPEM  string
		wantErr    bool
		wantEvents int
	}{
		{
			name:       "creates key pair and records domain event",
			id:         validId,
			prId:       validPrId,
			algorithm:  "RS256",
			privKeyPEM: stubPrivKeyPEM,
			pubKeyPEM:  stubPubKeyPEM,
			wantErr:    false,
			wantEvents: 1,
		},
		{
			name:       "rejects invalid id",
			id:         "not-a-uuid",
			prId:       validPrId,
			algorithm:  "RS256",
			privKeyPEM: stubPrivKeyPEM,
			pubKeyPEM:  stubPubKeyPEM,
			wantErr:    true,
		},
		{
			name:       "rejects invalid protected resource id",
			id:         validId,
			prId:       "not-a-uuid",
			algorithm:  "RS256",
			privKeyPEM: stubPrivKeyPEM,
			pubKeyPEM:  stubPubKeyPEM,
			wantErr:    true,
		},
		{
			name:       "rejects empty keys",
			id:         uuid.New().String(),
			prId:       uuid.New().String(),
			algorithm:  "RS256",
			privKeyPEM: "",
			pubKeyPEM:  "",
			wantErr:    true,
		},
		{
			name:       "rejects unsupported algorithm",
			id:         uuid.New().String(),
			prId:       uuid.New().String(),
			algorithm:  "HS256",
			privKeyPEM: stubPrivKeyPEM,
			pubKeyPEM:  stubPubKeyPEM,
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			kp, err := key_pairs.NewKeyPair(tc.id, tc.prId, tc.algorithm, tc.privKeyPEM, tc.pubKeyPEM)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, kp)
			assert.Len(t, kp.PullDomainEvents(), tc.wantEvents)
		})
	}
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs from primitives without domain events", func(t *testing.T) {
		t.Parallel()
		id := uuid.New().String()
		prId := uuid.New().String()

		p := key_pairs.KeyPairPrimitives{
			Id:                  id,
			ProtectedResourceId: prId,
			Algorithm:           "RS256",
			PrivateKey:          []byte("raw-private-key-bytes"),
			PublicKey:           []byte("raw-public-key-bytes"),
		}

		kp, err := key_pairs.FromPrimitives(p)
		require.NoError(t, err)
		require.NotNil(t, kp)
		assert.Empty(t, kp.PullDomainEvents())

		prim := kp.ToPrimitives()
		assert.Equal(t, id, prim.Id)
		assert.Equal(t, prId, prim.ProtectedResourceId)
		assert.Equal(t, "RS256", prim.Algorithm)
	})
}

func TestEnsureOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns nil for matching ownership", func(t *testing.T) {
		t.Parallel()
		kp := newKeyPairMother(t)
		prId := kp.ToPrimitives().ProtectedResourceId
		ownership, err := key_pairs.NewOwnership(prId)
		require.NoError(t, err)
		assert.NoError(t, kp.EnsureOwnership(ownership))
	})

	t.Run("returns error for different protected resource id", func(t *testing.T) {
		t.Parallel()
		kp := newKeyPairMother(t)
		otherOwnership, err := key_pairs.NewOwnership(uuid.New().String())
		require.NoError(t, err)
		assert.Error(t, kp.EnsureOwnership(otherOwnership))
	})
}

func TestNewOwnership_InvalidUUID(t *testing.T) {
	t.Parallel()
	_, err := key_pairs.NewOwnership("invalid-uuid")
	assertErrCode(t, err, key_pairs.ErrInvalidOwnership)
}

func TestKeyPair_Getters(t *testing.T) {
	t.Parallel()
	kp := newKeyPairMother(t)
	assert.NotEmpty(t, kp.Algorithm().String())
	assert.NotNil(t, kp.PrivateKey())
	assert.NotNil(t, kp.PublicKey())
}

func TestNewOwnershipFromVos(t *testing.T) {
	t.Parallel()
	prId, err := ids.NewProtectedResourceIdFromString(uuid.New().String())
	require.NoError(t, err)
	o := key_pairs.NewOwnershipFromVos(prId)
	assert.Equal(t, prId.String(), o.ProtectedResourceId.String())
}

func TestKeyPairOwnership_ToPrimitives(t *testing.T) {
	t.Parallel()
	kp := newKeyPairMother(t)
	prId := kp.ToPrimitives().ProtectedResourceId
	o, err := key_pairs.NewOwnership(prId)
	require.NoError(t, err)
	p := o.ToPrimitives()
	assert.Equal(t, prId, p.ProtectedResourceId)
}

func TestKeyPair_ToPrimitivesRoundTrip(t *testing.T) {
	t.Parallel()
	kp := newKeyPairMother(t)
	p := kp.ToPrimitives()
	assert.NotEmpty(t, p.Algorithm)
	assert.NotEmpty(t, p.PrivateKey)
	assert.NotEmpty(t, p.PublicKey)
}
