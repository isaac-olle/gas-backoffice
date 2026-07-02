package soft_delete_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/application/soft_delete"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId string
}

func withScopeId(scopeId string) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

// aCommand builds a soft_delete.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *soft_delete.Command {
	t.Helper()
	b := &commandBuilder{scopeId: uuid.New().String()}
	for _, o := range opts {
		o(b)
	}
	return soft_delete.NewCommand(b.scopeId)
}
