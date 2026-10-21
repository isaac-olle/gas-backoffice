package jwt_verifiers

import (
	"fmt"
	"time"
)

// ExpirationVerifier checks that a claims map's "exp" claim hasn't passed.
// It knows nothing about signatures or any other claim — only whether this
// token is still within its validity window.
type ExpirationVerifier struct{}

func NewExpirationVerifier() *ExpirationVerifier {
	return &ExpirationVerifier{}
}

func (v *ExpirationVerifier) Verify(claims map[string]any) error {
	expRaw, ok := claims["exp"]
	if !ok {
		return fmt.Errorf("missing exp claim")
	}

	exp, ok := expRaw.(float64)
	if !ok {
		return fmt.Errorf("invalid exp claim")
	}

	if time.Now().Unix() > int64(exp) {
		return fmt.Errorf("token expired")
	}

	return nil
}
