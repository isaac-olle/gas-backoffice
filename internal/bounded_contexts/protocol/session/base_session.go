package session

import (
	aggregate2 "github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	sessionvos "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/vos"
)

type BaseSession struct {
	clientID ids2.ClientId
	scope    sessionvos.Scope
	flow     sessionvos.Flow
	*aggregate2.Root
}

func newBaseSession(clientID ids2.ClientId, scope sessionvos.Scope, flow sessionvos.Flow, aggregateName aggregate2.Name) BaseSession {
	return BaseSession{
		clientID: clientID,
		scope:    scope,
		flow:     flow,
		Root:     aggregate2.NewRoot(aggregateName),
	}
}
