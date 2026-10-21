package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUrl(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid http", "http://example.com/path", false},
		{"valid https", "https://example.com/path", false},
		{"valid with port", "https://example.com:8080/path", false},
		{"invalid scheme ftp", "ftp://example.com", true},
		{"no scheme", "example.com", true},
		{"empty", "", true},
		{"no host", "http://", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := NewUrl(tc.input)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.input, u.String())
			}
		})
	}
}

func TestNewUrl_DecodesPercentEncodedInput(t *testing.T) {
	t.Parallel()

	u, err := NewUrl("http%3A%2F%2F127.0.0.1%3A5173%2Fcallback")
	require.NoError(t, err)
	assert.Equal(t, "http://127.0.0.1:5173/callback", u.String())
}

func TestNewUrl_FragmentNotAllowedOption(t *testing.T) {
	t.Parallel()
	opt := NewFragmentNotAllowedOption()

	_, err := NewUrl("https://example.com/path#section", opt)
	require.Error(t, err, "url with fragment should fail")

	_, err = NewUrl("https://example.com/path", opt)
	require.NoError(t, err, "url without fragment should pass")
}

func TestNewUrl_LocalHostOption(t *testing.T) {
	t.Parallel()
	opt := NewLocalHostOption()

	_, err := NewUrl("http://localhost:3000", opt)
	require.Error(t, err, "localhost url should fail")

	_, err = NewUrl("https://example.com", opt)
	require.NoError(t, err, "non-localhost url should pass")
}

func TestNewUrl_HttpsOnlyOption(t *testing.T) {
	t.Parallel()
	opt := NewHttpsOnlyOption()

	_, err := NewUrl("http://example.com", opt)
	require.Error(t, err, "http should fail with https-only option")

	_, err = NewUrl("https://example.com", opt)
	require.NoError(t, err, "https should pass with https-only option")
}

func TestUrlString(t *testing.T) {
	t.Parallel()
	u, err := NewUrl("https://example.com/callback")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/callback", u.String())
}
