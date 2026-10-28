package vo

import "fmt"

type Role string

const (
	SystemAdministratorRole Role = "system_administrator"
	SubscriptionUserRole    Role = "subscription_user"
)

const (
	SystemAdministratorClaim = "system_administrator"
	SubscriptionUserClaim    = "subscription_user"
)

var allowedRoles = map[string]Role{
	string(SubscriptionUserRole):    SubscriptionUserRole,
	string(SystemAdministratorRole): SystemAdministratorRole,
}

func NewRole(s string) (Role, error) {
	role, ok := allowedRoles[s]
	if !ok {
		return "", fmt.Errorf("invalid role: %s", s)
	}
	return role, nil
}

func (r Role) IsSystemAdministrator() bool { return r == SystemAdministratorRole }
func (r Role) IsSubscriptionUser() bool    { return r == SubscriptionUserRole }
func (r Role) String() string              { return string(r) }
func (r Role) Equals(other Role) bool      { return r.String() == other.String() }

// SystemRoleClaim is the public JWT representation of the global system role.
func (r Role) SystemRoleClaim() string {
	if r.IsSystemAdministrator() {
		return SystemAdministratorClaim
	}
	return SubscriptionUserClaim
}

func RoleFromSystemRoleClaim(s string) (Role, error) {
	switch s {
	case SystemAdministratorClaim:
		return SystemAdministratorRole, nil
	case SubscriptionUserClaim:
		return SubscriptionUserRole, nil
	default:
		return "", fmt.Errorf("invalid system role: %s", s)
	}
}
