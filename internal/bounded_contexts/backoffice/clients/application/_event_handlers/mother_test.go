package _event_handlers_test

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients"
	clientmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/mocks"
)

// ── Client Mothers ─────────────────────────────────────────────────────────────

func aClient(t *testing.T) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		"valid-client",
		clients.NewOwnership(ids2.NewAdminId(), ids2.NewTenantId()),
		nil,
		nil,
	)
	require.NoError(t, err)
	c.PullDomainEvents() // clear creation event
	return c, clientId
}

func aClientWithScope(t *testing.T, scopeId ids2.ScopeId) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	prims := clients.ClientPrimitives{
		Id:            clientId.String(),
		Name:          "valid-client",
		AdminId:       ids2.NewAdminId().String(),
		TenantId:      ids2.NewTenantId().String(),
		AllowedScopes: []string{scopeId.String()},
	}
	c, err := clients.FromPrimitives(prims)
	require.NoError(t, err)
	return c, clientId
}

// ── Mock Helpers ───────────────────────────────────────────────────────────────

func newMocks(t *testing.T) (*clientmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return clientmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
