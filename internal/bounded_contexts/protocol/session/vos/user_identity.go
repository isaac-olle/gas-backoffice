package vos

import (
	"fmt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

var ErrInvalidExternalOidcUserIdentity app_errors.ErrorCode = "INVALID_EXTERNAL_OIDC_USER_IDENTITY"

// IdpClaims holds the identity claims resolved from the identity provider.
type IdpClaims struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Picture string `json:"picture,omitempty"`
}

// ExternalOidcUserIdentity is the resolved external identity (tokens + claims) for a
// federated login. It lives embedded in AuthorizationCodeSession instead of
// being its own aggregate, since the relationship is a true 1:1.
type ExternalOidcUserIdentity struct {
	idpId        ids.IdentityProviderId
	idToken      string
	accessToken  string
	refreshToken *string
	claims       IdpClaims
	expiresAt    time.Time
}

func NewExternalOidcUserIdentity(
	idpIdString string,
	idToken string,
	accessToken string,
	refreshToken *string,
	claims IdpClaims,
	expiresAt time.Time,
) (ExternalOidcUserIdentity, error) {
	errArr := app_errors.NewEmptyErrorArray()

	idpId, err := ids.NewIdentityProviderIdFromString(idpIdString)
	if err != nil {
		errArr.AddError(err)
	}

	if idToken == "" {
		errArr.AddError(fmt.Errorf("id_token cannot be empty"))
	}
	if accessToken == "" {
		errArr.AddError(fmt.Errorf("access_token cannot be empty"))
	}
	if claims.Sub == "" {
		errArr.AddError(fmt.Errorf("claims sub cannot be empty"))
	}
	if claims.Email == "" {
		errArr.AddError(fmt.Errorf("claims email cannot be empty"))
	}

	if errArr.HasErrors() {
		return ExternalOidcUserIdentity{}, app_errors.BadRequestError(ErrInvalidExternalOidcUserIdentity).
			WithDescription("several errors building user identity").
			WithCause(errArr, true)
	}

	return ExternalOidcUserIdentity{
		idpId:        idpId,
		idToken:      idToken,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		claims:       claims,
		expiresAt:    expiresAt,
	}, nil
}

func (i ExternalOidcUserIdentity) IdpId() ids.IdentityProviderId { return i.idpId }
func (i ExternalOidcUserIdentity) IDToken() string               { return i.idToken }
func (i ExternalOidcUserIdentity) AccessToken() string           { return i.accessToken }
func (i ExternalOidcUserIdentity) RefreshToken() *string         { return i.refreshToken }
func (i ExternalOidcUserIdentity) Claims() IdpClaims             { return i.claims }
func (i ExternalOidcUserIdentity) ExpiresAt() time.Time          { return i.expiresAt }
