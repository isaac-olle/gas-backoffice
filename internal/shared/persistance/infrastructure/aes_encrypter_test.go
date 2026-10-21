package infrastructure_test

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
)

func newTestKey(t *testing.T) string {
	t.Helper()
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(key)
}

func TestAESEncrypter_RoundTrip(t *testing.T) {
	t.Parallel()
	enc, err := infrastructure.NewAESEncrypter(newTestKey(t))
	require.NoError(t, err)

	plaintext := "super-secret-client-secret"
	ciphertext, err := enc.Encrypt(plaintext)
	require.NoError(t, err)
	assert.NotEqual(t, plaintext, ciphertext)

	got, err := enc.Decrypt(ciphertext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, got)
}

func TestAESEncrypter_NonceIsRandom(t *testing.T) {
	t.Parallel()
	enc, err := infrastructure.NewAESEncrypter(newTestKey(t))
	require.NoError(t, err)

	a, err := enc.Encrypt("same")
	require.NoError(t, err)
	b, err := enc.Encrypt("same")
	require.NoError(t, err)

	// Same plaintext must produce different ciphertexts (fresh nonce each time).
	assert.NotEqual(t, a, b)
}

func TestAESEncrypter_WrongKeySize(t *testing.T) {
	t.Parallel()
	shortKey := base64.StdEncoding.EncodeToString(make([]byte, 16))
	_, err := infrastructure.NewAESEncrypter(shortKey)
	assert.Error(t, err)
}

func TestAESEncrypter_InvalidBase64Key(t *testing.T) {
	t.Parallel()
	_, err := infrastructure.NewAESEncrypter("not-base64!!!")
	assert.Error(t, err)
}

func TestAESEncrypter_TamperedCiphertextFails(t *testing.T) {
	t.Parallel()
	enc, err := infrastructure.NewAESEncrypter(newTestKey(t))
	require.NoError(t, err)

	ciphertext, err := enc.Encrypt("data")
	require.NoError(t, err)

	// Flip the last base64 char to simulate tampering; GCM auth must reject it.
	tampered := ciphertext[:len(ciphertext)-1]
	if ciphertext[len(ciphertext)-1] == 'A' {
		tampered += "B"
	} else {
		tampered += "A"
	}
	_, err = enc.Decrypt(tampered)
	assert.Error(t, err)
}

func TestAESEncrypter_DecryptWithDifferentKeyFails(t *testing.T) {
	t.Parallel()
	enc1, err := infrastructure.NewAESEncrypter(newTestKey(t))
	require.NoError(t, err)
	enc2, err := infrastructure.NewAESEncrypter(newTestKey(t))
	require.NoError(t, err)

	ciphertext, err := enc1.Encrypt("data")
	require.NoError(t, err)

	_, err = enc2.Decrypt(ciphertext)
	assert.Error(t, err)
}
