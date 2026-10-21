package oidc

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
)

// verifierByteLength: 32 bytes en base64url són exactament 43 caràcters, que és
// el mínim que exigeix l'RFC 7636 — el 43 en surt d'aquí.
const verifierByteLength = 32

// RandomCodeVerifierGenerator obté el verifier de crypto/rand. No coneix la
// forma que ha de tenir: construeix el candidat i deixa que el VO el validi.
type RandomCodeVerifierGenerator struct{}

func NewRandomCodeVerifierGenerator() *RandomCodeVerifierGenerator {
	return &RandomCodeVerifierGenerator{}
}

func (g *RandomCodeVerifierGenerator) GenerateCodeVerifier() (vos.CodeVerifier, error) {
	b := make([]byte, verifierByteLength)
	// crypto/rand.Read mai retorna error: si la font del sistema falla, el
	// programa peta irrecuperablement (Go 1.24+).
	_, _ = rand.Read(b)

	return vos.NewCodeVerifierFromString(base64.RawURLEncoding.EncodeToString(b))
}
