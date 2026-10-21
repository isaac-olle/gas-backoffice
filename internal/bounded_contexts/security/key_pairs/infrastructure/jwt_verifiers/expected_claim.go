package jwt_verifiers

import (
	"fmt"
	"reflect"
)

// ExpectedClaim checks that a claims map contains the expected
// values. It doesn't know which claim names matter — the caller decides
// that — only how to compare them, including the JWT/OIDC peculiarity that
// a claim like "aud" may be a single string or an array of strings.
type ExpectedClaim struct{}

func NewExpectedClaimVerifier() *ExpectedClaim {
	return &ExpectedClaim{}
}

func (v *ExpectedClaim) Verify(claims map[string]any, expectedClaims map[string]any) error {
	for name, expectedValue := range expectedClaims {
		actualValue, exists := claims[name]
		if !exists {
			return fmt.Errorf("missing claim: %s", name)
		}

		if matchesArrayClaim(actualValue, expectedValue) {
			continue
		}

		if !reflect.DeepEqual(actualValue, expectedValue) {
			return fmt.Errorf("invalid claim: %s", name)
		}
	}

	return nil
}

// matchesArrayClaim handles claims that per JWT/OIDC may be either a single
// string or an array of strings (e.g. "aud"), matching if expectedValue is
// a string equal to one of the array's elements.
func matchesArrayClaim(actualValue, expectedValue any) bool {
	expectedStr, ok := expectedValue.(string)
	if !ok {
		return false
	}

	values, ok := actualValue.([]any)
	if !ok {
		return false
	}

	for _, v := range values {
		if s, ok := v.(string); ok && s == expectedStr {
			return true
		}
	}

	return false
}
