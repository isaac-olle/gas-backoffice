package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

func TestNewScope(t *testing.T) {
	t.Parallel()
	t.Run("accepts non-empty scope", func(t *testing.T) {
		s, err := sessionvos.NewScope("openid profile")
		require.NoError(t, err)
		assert.Equal(t, "openid profile", s.Value())
		assert.Equal(t, "openid profile", s.String())
	})

	t.Run("rejects empty scope", func(t *testing.T) {
		_, err := sessionvos.NewScope("")
		assert.Error(t, err)
	})
}
