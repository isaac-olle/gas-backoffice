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
		{"system administrator role", "system_administrator", SystemAdministratorRole, false},
		{"subscription user role", "subscription_user", SubscriptionUserRole, false},
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
	assert.True(t, SystemAdministratorRole.IsSystemAdministrator())
	assert.False(t, SystemAdministratorRole.IsSubscriptionUser())
	assert.True(t, SubscriptionUserRole.IsSubscriptionUser())
	assert.False(t, SubscriptionUserRole.IsSystemAdministrator())
}

func TestRoleString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "system_administrator", SystemAdministratorRole.String())
	assert.Equal(t, "subscription_user", SubscriptionUserRole.String())
}

func TestRoleEquals(t *testing.T) {
	t.Parallel()
	assert.True(t, SystemAdministratorRole.Equals(SystemAdministratorRole))
	assert.False(t, SystemAdministratorRole.Equals(SubscriptionUserRole))
}
