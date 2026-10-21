package client_projection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection"
)

func TestClientProjectionIsOwnedBy(t *testing.T) {
	t.Parallel()

	const (
		admin   = "admin"
		tenant  = "tenant"
		sibling = "sibling"
		other   = "other-admin"
	)
	p := &client_projection.ClientProjection{ClientId: "client", UserId: tenant, AdminId: admin}

	tests := []struct {
		name                        string
		callerUserId, callerAdminId string
		want                        bool
	}{
		{name: "owning tenant", callerUserId: tenant, callerAdminId: admin, want: true},
		{name: "owning admin acting as itself", callerUserId: admin, callerAdminId: admin, want: true},
		{name: "sibling tenant under the same admin", callerUserId: sibling, callerAdminId: admin, want: false},
		{name: "different admin", callerUserId: other, callerAdminId: other, want: false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, p.IsOwnedBy(tc.callerUserId, tc.callerAdminId))
		})
	}
}
