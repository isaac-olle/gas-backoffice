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
		{"tenant role", "tenant", TenantRole, false},
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
	assert.False(t, AdminRole.IsTenant())
	assert.False(t, AdminRole.IsClient())

	assert.True(t, TenantRole.IsTenant())
	assert.False(t, TenantRole.IsAdmin())
	assert.False(t, TenantRole.IsClient())

	assert.True(t, ClientRole.IsClient())
	assert.False(t, ClientRole.IsAdmin())
	assert.False(t, ClientRole.IsTenant())
}

func TestRoleString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "admin", AdminRole.String())
	assert.Equal(t, "tenant", TenantRole.String())
	assert.Equal(t, "client", ClientRole.String())
}

func TestRoleEquals(t *testing.T) {
	t.Parallel()
	assert.True(t, AdminRole.Equals(AdminRole))
	assert.False(t, AdminRole.Equals(TenantRole))
}
