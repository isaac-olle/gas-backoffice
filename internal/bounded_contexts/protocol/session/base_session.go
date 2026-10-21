package session

import (
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

type BaseSession struct {
	clientID            ids2.ClientId
	scope               sessionvos.Scope
	grantType           sessionvos.GrantType
	protectedResourceId string
	*aggregate2.Root
}

func newBaseSession(
	clientID ids2.ClientId,
	scope sessionvos.Scope,
	grantType sessionvos.GrantType,
	aggregateName aggregate2.Name,
	protectedResourceId string,
) *BaseSession {
	return &BaseSession{
		clientID:            clientID,
		scope:               scope,
		grantType:           grantType,
		protectedResourceId: protectedResourceId,
		Root:                aggregate2.NewRoot(aggregateName),
	}
}

func (b BaseSession) ClientID() ids2.ClientId {
	return b.clientID
}

func (b BaseSession) HasScope(scope string) bool {
	return b.scope.Has(scope)
}

func (b BaseSession) ToPrimitives() *BaseSessionPrimitives {
	return &BaseSessionPrimitives{
		ClientID:            b.clientID.String(),
		Scope:               b.scope.String(),
		GrantType:           b.grantType.String(),
		ProtectedResourceId: b.protectedResourceId,
	}
}
