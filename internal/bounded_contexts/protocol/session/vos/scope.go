package vos

import (
	"fmt"
	"strings"
)

type Scope struct{ value string }

func NewScope(v string) (Scope, error) {
	if v == "" {
		return Scope{}, fmt.Errorf("scope cannot be empty")
	}
	return Scope{value: v}, nil
}

func (s Scope) Value() string  { return s.value }
func (s Scope) String() string { return s.value }

func (s Scope) Has(scope string) bool {
	for _, grantedScope := range strings.Fields(s.value) {
		if grantedScope == scope {
			return true
		}
	}
	return false
}
