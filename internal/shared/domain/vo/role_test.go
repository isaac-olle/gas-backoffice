package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRole(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		wantRole  Role
		wantError bool
	}{
		{"admin role", "admin", AdminRole, false},
		{"user role", "user", UserRole, false},
		{"client role", "client", ClientRole, false},
		{"invalid role", "superuser", "", true},
		{"empty string", "", "", true},
		{"mixed case", "Admin", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			role, err := NewRole(tc.input)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantRole, role)
			}
		})
	}
}

func TestRoleChecks(t *testing.T) {
	t.Parallel()
	assert.True(t, AdminRole.IsAdmin())
	assert.False(t, AdminRole.IsUser())
	assert.False(t, AdminRole.IsClient())

	assert.True(t, UserRole.IsUser())
	assert.False(t, UserRole.IsAdmin())
	assert.False(t, UserRole.IsClient())

	assert.True(t, ClientRole.IsClient())
	assert.False(t, ClientRole.IsAdmin())
	assert.False(t, ClientRole.IsUser())
}

func TestRoleString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "admin", AdminRole.String())
	assert.Equal(t, "user", UserRole.String())
	assert.Equal(t, "client", ClientRole.String())
}

func TestRoleEquals(t *testing.T) {
	t.Parallel()
	assert.True(t, AdminRole.Equals(AdminRole))
	assert.False(t, AdminRole.Equals(UserRole))
}
