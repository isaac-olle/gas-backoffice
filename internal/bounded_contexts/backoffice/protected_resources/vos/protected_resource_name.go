package vos

import "fmt"

type ProtectedResourceName string

func NewProtectedResourceName(raw string) (ProtectedResourceName, error) {
	if raw == "" {
		return "", fmt.Errorf("protected resource name cannot be empty")
	}
	return ProtectedResourceName(raw), nil
}

func ProtectedResourceNameFromPrimitives(raw string) ProtectedResourceName {
	return ProtectedResourceName(raw)
}

func (n ProtectedResourceName) String() string { return string(n) }
