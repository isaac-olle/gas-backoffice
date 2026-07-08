package create_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/create"
)

type cmdOpt func(*cmdBuilder)

type cmdBuilder struct {
	prId      ids2.ProtectedResourceId
	name      string
	ownership protected_resources.Ownership
}

func withCmdName(name string) cmdOpt {
	return func(b *cmdBuilder) { b.name = name }
}

func aCommand(t *testing.T, opts ...cmdOpt) *create.Command {
	t.Helper()
	b := &cmdBuilder{
		prId:      ids2.NewProtectedResourceId(),
		name:      "My Resource",
		ownership: protected_resources.NewOwnership(ids2.NewUserId(), ids2.NewAdminId()),
	}
	for _, o := range opts {
		o(b)
	}
	return create.NewCommand(b.prId, b.name, b.ownership)
}
