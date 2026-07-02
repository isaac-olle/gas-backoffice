package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTenantName(t *testing.T) {
	t.Parallel()

	t.Run("accepts valid name", func(t *testing.T) {
		t.Parallel()
		n, err := NewTenantName("My Tenant")
		require.NoError(t, err)
		assert.Equal(t, "My Tenant", n.String())
	})

	invalidCases := []struct {
		desc string
		name string
	}{
		{"rejects empty name", ""},
		{"rejects reserved name Gas-Tenant", "Gas-Tenant"},
	}
	for _, tc := range invalidCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, err := NewTenantName(tc.name)
			assert.Error(t, err)
		})
	}
}

func TestTenantNameFromPrimitives(t *testing.T) {
	t.Parallel()
	n := TenantNameFromPrimitives("My Tenant")
	assert.Equal(t, "My Tenant", n.String())
}
