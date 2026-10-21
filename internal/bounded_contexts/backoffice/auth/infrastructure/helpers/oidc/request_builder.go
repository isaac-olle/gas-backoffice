package oidc

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
)

type RequestBuilder struct {
	codeVerifierGenerator *RandomCodeVerifierGenerator
	stateEncoder          *StateEncoder
}

func NewRequestBuilder(codeVerifierGenerator *RandomCodeVerifierGenerator,
	stateEncoder *StateEncoder) *RequestBuilder {
	return &RequestBuilder{
		codeVerifierGenerator: codeVerifierGenerator,
		stateEncoder:          stateEncoder,
	}
}

func (r *RequestBuilder) Build(
	redirectUri,
	clientId,
	idpReferenceName,
	authorizeEndpoint string,
) (string, error) {
	codeVerifier, err := r.codeVerifierGenerator.GenerateCodeVerifier()
	if err != nil {
		return "", fmt.Errorf("oidc: generate code verifier: %w", err)
	}

	state, err := r.stateEncoder.EncodeState(codeVerifier)
	if err != nil {
		return "", fmt.Errorf("oidc: encode state: %w", err)
	}

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", clientId)
	params.Set("redirect_uri", redirectUri)
	params.Set("scope", "openid email")
	params.Set("state", state)
	params.Set("idp_reference_name", idpReferenceName)
	params.Set("code_challenge", challengeFromVerifier(codeVerifier.String()))
	params.Set("code_challenge_method", "S256")

	return authorizeEndpoint + "?" + params.Encode(), nil
}

func challengeFromVerifier(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
