package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

func TestNewRedirectURI(t *testing.T) {
	t.Parallel()
	t.Run("accepts valid https url", func(t *testing.T) {
		r, err := sessionvos.NewRedirectURI("https://example.com/callback")
		require.NoError(t, err)
		assert.Equal(t, "https://example.com/callback", r.String())
		assert.Equal(t, "https://example.com/callback", r.Value())
	})

	t.Run("rejects empty string", func(t *testing.T) {
		_, err := sessionvos.NewRedirectURI("")
		assert.Error(t, err)
	})

	t.Run("rejects url with fragment", func(t *testing.T) {
		_, err := sessionvos.NewRedirectURI("https://example.com/cb#fragment")
		assert.Error(t, err)
	})
}
