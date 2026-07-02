package vos

import "fmt"

type ProtectedResourceName string

var forbiddenProtectedResourceNames = map[string]struct{}{
	"Gas-Backoffice": {},
}

func NewProtectedResourceName(raw string) (ProtectedResourceName, error) {
	if raw == "" {
		return "", fmt.Errorf("protected resource name cannot be empty")
	}
	if _, ok := forbiddenProtectedResourceNames[raw]; ok {
		return "", fmt.Errorf("protected resource name '%s' is reserved", raw)
	}
	return ProtectedResourceName(raw), nil
}

func ProtectedResourceNameFromPrimitives(raw string) ProtectedResourceName {
	return ProtectedResourceName(raw)
}

func (n ProtectedResourceName) CheckForbiddenName() error {
	if _, ok := forbiddenProtectedResourceNames[string(n)]; ok {
		return fmt.Errorf("protected resource name '%s' is reserved", n)
	}
	return nil
}

func (n ProtectedResourceName) String() string { return string(n) }
