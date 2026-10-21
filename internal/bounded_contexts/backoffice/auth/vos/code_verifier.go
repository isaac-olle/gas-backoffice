package vos

import (
	"fmt"
	"regexp"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// Forma del code_verifier segons l'RFC 7636 §4.1:
//
//	code-verifier = 43*128unreserved
//	unreserved    = ALPHA / DIGIT / "-" / "." / "_" / "~"
const (
	MinCodeVerifierLength = 43
	MaxCodeVerifierLength = 128
)

var ErrInvalidCodeVerifier app_errors.ErrorCode = "INVALID_CODE_VERIFIER"

var unreservedPattern = regexp.MustCompile(`^[A-Za-z0-9\-._~]+$`)

// CodeVerifier és el secret que el client guarda mentre dura el flux
// d'autorització i que presenta al token endpoint per demostrar que és qui va
// iniciar el login. Aquest tipus només el representa i el valida; generar-lo és
// feina d'un CodeVerifierGenerator.
type CodeVerifier struct{ value string }

func NewCodeVerifierFromString(raw string) (CodeVerifier, error) {
	if len(raw) < MinCodeVerifierLength || len(raw) > MaxCodeVerifierLength {
		return CodeVerifier{}, app_errors.BadRequestError(ErrInvalidCodeVerifier).
			WithDescription(fmt.Sprintf(
				"code verifier must be between %d and %d characters, got %d",
				MinCodeVerifierLength, MaxCodeVerifierLength, len(raw),
			))
	}

	if !unreservedPattern.MatchString(raw) {
		return CodeVerifier{}, app_errors.BadRequestError(ErrInvalidCodeVerifier).
			WithDescription("code verifier must only contain unreserved characters")
	}

	return CodeVerifier{value: raw}, nil
}

func (c CodeVerifier) String() string { return c.value }

func (c CodeVerifier) Equals(other CodeVerifier) bool { return c.value == other.value }
