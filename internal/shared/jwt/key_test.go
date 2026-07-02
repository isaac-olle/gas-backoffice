package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		bytes     []byte
		method    string
		id        string
		wantError bool
	}{
		{"valid key", []byte("keydata"), "RS256", "key-id-1", false},
		{"empty bytes", []byte{}, "RS256", "key-id-1", true},
		{"nil bytes", nil, "RS256", "key-id-1", true},
		{"empty method", []byte("keydata"), "", "key-id-1", true},
		{"empty id", []byte("keydata"), "RS256", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			k, err := NewKey(tc.bytes, tc.method, tc.id)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.bytes, k.Bytes())
				assert.Equal(t, tc.method, k.Method())
				assert.Equal(t, tc.id, k.ID())
			}
		})
	}
}

func TestKeyIsRSA(t *testing.T) {
	t.Parallel()
	tests := []struct {
		method string
		isRSA  bool
	}{
		{"RS256", true},
		{"RS384", true},
		{"RS512", true},
		{"ES256", false},
		{"HS256", false},
	}
	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			k, err := NewKey([]byte("data"), tc.method, "id")
			require.NoError(t, err)
			assert.Equal(t, tc.isRSA, k.IsRSA())
		})
	}
}

func TestKeyIsECDSA(t *testing.T) {
	t.Parallel()
	tests := []struct {
		method  string
		isECDSA bool
	}{
		{"ES256", true},
		{"ES384", true},
		{"ES512", true},
		{"RS256", false},
		{"HS256", false},
	}
	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			k, err := NewKey([]byte("data"), tc.method, "id")
			require.NoError(t, err)
			assert.Equal(t, tc.isECDSA, k.IsECDSA())
		})
	}
}
