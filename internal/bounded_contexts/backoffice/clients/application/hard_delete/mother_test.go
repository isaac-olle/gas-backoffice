package hard_delete_test

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/hard_delete"
)

// ── Command Mother ─────────────────────────────────────────────────────────────

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	clientId string
	adminId  string
	userId   string
}

func withClientId(id string) commandOpt {
	return func(b *commandBuilder) { b.clientId = id }
}

func withUserId(id string) commandOpt {
	return func(b *commandBuilder) { b.userId = id }
}

func aCommand(t *testing.T, opts ...commandOpt) *hard_delete.Command {
	t.Helper()
	b := &commandBuilder{
		clientId: ids.NewClientId().String(),
		adminId:  ids.NewUserId().String(),
		userId:   ids.NewUserId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return hard_delete.NewCommand(b.clientId, b.adminId, b.userId)
}
