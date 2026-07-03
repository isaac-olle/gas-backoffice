package issue_client_credentials

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	clientFinder   *domain_services.ClientFinder
	tokenGenerator session.TokenGenerator
	eventBus       event.Bus
}

func NewService(
	clientFinder *domain_services.ClientFinder,
	tokenGenerator session.TokenGenerator,
	eventBus event.Bus,
) *Service {
	return &Service{
		clientFinder:   clientFinder,
		tokenGenerator: tokenGenerator,
		eventBus:       eventBus,
	}
}

func (s *Service) IssueClientCredentials(ctx context.Context, qry *Query) (any, error) {
	client, err := s.clientFinder.Find(ctx, qry.clientID)
	if err != nil {
		return nil, err
	}

	if err = client.CheckClientCredentials(qry.scope); err != nil {
		return nil, err
	}

	ccSession, err := session.NewClientCredentialsSession(client.ClientID, client.TenantID, client.AdminID, qry.scope)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenGenerator.GenerateFromClientCredentials(ctx, ccSession)
	if err != nil {
		return nil, app_errors.InternalError(session.ErrCannotSignToken).
			WithDescription("failed to generate token").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, ccSession.PullDomainEvents())

	return &session.TokenResponse{Token: token, Scope: qry.scope}, nil
}
