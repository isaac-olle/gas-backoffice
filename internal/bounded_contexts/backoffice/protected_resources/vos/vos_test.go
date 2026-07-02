package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProtectedResourceName_Valid(t *testing.T) {
	t.Parallel()
	n, err := NewProtectedResourceName("My Resource")
	require.NoError(t, err)
	assert.Equal(t, "My Resource", n.String())
}

func TestNewProtectedResourceName_Empty_ReturnsError(t *testing.T) {
	t.Parallel()
	_, err := NewProtectedResourceName("")
	assert.Error(t, err)
}

func TestNewProtectedResourceName_Reserved_GasBackoffice_ReturnsError(t *testing.T) {
	t.Parallel()
	_, err := NewProtectedResourceName("Gas-Backoffice")
	assert.Error(t, err)
}

func TestProtectedResourceNameFromPrimitives(t *testing.T) {
	t.Parallel()
	n := ProtectedResourceNameFromPrimitives("My Resource")
	assert.Equal(t, "My Resource", n.String())
}

func TestCheckForbiddenName(t *testing.T) {
	t.Parallel()
	t.Run("returns nil for allowed name", func(t *testing.T) {
		n := ProtectedResourceName("Allowed")
		assert.NoError(t, n.CheckForbiddenName())
	})

	t.Run("returns error for forbidden name Gas-Backoffice", func(t *testing.T) {
		n := ProtectedResourceName("Gas-Backoffice")
		assert.Error(t, n.CheckForbiddenName())
	})
}
