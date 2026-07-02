package vo

import (
	"fmt"
)

type Role string

const (
	TenantRole Role = "tenant"
	AdminRole  Role = "admin"
	ClientRole Role = "client"
)

var allowedRoles = map[string]Role{
	string(TenantRole): TenantRole,
	string(AdminRole):  AdminRole,
	string(ClientRole): ClientRole,
}

func NewRole(s string) (Role, error) {
	role, ok := allowedRoles[s]
	if !ok {
		return "", fmt.Errorf("invalid role: %s", s)
	}
	return role, nil
}

func (r Role) IsAdmin() bool          { return r == AdminRole }
func (r Role) IsTenant() bool         { return r == TenantRole }
func (r Role) IsClient() bool         { return r == ClientRole }
func (r Role) String() string         { return string(r) }
func (r Role) Equals(other Role) bool { return r.String() == other.String() }
