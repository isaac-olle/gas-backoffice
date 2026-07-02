package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScopeName(t *testing.T) {
	t.Parallel()

	t.Run("accepts valid name", func(t *testing.T) {
		t.Parallel()
		n, err := NewScopeName("read:users")
		require.NoError(t, err)
		assert.Equal(t, "read:users", n.String())
	})

	t.Run("rejects empty name", func(t *testing.T) {
		t.Parallel()
		_, err := NewScopeName("")
		assert.Error(t, err)
	})
}

func TestScopeNameFromPrimitives(t *testing.T) {
	t.Parallel()
	n := ScopeNameFromPrimitives("read:items")
	assert.Equal(t, "read:items", n.String())
}
