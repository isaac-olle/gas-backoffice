package vos

import "fmt"

type ScopeName string

func NewScopeName(raw string) (ScopeName, error) {
	if raw == "" {
		return "", fmt.Errorf("scope name cannot be empty")
	}
	return ScopeName(raw), nil
}

func ScopeNameFromPrimitives(raw string) ScopeName {
	return ScopeName(raw)
}

func (n ScopeName) String() string { return string(n) }
