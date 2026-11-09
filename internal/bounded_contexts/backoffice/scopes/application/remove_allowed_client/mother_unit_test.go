package remove_allowed_client_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/remove_allowed_client"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId             string
	clientId            string
	protectedResourceId string
	userId              string
}

func withScopeId(scopeId string) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

func withClientId(clientId string) commandOpt {
	return func(b *commandBuilder) { b.clientId = clientId }
}

func withProtectedResourceId(protectedResourceId string) commandOpt {
	return func(b *commandBuilder) { b.protectedResourceId = protectedResourceId }
}

func withUserId(userId string) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

func aCommand(t *testing.T, opts ...commandOpt) *remove_allowed_client.Command {
	t.Helper()
	b := &commandBuilder{
		scopeId:             ids2.NewScopeId().String(),
		clientId:            ids2.NewClientId().String(),
		protectedResourceId: ids2.NewProtectedResourceId().String(),
		userId:              ids2.NewUserId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return remove_allowed_client.NewCommand(b.scopeId, b.clientId, b.protectedResourceId, ids2.NewUserId().String(), b.userId)
}

func aScopeWithClient(t *testing.T, clientId ids2.ClientId, prId ids2.ProtectedResourceId, userId ids2.UserId) (*scopes.Scope, ids2.ScopeId) {
	t.Helper()
	scopeId := ids2.NewScopeId()
	prims := scopes.ScopePrimitives{
		Id:                  scopeId.String(),
		Name:                "read:users",
		ProtectedResourceId: prId.String(),
		AdminId:             ids2.NewUserId().String(),
		UserId:              userId.String(),
		AllowedClientIds:    []string{clientId.String()},
	}
	scope, err := scopes.FromPrimitives(prims)
	require.NoError(t, err)
	return scope, scopeId
}
