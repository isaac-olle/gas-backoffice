package initiate_login

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Result struct {
	AuthorizationURL string
}

type Service struct {
	idpConfigRepo auth.IdpConfigRepository
	connector     auth.IdpConnector
}

func NewService(
	idpConfigRepo auth.IdpConfigRepository,
	connector auth.IdpConnector,
) *Service {
	return &Service{
		idpConfigRepo: idpConfigRepo,
		connector:     connector,
	}
}

func (s *Service) Initiate(ctx context.Context, qry *Query) (any, error) {
	idpName, err := vos.NewFromString(qry.idpReferenceName)
	if err != nil {
		return nil, app_errors.BadRequestError(auth.ErrInvalidIdpReferenceName).
			WithDescription("invalid idp reference name").WithCause(err, false).
			WithDetails(map[string]interface{}{
				"idpReferenceName": qry.idpReferenceName,
			})
	}
	idp, err := s.idpConfigRepo.FindByReferenceName(ctx, idpName)
	if err != nil {
		return nil, app_errors.InternalError(auth.ErrInitiateLoginFailed).
			WithDescription("idp not found").WithCause(err, false)
	}

	authURL, err := s.connector.Connect(ctx, idp)
	if err != nil {
		return nil, app_errors.InternalError(auth.ErrInitiateLoginFailed).
			WithDescription("failed to initiate login").WithCause(err, false)
	}

	return Result{AuthorizationURL: authURL}, nil
}
