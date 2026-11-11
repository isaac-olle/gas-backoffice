package oidc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"time"
)

type statePayload struct {
	CodeVerifier     string `json:"code_verifier"`
	Exp              int64  `json:"exp"`
	IdpReferenceName string `json:"idp_reference_name"`
}

type StateEncoder struct {
	stateTTl  time.Duration
	encrypter persistance.Encrypter
}

func NewStateEncoder(stateTTl time.Duration,
	encrypter persistance.Encrypter) *StateEncoder {
	return &StateEncoder{
		stateTTl:  stateTTl,
		encrypter: encrypter,
	}
}

func (s *StateEncoder) EncodeState(codeVerifier vos.CodeVerifier, idpReferenceName string) (string, error) {
	payload := statePayload{
		CodeVerifier:     codeVerifier.String(),
		Exp:              time.Now().Add(s.stateTTl).Unix(),
		IdpReferenceName: idpReferenceName,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encrypted, err := s.encrypter.Encrypt(string(raw))
	if err != nil {
		return "", err
	}

	return toURLSafeState(encrypted)
}

func (s *StateEncoder) DecodeState(encodedState string) (auth.AuthorizationState, error) {
	encrypted, err := fromURLSafeState(encodedState)
	if err != nil {
		return auth.AuthorizationState{}, fmt.Errorf("invalid state: %w", err)
	}

	raw, err := s.encrypter.Decrypt(encrypted)
	if err != nil {
		return auth.AuthorizationState{}, fmt.Errorf("invalid state: %w", err)
	}

	var payload statePayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return auth.AuthorizationState{}, fmt.Errorf("malformed state payload: %w", err)
	}

	if time.Now().Unix() > payload.Exp {
		return auth.AuthorizationState{}, fmt.Errorf("state token expired")
	}

	if payload.CodeVerifier == "" {
		return auth.AuthorizationState{}, fmt.Errorf("missing code_verifier in state")
	}

	return auth.AuthorizationState{CodeVerifier: payload.CodeVerifier, IdpReferenceName: payload.IdpReferenceName}, nil
}

func toURLSafeState(encrypted string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("transcode encrypted state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func fromURLSafeState(state string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return "", fmt.Errorf("transcode state: %w", err)
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}
