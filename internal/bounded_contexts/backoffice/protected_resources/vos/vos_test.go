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

func TestProtectedResourceNameFromPrimitives(t *testing.T) {
	t.Parallel()
	n := ProtectedResourceNameFromPrimitives("My Resource")
	assert.Equal(t, "My Resource", n.String())
}
