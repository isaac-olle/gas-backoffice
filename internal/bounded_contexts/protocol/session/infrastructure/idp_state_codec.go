package infrastructure

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

// IdpStateCodec encodes/decodes the self-contained, encrypted "state" value
// sent to and received back from the identity provider. It uses the shared
// AES-256-GCM persistance.Encrypter's Seal/Open (raw bytes) instead of
// Encrypt/Decrypt, to avoid a double base64 transcode: the sealed bytes are
// base64 URL-safe encoded directly, ready to travel in a query string.
type IdpStateCodec struct {
	ttl       time.Duration
	encrypter persistance.Encrypter
}

func NewIdpStateCodec(ttl time.Duration, encrypter persistance.Encrypter) *IdpStateCodec {
	return &IdpStateCodec{ttl: ttl, encrypter: encrypter}
}

func (c *IdpStateCodec) Encode(payload session.IdpStatePayload) (string, error) {
	payload.Exp = time.Now().Add(c.ttl).Unix()

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal idp state payload: %w", err)
	}

	sealed, err := c.encrypter.Seal(raw)
	if err != nil {
		return "", fmt.Errorf("seal idp state payload: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (c *IdpStateCodec) Decode(encoded string) (session.IdpStatePayload, error) {
	sealed, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return session.IdpStatePayload{}, fmt.Errorf("decode idp state: %w", err)
	}

	raw, err := c.encrypter.Open(sealed)
	if err != nil {
		return session.IdpStatePayload{}, fmt.Errorf("open idp state: %w", err)
	}

	var payload session.IdpStatePayload
	if err = json.Unmarshal(raw, &payload); err != nil {
		return session.IdpStatePayload{}, fmt.Errorf("unmarshal idp state payload: %w", err)
	}

	if time.Now().Unix() > payload.Exp {
		return session.IdpStatePayload{}, fmt.Errorf("idp state expired")
	}

	return payload, nil
}
