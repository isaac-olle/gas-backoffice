package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs/application/create"
)

type CreateKeyPairOnProtectedResourceCreatedEventHandler struct {
	service *create.Service
}

func NewCreateKeyPairOnProtectedResourceCreatedEventHandler(service *create.Service) *CreateKeyPairOnProtectedResourceCreatedEventHandler {
	return &CreateKeyPairOnProtectedResourceCreatedEventHandler{service: service}
}

func (h *CreateKeyPairOnProtectedResourceCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.ProtectedResourceCreated) error {
	keyPairId := ids.NewKeyPairId()
	return h.service.Create(ctx, create.NewCommand(keyPairId.String(), evt.AggregateID(), "RS256"))
}
