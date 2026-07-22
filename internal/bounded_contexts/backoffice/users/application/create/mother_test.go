package create_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	name string
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func aCommand(t *testing.T, opts ...commandOpt) *create.Command {
	t.Helper()
	b := &commandBuilder{name: "my-user"}
	for _, o := range opts {
		o(b)
	}

	tenantId, err := ids2.NewUserIdFromString(uuid.New().String())
	if err != nil {
		t.Fatal(err)
	}
	adminId := ids2.NewUserId()

	return create.NewCommand(tenantId, adminId, b.name, "tenant")
}
