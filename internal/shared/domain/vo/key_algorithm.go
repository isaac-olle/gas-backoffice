package vo

import (
	"fmt"
)

var AllowedMethods = []string{"RS256", "RS384", "RS512", "ES256", "ES384", "ES512"}

type Algorithm struct {
	value string
}

func NewAlgorithm(value string) (Algorithm, error) {
	for _, m := range AllowedMethods {
		if m == value {
			return Algorithm{value: value}, nil
		}
	}
	return Algorithm{}, fmt.Errorf("unsupported algorithm '%s'", value)
}

func (a Algorithm) String() string { return a.value }
