package refresh_session

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Result struct {
	AccessToken  string
	RefreshToken string
}

type Service struct {
	refreshTokenVerifier  auth.RefreshTokenVerifier
	accessTokenGenerator  auth.TokenGenerator
	refreshTokenGenerator auth.RefreshTokenGenerator
	sessionRepository     auth.SessionRepository
	eventBus              event.Bus
}

func NewService(refreshTokenVerifier auth.RefreshTokenVerifier, accessTokenGenerator auth.TokenGenerator, refreshTokenGenerator auth.RefreshTokenGenerator, sessionRepository auth.SessionRepository, eventBus event.Bus) *Service {
	return &Service{refreshTokenVerifier: refreshTokenVerifier, accessTokenGenerator: accessTokenGenerator, refreshTokenGenerator: refreshTokenGenerator, sessionRepository: sessionRepository, eventBus: eventBus}
}

func (s *Service) Refresh(ctx context.Context, qry *Query) (any, error) {
	claims, err := s.refreshTokenVerifier.Verify(ctx, qry.refreshToken)
	if err != nil {
		return nil, app_errors.UnauthorizedError("INVALID_REFRESH_TOKEN").WithCause(err, false)
	}

	if claims.UserID == "" || claims.SystemRole == "" {
		return nil, app_errors.UnauthorizedError("INVALID_REFRESH_TOKEN")
	}
	if s.sessionRepository != nil {
		active, checkErr := s.sessionRepository.IsActive(ctx, claims.SessionID)
		if checkErr != nil || !active {
			return nil, app_errors.UnauthorizedError("INVALID_REFRESH_TOKEN")
		}
	}

	session, err := auth.NewSessionWithIdAndMemberships(claims.SessionID, claims.UserID, claims.SystemRole, claims.UserID, claims.Memberships)
	if err != nil {
		return nil, app_errors.UnauthorizedError("INVALID_REFRESH_TOKEN").WithCause(err, false)
	}
	session.SetPicture(claims.Picture)
	accessToken, err := s.accessTokenGenerator.Generate(ctx, session)
	if err != nil {
		return nil, app_errors.InternalError(auth.ErrCannotCreateLoginSession).WithCause(err, false)
	}
	refreshToken, err := s.refreshTokenGenerator.GenerateRefreshToken(ctx, session)
	if err != nil {
		return nil, app_errors.InternalError(auth.ErrCannotCreateLoginSession).WithCause(err, false)
	}
	s.eventBus.Publish(ctx, session.PullDomainEvents())
	return Result{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
