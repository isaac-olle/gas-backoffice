package vos

import (
	"fmt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

var ErrInvalidUserIdentity app_errors.ErrorCode = "INVALID_USER_IDENTITY"

// IdpClaims holds the identity claims resolved from the identity provider.
type IdpClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

// UserIdentity is the resolved external identity (tokens + claims) for a
// federated login. It lives embedded in AuthorizationCodeSession instead of
// being its own aggregate, since the relationship is a true 1:1.
type UserIdentity struct {
	idpId        ids.IdentityProviderId
	idToken      string
	accessToken  string
	refreshToken *string
	claims       IdpClaims
	expiresAt    time.Time
}

func NewUserIdentity(
	idpIdString string,
	idToken string,
	accessToken string,
	refreshToken *string,
	claims IdpClaims,
	expiresAt time.Time,
) (UserIdentity, error) {
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
		return UserIdentity{}, app_errors.BadRequestError(ErrInvalidUserIdentity).
			WithDescription("several errors building user identity").
			WithCause(errArr, true)
	}

	return UserIdentity{
		idpId:        idpId,
		idToken:      idToken,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		claims:       claims,
		expiresAt:    expiresAt,
	}, nil
}

func (i UserIdentity) IdpId() ids.IdentityProviderId { return i.idpId }
func (i UserIdentity) IDToken() string               { return i.idToken }
func (i UserIdentity) AccessToken() string           { return i.accessToken }
func (i UserIdentity) RefreshToken() *string         { return i.refreshToken }
func (i UserIdentity) Claims() IdpClaims             { return i.claims }
func (i UserIdentity) ExpiresAt() time.Time          { return i.expiresAt }
