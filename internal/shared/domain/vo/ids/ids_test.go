package ids

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	validUUID   = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	invalidUUID = "not-a-uuid"
)

// UserId

func TestNewUserId(t *testing.T) {
	t.Parallel()
	id1 := NewUserId()
	id2 := NewUserId()
	assert.NotEmpty(t, id1.String())
	assert.NotEqual(t, id1.String(), id2.String())
}

func TestNewUserIdFromString(t *testing.T) {
	t.Parallel()
	id, err := NewUserIdFromString(validUUID)
	require.NoError(t, err)
	assert.Equal(t, validUUID, id.String())

	_, err = NewUserIdFromString(invalidUUID)
	require.Error(t, err)
}

// ClientId

func TestNewClientId(t *testing.T) {
	t.Parallel()
	id1 := NewClientId()
	id2 := NewClientId()
	assert.NotEmpty(t, id1.String())
	assert.NotEqual(t, id1.String(), id2.String())
}

func TestNewClientIdFromString(t *testing.T) {
	t.Parallel()
	id, err := NewClientIdFromString(validUUID)
	require.NoError(t, err)
	assert.Equal(t, validUUID, id.String())

	_, err = NewClientIdFromString(invalidUUID)
	require.Error(t, err)
}

// ScopeId

func TestNewScopeId(t *testing.T) {
	t.Parallel()
	id1 := NewScopeId()
	id2 := NewScopeId()
	assert.NotEmpty(t, id1.String())
	assert.NotEqual(t, id1.String(), id2.String())
}

func TestNewScopeIdFromString(t *testing.T) {
	t.Parallel()
	id, err := NewScopeIdFromString(validUUID)
	require.NoError(t, err)
	assert.Equal(t, validUUID, id.String())

	_, err = NewScopeIdFromString(invalidUUID)
	require.Error(t, err)
}

// ProtectedResourceId

func TestNewProtectedResourceId(t *testing.T) {
	t.Parallel()
	id1 := NewProtectedResourceId()
	id2 := NewProtectedResourceId()
	assert.NotEmpty(t, id1.String())
	assert.NotEqual(t, id1.String(), id2.String())
}

func TestNewProtectedResourceIdFromString(t *testing.T) {
	t.Parallel()
	id, err := NewProtectedResourceIdFromString(validUUID)
	require.NoError(t, err)
	assert.Equal(t, validUUID, id.String())

	_, err = NewProtectedResourceIdFromString(invalidUUID)
	require.Error(t, err)
}

// KeyPairId

func TestNewKeyPairId(t *testing.T) {
	t.Parallel()
	id1 := NewKeyPairId()
	id2 := NewKeyPairId()
	assert.NotEmpty(t, id1.String())
	assert.NotEqual(t, id1.String(), id2.String())
}

func TestNewKeyPairIdFromString(t *testing.T) {
	t.Parallel()
	id, err := NewKeyPairIdFromString(validUUID)
	require.NoError(t, err)
	assert.Equal(t, validUUID, id.String())

	_, err = NewKeyPairIdFromString(invalidUUID)
	require.Error(t, err)
}

// ClientAssertionKeyId

func TestNewClientAssertionKeyId(t *testing.T) {
	t.Parallel()
	id1 := NewClientAssertionKeyId()
	id2 := NewClientAssertionKeyId()
	assert.NotEmpty(t, id1.String())
	assert.NotEqual(t, id1.String(), id2.String())
}

func TestNewClientAssertionKeyIdFromString(t *testing.T) {
	t.Parallel()
	id, err := NewClientAssertionKeyIdFromString(validUUID)
	require.NoError(t, err)
	assert.Equal(t, validUUID, id.String())

	_, err = NewClientAssertionKeyIdFromString(invalidUUID)
	require.Error(t, err)
}

// ClientSecretId

func TestNewClientSecretId(t *testing.T) {
	t.Parallel()
	id1 := NewClientSecretId()
	id2 := NewClientSecretId()
	assert.NotEmpty(t, id1.String())
	assert.NotEqual(t, id1.String(), id2.String())
}

func TestNewClientSecretIdFromString(t *testing.T) {
	t.Parallel()
	id, err := NewClientSecretIdFromString(validUUID)
	require.NoError(t, err)
	assert.Equal(t, validUUID, id.String())

	_, err = NewClientSecretIdFromString(invalidUUID)
	require.Error(t, err)
}

// SessionId

func TestNewSessionId(t *testing.T) {
	t.Parallel()
	id1 := NewSessionId()
	id2 := NewSessionId()
	assert.NotEmpty(t, id1.String())
	assert.NotEqual(t, id1.String(), id2.String())
}

func TestNewSessionIdFromString(t *testing.T) {
	t.Parallel()
	id, err := NewSessionIdFromString(validUUID)
	require.NoError(t, err)
	assert.Equal(t, validUUID, id.String())

	_, err = NewSessionIdFromString(invalidUUID)
	require.Error(t, err)
}

// Table-driven test for error codes on invalid input
func TestIdFromString_ErrorCodes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		fn   func(string) error
	}{
		{"user", func(s string) error { _, err := NewUserIdFromString(s); return err }},
		{"client", func(s string) error { _, err := NewClientIdFromString(s); return err }},
		{"scope", func(s string) error { _, err := NewScopeIdFromString(s); return err }},
		{"protected_resource", func(s string) error { _, err := NewProtectedResourceIdFromString(s); return err }},
		{"key_pair", func(s string) error { _, err := NewKeyPairIdFromString(s); return err }},
		{"client_assertion_key", func(s string) error { _, err := NewClientAssertionKeyIdFromString(s); return err }},
		{"client_secret", func(s string) error { _, err := NewClientSecretIdFromString(s); return err }},
		{"session", func(s string) error { _, err := NewSessionIdFromString(s); return err }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn(invalidUUID)
			require.Error(t, err, "expected error for invalid UUID")
		})
	}
}
