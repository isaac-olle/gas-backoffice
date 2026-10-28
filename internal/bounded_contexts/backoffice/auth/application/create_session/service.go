package create_session

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/_shared/user_projection"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Result struct {
	AccessToken string
}

type Service struct {
	tokenClient    auth.TokenEndpointClient
	stateDecoder   auth.StateDecoder
	userProjection user_projection.Repository
	tokenGenerator auth.TokenGenerator
	eventBus       event.Bus
}

func NewService(
	tokenClient auth.TokenEndpointClient,
	stateDecoder auth.StateDecoder,
	userProjection user_projection.Repository,
	tokenGenerator auth.TokenGenerator,
	eventBus event.Bus,
) *Service {
	return &Service{
		tokenClient:    tokenClient,
		stateDecoder:   stateDecoder,
		userProjection: userProjection,
		tokenGenerator: tokenGenerator,
		eventBus:       eventBus,
	}
}

func (s *Service) Exchange(ctx context.Context, qry *Query) (any, error) {
	//TODO aqui com apunta un redirect sha de mirar si el redirect conte error
	verifier, err := s.stateDecoder.DecodeState(qry.state)
	if err != nil {
		return nil, app_errors.BadRequestError(auth.ErrCannotExchangeToken).
			WithDescription("error while decoding token").WithCause(err, false)
	}

	tokenResult, err := s.tokenClient.ExchangeCode(ctx, qry.code, verifier)
	if err != nil {
		return nil, app_errors.InternalError(auth.ErrCannotExchangeToken).WithCause(err, false).WithDescription("error while exchanging code for token")
	}

	projection, err := s.userProjection.FindByEmail(ctx, tokenResult.Email)
	if err != nil {
		return nil, app_errors.ForbiddenError(auth.ErrCannotExchangeToken).
			WithDescription("user is not part of the system").WithCause(err, false)
	}

	session, err := auth.NewSessionWithMemberships(projection.UserId, projection.SystemRole, projection.AdminId, projection.Memberships)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	accessToken, err := s.tokenGenerator.Generate(ctx, session)
	if err != nil {
		return nil, app_errors.InternalError(auth.ErrCannotCreateLoginSession).WithCause(err, false)
	}

	s.eventBus.Publish(ctx, session.PullDomainEvents())

	return Result{AccessToken: accessToken}, nil
}
