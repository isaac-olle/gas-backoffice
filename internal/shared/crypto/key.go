package crypto

// Key is the common capability of every key type: exporting itself as PEM
// or as a JWK (RFC 7517) — private keys include their private members,
// public/symmetric keys don't.
type Key interface {
	Algorithm() string
	ToPEMFormat() ([]byte, error)
	ToJWK() (map[string]any, error)
}

// SigningKey signs data with whatever algorithm and secret material it
// holds. SignBytes is the raw crypto primitive; SignString signs a
// pre-built signing input (base64(header).base64(payload)) and appends
// the base64 signature — it knows nothing about claims or JWT structure,
// that's TokenBuilder's job.
type SigningKey interface {
	Key
	SignBytes(data []byte) ([]byte, error)
	SignString(signingInput string) (string, error)
}

// VerifyingKey verifies data signed by a SigningKey. VerifyBytes is the raw
// crypto primitive; VerifyString parses and verifies a full JWT string.
type VerifyingKey interface {
	Key
	VerifyBytes(data, signature []byte) error
	VerifyString(tokenString string) (map[string]any, error)
}
