package key_pairs

// TokenSpec is everything a TokenGenerator caller alone can know about the
// token it wants: who it's for (Sub), who should consume it (Aud), a
// caller-chosen correlation id (JTI), domain-specific claims (ExtraClaims),
// and whether the header should embed the signing key's material (JWK/X5c).
// Everything else a JWT needs (iss, iat, exp, kid, alg) is decided by
// whoever resolves the KeyPair, since it depends on key material the caller
// doesn't have access to.
type TokenSpec struct {
	Sub         string
	Aud         string
	JTI         string
	ExtraClaims map[string]any
	IncludeX5c  bool
	IncludeJWK  bool
}
