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

func TestNewEmail(t *testing.T) {
	t.Parallel()

	t.Run("accepts empty email", func(t *testing.T) {
		t.Parallel()
		e, err := NewEmail("")
		require.NoError(t, err)
		assert.Equal(t, "", e.String())
	})

	t.Run("accepts valid email", func(t *testing.T) {
		t.Parallel()
		e, err := NewEmail("user@example.com")
		require.NoError(t, err)
		assert.Equal(t, "user@example.com", e.String())
	})

	t.Run("rejects malformed email", func(t *testing.T) {
		t.Parallel()
		_, err := NewEmail("not-an-email")
		assert.Error(t, err)
	})
}

func TestEmailFromPrimitives(t *testing.T) {
	t.Parallel()
	e := EmailFromPrimitives("user@example.com")
	assert.Equal(t, "user@example.com", e.String())
}
