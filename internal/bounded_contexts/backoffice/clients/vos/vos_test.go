package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
)

func TestNewClientName(t *testing.T) {
	t.Parallel()

	t.Run("accepts a non-empty name", func(t *testing.T) {
		t.Parallel()
		n, err := vos.NewClientName("My Client")
		require.NoError(t, err)
		assert.Equal(t, "My Client", n.String())
	})

	t.Run("rejects empty name", func(t *testing.T) {
		t.Parallel()
		_, err := vos.NewClientName("")
		assert.Error(t, err)
	})
}

func TestNewGrantType(t *testing.T) {
	t.Parallel()

	validGrantTypes := []string{
		vos.AuthorizationCodeGrantType,
		vos.ClientCredentialsGrantType,
	}
	for _, raw := range validGrantTypes {
		raw := raw
		t.Run("accepts grant type "+raw, func(t *testing.T) {
			t.Parallel()
			grantType, err := vos.NewGrantType(raw)
			require.NoError(t, err)
			assert.Equal(t, raw, grantType.String())
		})
	}

	t.Run("rejects unsupported grant type", func(t *testing.T) {
		t.Parallel()
		_, err := vos.NewGrantType("implicit")
		assert.Error(t, err)
	})
}

func TestNewRedirectUri(t *testing.T) {
	t.Parallel()

	t.Run("accepts valid https uri", func(t *testing.T) {
		t.Parallel()
		uri, err := vos.NewRedirectUri("https://example.com/callback")
		require.NoError(t, err)
		assert.Equal(t, "https://example.com/callback", uri.String())
	})

	errCases := []struct {
		desc string
		raw  string
	}{
		{"rejects uri with fragment", "https://example.com/callback#fragment"},
		{"rejects ftp scheme", "ftp://example.com/callback"},
		{"rejects empty string", ""},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, err := vos.NewRedirectUri(tc.raw)
			assert.Error(t, err)
		})
	}
}
