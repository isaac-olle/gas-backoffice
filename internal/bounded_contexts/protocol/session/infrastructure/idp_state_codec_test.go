package infrastructure_test

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	persistanceInfra "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
)

func newTestEncrypter(t *testing.T) *persistanceInfra.AESEncrypter {
	t.Helper()
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encrypter, err := persistanceInfra.NewAESEncrypter(base64.StdEncoding.EncodeToString(key))
	require.NoError(t, err)
	return encrypter
}

func aPayload() session.IdpStatePayload {
	return session.IdpStatePayload{
		ClientID:            "client-id",
		RedirectURI:         "https://example.com/cb",
		Scope:               "openid",
		ResponseType:        "code",
		CodeChallenge:       "challenge",
		CodeChallengeMethod: "S256",
		ClientState:         "client-state",
		IdpReferenceName:    "google",
	}
}

func TestIdpStateCodec_RoundTrip(t *testing.T) {
	t.Parallel()

	codec := infrastructure.NewIdpStateCodec(10*time.Minute, newTestEncrypter(t))

	encoded, err := codec.Encode(aPayload())
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := codec.Decode(encoded)
	require.NoError(t, err)
	assert.Equal(t, "client-id", decoded.ClientID)
	assert.Equal(t, "https://example.com/cb", decoded.RedirectURI)
	assert.Equal(t, "openid", decoded.Scope)
	assert.Equal(t, "code", decoded.ResponseType)
	assert.Equal(t, "challenge", decoded.CodeChallenge)
	assert.Equal(t, "S256", decoded.CodeChallengeMethod)
	assert.Equal(t, "client-state", decoded.ClientState)
	assert.Equal(t, "google", decoded.IdpReferenceName)
}

func TestIdpStateCodec_RejectsExpired(t *testing.T) {
	t.Parallel()

	codec := infrastructure.NewIdpStateCodec(-1*time.Minute, newTestEncrypter(t))

	encoded, err := codec.Encode(aPayload())
	require.NoError(t, err)

	_, err = codec.Decode(encoded)
	assert.Error(t, err)
}

func TestIdpStateCodec_RejectsTamperedState(t *testing.T) {
	t.Parallel()

	codec := infrastructure.NewIdpStateCodec(10*time.Minute, newTestEncrypter(t))

	encoded, err := codec.Encode(aPayload())
	require.NoError(t, err)

	tampered := []byte(encoded)
	// Flip a byte well inside the payload so the GCM authentication tag no
	// longer matches — tampering must be detected, not silently decoded.
	flip := len(tampered) / 2
	if tampered[flip] == 'A' {
		tampered[flip] = 'B'
	} else {
		tampered[flip] = 'A'
	}

	_, err = codec.Decode(string(tampered))
	assert.Error(t, err)
}

func TestIdpStateCodec_RejectsGarbage(t *testing.T) {
	t.Parallel()

	codec := infrastructure.NewIdpStateCodec(10*time.Minute, newTestEncrypter(t))

	_, err := codec.Decode("not-a-valid-state")
	assert.Error(t, err)
}
