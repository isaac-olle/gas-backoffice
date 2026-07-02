package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sessionvos "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/vos"
)

func TestNewConsumed(t *testing.T) {
	t.Parallel()
	t.Run("false", func(t *testing.T) {
		c := sessionvos.NewConsumed(false)
		assert.False(t, c.IsConsumed())
	})
	t.Run("true", func(t *testing.T) {
		c := sessionvos.NewConsumed(true)
		assert.True(t, c.IsConsumed())
	})
}
