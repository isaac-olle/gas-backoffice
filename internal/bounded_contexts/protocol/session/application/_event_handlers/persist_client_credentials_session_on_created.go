package _event_handlers

import (
	"context"
	protocolEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/protocol"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

type PersistClientCredentialsSessionOnCreatedEventHandler struct {
	repository session.ClientCredentialsSessionRepository
}

func NewPersistClientCredentialsSessionOnCreatedEventHandler(repository session.ClientCredentialsSessionRepository) *PersistClientCredentialsSessionOnCreatedEventHandler {
	return &PersistClientCredentialsSessionOnCreatedEventHandler{repository: repository}
}

func (h *PersistClientCredentialsSessionOnCreatedEventHandler) Handle(ctx context.Context, evt protocolEvents.ClientCredentialsSessionCreated) error {
	sess, err := session.ClientCredentialsSessionFromPrimitives(
		&session.ClientCredentialsSessionPrimitives{
			ID:       evt.AggregateID(),
			ClientID: evt.ClientID,
			TenantID: evt.TenantID,
			AdminID:  evt.AdminID,
			Scope:    evt.Scope,
			IssuedAt: evt.OccurredOn().Format("2006-01-02T15:04:05Z07:00"),
		})
	if err != nil {
		return err
	}
	return h.repository.Save(ctx, sess)
}
