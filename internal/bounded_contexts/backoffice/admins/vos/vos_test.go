package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdminName_Valid(t *testing.T) {
	t.Parallel()
	n, err := NewAdminName("Alice")
	require.NoError(t, err)
	assert.Equal(t, "Alice", n.String())
}

func TestNewAdminName_Empty_ReturnsError(t *testing.T) {
	t.Parallel()
	_, err := NewAdminName("")
	assert.Error(t, err)
}

func TestAdminNameFromPrimitives(t *testing.T) {
	t.Parallel()
	n := AdminNameFromPrimitives("Alice")
	assert.Equal(t, "Alice", n.String())
}
