package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/vos"
)

func TestNewState(t *testing.T) {
	t.Parallel()
	s, err := sessionvos.NewState("random-state-xyz")
	require.NoError(t, err)
	assert.Equal(t, "random-state-xyz", s.Value())
	assert.Equal(t, "random-state-xyz", s.String())
}
