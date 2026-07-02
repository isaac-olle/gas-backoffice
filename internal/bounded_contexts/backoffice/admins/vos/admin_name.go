package vos

import "fmt"

type AdminName string

var forbiddenNames = map[string]struct{}{
	"Gas-Seed-Admin": {},
}

func NewAdminName(raw string) (AdminName, error) {
	if raw == "" {
		return "", fmt.Errorf("admin name cannot be empty")
	}
	if _, ok := forbiddenNames[raw]; ok {
		return "", fmt.Errorf("admin name '%s' is reserved", raw)
	}
	return AdminName(raw), nil
}

func AdminNameFromPrimitives(raw string) AdminName {
	return AdminName(raw)
}

func (n AdminName) String() string { return string(n) }
