package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
)

type Service struct {
	repository   key_pairs.Repository
	eventBus     event.Bus
	keyGenerator key_pairs.KeyGenerator
}

func NewService(repository key_pairs.Repository, eventBus event.Bus, keyGenerator key_pairs.KeyGenerator) *Service {
	return &Service{repository: repository, eventBus: eventBus, keyGenerator: keyGenerator}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	privateKeyPEM, publicKeyPEM, err := s.keyGenerator.Generate(cmd.algorithm)
	if err != nil {
		return app_errors.InternalError(key_pairs.ErrKeyPairNotCreated).
			WithDescription("failed to generate key pair").WithCause(

			err, false)

	}

	keyPair, err := key_pairs.NewKeyPair(cmd.id, cmd.protectedResourceId, cmd.algorithm, privateKeyPEM, publicKeyPEM)
	if err != nil {
		return app_errors.BadRequestError(key_pairs.ErrKeyPairNotCreated).
			WithDescription("failed to create key pair").WithCause(

			err, false)

	}

	if err = s.repository.Save(ctx, keyPair); err != nil {
		return app_errors.InternalError(key_pairs.ErrKeyPairNotSaved).
			WithDescription("failed to save key pair").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, keyPair.PullDomainEvents())
	return nil
}
