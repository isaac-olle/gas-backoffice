package oidc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"time"
)

type statePayload struct {
	CV  string `json:"cv"`
	Exp int64  `json:"exp"`
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

func (s *StateEncoder) EncodeState(codeVerifier vos.CodeVerifier) (string, error) {
	payload := statePayload{
		CV:  codeVerifier.String(),
		Exp: time.Now().Add(s.stateTTl).Unix(),
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

func (s *StateEncoder) DecodeState(encodedState string) (string, error) {
	encrypted, err := fromURLSafeState(encodedState)
	if err != nil {
		return "", fmt.Errorf("invalid state: %w", err)
	}

	raw, err := s.encrypter.Decrypt(encrypted)
	if err != nil {
		return "", fmt.Errorf("invalid state: %w", err)
	}

	var payload statePayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return "", fmt.Errorf("malformed state payload: %w", err)
	}

	if time.Now().Unix() > payload.Exp {
		return "", fmt.Errorf("state token expired")
	}

	if payload.CV == "" {
		return "", fmt.Errorf("missing code_verifier in state")
	}

	return payload.CV, nil
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
