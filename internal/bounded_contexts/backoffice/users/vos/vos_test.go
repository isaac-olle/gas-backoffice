package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserName(t *testing.T) {
	t.Parallel()

	t.Run("accepts valid name", func(t *testing.T) {
		t.Parallel()
		n, err := NewUserName("My User")
		require.NoError(t, err)
		assert.Equal(t, "My User", n.String())
	})

	invalidCases := []struct {
		desc string
		name string
	}{
		{"rejects empty name", ""},
		{"rejects reserved name Gas-User", "Gas-User"},
	}
	for _, tc := range invalidCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, err := NewUserName(tc.name)
			assert.Error(t, err)
		})
	}
}

func TestUserNameFromPrimitives(t *testing.T) {
	t.Parallel()
	n := UserNameFromPrimitives("My User")
	assert.Equal(t, "My User", n.String())
}
