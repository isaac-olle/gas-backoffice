package identity_providers

import (
	"net/url"
	"strings"

	idpvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var identityProviderAggregateName aggregate.Name = aggregate.IdentityProvider

type IdentityProvider struct {
	*aggregate.Root
	id                    ids.IdentityProviderId
	referenceName         idpvos.ReferenceName
	name                  idpvos.Name
	issuer                idpvos.Endpoint
	clientID              idpvos.ExternalClientID
	clientSecret          string
	authorizationEndpoint idpvos.Endpoint
	tokenEndpoint         idpvos.Endpoint
	userinfoEndpoint      *idpvos.Endpoint
	jwksURI               idpvos.Endpoint
	scopes                idpvos.Scopes
	backofficeRedirectURI idpvos.Endpoint
	enabled               bool
}

func NewIdentityProvider(
	id ids.IdentityProviderId,
	referenceName, name, issuer,
	clientID, clientSecret,
	authorizationEndpoint, tokenEndpoint, userinfoEndpoint, jwksURI string,
	scopes []string,
	backofficeRedirectURI string,
) (*IdentityProvider, error) {
	errArr := app_errors.NewEmptyErrorArray()

	referenceNameVO, err := idpvos.NewReferenceName(referenceName)
	if err != nil {
		errArr.AddError(err)
	}
	nameVO, err := idpvos.NewName(name)
	if err != nil {
		errArr.AddError(err)
	}
	issuerVO, err := idpvos.NewEndpoint(issuer)
	if err != nil {
		errArr.AddError(err)
	}
	clientIDVO, err := idpvos.NewExternalClientID(clientID)
	if err != nil {
		errArr.AddError(err)
	}
	authEndpointVO, err := idpvos.NewEndpoint(authorizationEndpoint)
	if err != nil {
		errArr.AddError(err)
	}
	tokenEndpointVO, err := idpvos.NewEndpoint(tokenEndpoint)
	if err != nil {
		errArr.AddError(err)
	}
	jwksVO, err := idpvos.NewEndpoint(jwksURI)
	if err != nil {
		errArr.AddError(err)
	}
	scopesVO, err := idpvos.NewScopes(scopes)
	if err != nil {
		errArr.AddError(err)
	}
	backofficeRedirectURIVO, err := idpvos.NewEndpoint(backofficeRedirectURI)
	if err != nil {
		errArr.AddError(err)
	}

	var userinfoEndpointVO *idpvos.Endpoint
	if userinfoEndpoint != "" {
		vo, uErr := idpvos.NewEndpoint(userinfoEndpoint)
		if uErr != nil {
			errArr.AddError(uErr)
		} else {
			userinfoEndpointVO = &vo
		}
	}

	if errArr.HasErrors() {
		return nil, app_errors.BadRequestError(ErrInvalidIdentityProvider).
			WithDescription("invalid identity provider").
			WithCause(errArr, true)
	}

	idp := &IdentityProvider{
		Root:                  aggregate.NewRoot(identityProviderAggregateName),
		id:                    id,
		referenceName:         referenceNameVO,
		name:                  nameVO,
		issuer:                issuerVO,
		clientID:              clientIDVO,
		clientSecret:          clientSecret,
		authorizationEndpoint: authEndpointVO,
		tokenEndpoint:         tokenEndpointVO,
		userinfoEndpoint:      userinfoEndpointVO,
		jwksURI:               jwksVO,
		scopes:                scopesVO,
		backofficeRedirectURI: backofficeRedirectURIVO,
		enabled:               true,
	}
	idp.Record(backoffice.NewIdentityProviderCreated(id.String(), referenceName, name))
	return idp, nil
}

func FromPrimitives(p IdentityProviderPrimitives) (*IdentityProvider, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids.NewIdentityProviderIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}
	referenceNameVO, err := idpvos.NewReferenceName(p.ReferenceName)
	if err != nil {
		errArr.AddError(err)
	}
	nameVO, err := idpvos.NewName(p.Name)
	if err != nil {
		errArr.AddError(err)
	}
	issuerVO, err := idpvos.NewEndpoint(p.Issuer)
	if err != nil {
		errArr.AddError(err)
	}
	clientIDVO, err := idpvos.NewExternalClientID(p.ClientId)
	if err != nil {
		errArr.AddError(err)
	}
	authEndpointVO, err := idpvos.NewEndpoint(p.AuthorizationEndpoint)
	if err != nil {
		errArr.AddError(err)
	}
	tokenEndpointVO, err := idpvos.NewEndpoint(p.TokenEndpoint)
	if err != nil {
		errArr.AddError(err)
	}
	jwksVO, err := idpvos.NewEndpoint(p.JwksUri)
	if err != nil {
		errArr.AddError(err)
	}
	scopesVO, err := idpvos.NewScopes(p.Scopes)
	if err != nil {
		errArr.AddError(err)
	}
	backofficeRedirectURIVO, err := idpvos.NewEndpoint(p.BackofficeRedirectUri)
	if err != nil {
		errArr.AddError(err)
	}

	var userinfoEndpointVO *idpvos.Endpoint
	if p.UserinfoEndpoint != "" {
		vo, uErr := idpvos.NewEndpoint(p.UserinfoEndpoint)
		if uErr != nil {
			errArr.AddError(uErr)
		} else {
			userinfoEndpointVO = &vo
		}
	}

	if errArr.HasErrors() {
		return nil, app_errors.InternalError(ErrIdentityProviderMisconfigured).
			WithDescription("malformed identity provider in database").
			WithCause(errArr, false)
	}

	return &IdentityProvider{
		Root:                  aggregate.NewRoot(identityProviderAggregateName),
		id:                    id,
		referenceName:         referenceNameVO,
		name:                  nameVO,
		issuer:                issuerVO,
		clientID:              clientIDVO,
		clientSecret:          p.ClientSecret,
		authorizationEndpoint: authEndpointVO,
		tokenEndpoint:         tokenEndpointVO,
		userinfoEndpoint:      userinfoEndpointVO,
		jwksURI:               jwksVO,
		scopes:                scopesVO,
		backofficeRedirectURI: backofficeRedirectURIVO,
		enabled:               p.Enabled,
	}, nil
}

func (p *IdentityProvider) ToPrimitives() IdentityProviderPrimitives {
	prim := IdentityProviderPrimitives{
		Id:                    p.id.String(),
		ReferenceName:         p.referenceName.String(),
		Name:                  p.name.String(),
		Issuer:                p.issuer.String(),
		ClientId:              p.clientID.String(),
		ClientSecret:          p.clientSecret,
		AuthorizationEndpoint: p.authorizationEndpoint.String(),
		TokenEndpoint:         p.tokenEndpoint.String(),
		JwksUri:               p.jwksURI.String(),
		Scopes:                p.scopes.ToSlice(),
		BackofficeRedirectUri: p.backofficeRedirectURI.String(),
		Enabled:               p.enabled,
	}
	if p.userinfoEndpoint != nil {
		prim.UserinfoEndpoint = p.userinfoEndpoint.String()
	}
	return prim
}

func (p *IdentityProvider) IsEnabled() bool { return p.enabled }

func (p *IdentityProvider) BackofficeRedirectURI() string { return p.backofficeRedirectURI.String() }

// AuthorizationURL builds the front-channel URL the browser is redirected to.
func (p *IdentityProvider) AuthorizationURL(redirectURI, state, nonce string) (string, error) {
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", p.clientID.String())
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", strings.Join(p.scopes.ToSlice(), " "))
	q.Set("state", state)
	q.Set("nonce", nonce)
	return p.authorizationEndpoint.String() + "?" + q.Encode(), nil
}
