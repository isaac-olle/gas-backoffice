package exchange_code

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Service struct {
	sessionFinder  *domain_services.SessionFinder
	clientFinder   *domain_services.ClientFinder
	repository     session.Repository
	tokenGenerator session.TokenGenerator
}

func NewService(
	sessionFinder *domain_services.SessionFinder,
	clientFinder *domain_services.ClientFinder,
	repository session.Repository,
	tokenGenerator session.TokenGenerator,
) *Service {
	return &Service{
		sessionFinder:  sessionFinder,
		clientFinder:   clientFinder,
		repository:     repository,
		tokenGenerator: tokenGenerator,
	}
}

func (s *Service) ExchangeCode(ctx context.Context, qry *Query) (any, error) {
	sess, err := s.sessionFinder.FindByCode(ctx, qry.code)
	if err != nil {
		return nil, err
	}

	if err = sess.CheckExchangeParams(qry.clientID, qry.redirectURI); err != nil {
		return nil, err
	}

	if _, err = sess.CheckPkce(qry.codeVerifier); err != nil {
		return nil, err
	}

	client, err := s.clientFinder.Find(ctx, qry.clientID)
	if err != nil {
		return nil, err
	}

	sess.Consume()
	if err = s.repository.Save(ctx, sess); err != nil {
		return nil, app_errors.InternalError(session.ErrCannotSaveSession).
			WithDescription("failed to consume session").
			WithCause(err)
	}

	token, err := s.tokenGenerator.GenerateFromSession(ctx, sess, client)
	if err != nil {
		return nil, app_errors.InternalError(session.ErrCannotSignToken).
			WithDescription("failed to generate token").
			WithCause(err)
	}

	return &session.TokenResponse{Token: token, Scope: sess.ToPrimitives().Scope}, nil
}
