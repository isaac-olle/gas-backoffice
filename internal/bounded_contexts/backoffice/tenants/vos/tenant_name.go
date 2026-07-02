package vos

import "fmt"

type TenantName string

var forbiddenTenantNames = map[string]struct{}{
	"Gas-Tenant": {},
}

func NewTenantName(raw string) (TenantName, error) {
	if raw == "" {
		return "", fmt.Errorf("tenant name cannot be empty")
	}
	if _, ok := forbiddenTenantNames[raw]; ok {
		return "", fmt.Errorf("tenant name '%s' is reserved", raw)
	}
	return TenantName(raw), nil
}

func TenantNameFromPrimitives(raw string) TenantName {
	return TenantName(raw)
}

func (n TenantName) String() string { return string(n) }
