package create_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

// ── Command Mother ─────────────────────────────────────────────────────────────

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	id     ids2.SessionId
	userId ids2.UserId
}

func withUserId(userId ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

func aCommand(t *testing.T, opts ...commandOpt) *create.Command {
	t.Helper()
	b := &commandBuilder{
		id:     ids2.NewSessionId(),
		userId: ids2.NewUserId(),
	}
	for _, o := range opts {
		o(b)
	}
	cmd := create.NewCommand(b.id, b.userId)
	require.NotNil(t, cmd)
	return cmd
}
