package oidc_url_builders

import "strings"

type OIDCAuthorizationURLBuilderFactory struct {
	builders map[string]OIDCAuthorizationURLBuilder
}

func NewOIDCAuthorizationURLBuilderFactory() *OIDCAuthorizationURLBuilderFactory {
	return &OIDCAuthorizationURLBuilderFactory{
		builders: map[string]OIDCAuthorizationURLBuilder{
			"accounts.google.com":       GoogleAuthorizationURLBuilder{},
			"login.microsoftonline.com": MicrosoftAuthorizationURLBuilder{},
		},
	}
}

func (f *OIDCAuthorizationURLBuilderFactory) ForHost(host string) OIDCAuthorizationURLBuilder {
	if builder, ok := f.builders[strings.ToLower(host)]; ok {
		return builder
	}
	return genericOIDCAuthorizationURLBuilder{}
}
