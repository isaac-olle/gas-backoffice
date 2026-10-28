package vo

import (
	"encoding/json"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// SubscriptionRole controls what a user may do within one subscription. It is
// intentionally separate from Role, which identifies the global user type.
type SubscriptionRole string

const UserSubscriptionRole SubscriptionRole = "user"

func NewSubscriptionRole(raw string) (SubscriptionRole, error) {
	role := SubscriptionRole(raw)
	switch role {
	case UserSubscriptionRole:
		return role, nil
	default:
		return "", fmt.Errorf("invalid subscription role: %s", raw)
	}
}

func (r SubscriptionRole) String() string { return string(r) }

func (r SubscriptionRole) Equals(other SubscriptionRole) bool { return r == other }

type SubscriptionMembership struct {
	SubscriptionId   ids.SubscriptionId `json:"subscription_id"`
	SubscriptionRole SubscriptionRole   `json:"subscription_role"`
}

func NewSubscriptionMembership(subscriptionId ids.SubscriptionId, role SubscriptionRole) SubscriptionMembership {
	return SubscriptionMembership{SubscriptionId: subscriptionId, SubscriptionRole: role}
}

func (m SubscriptionMembership) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SubscriptionId   string `json:"subscription_id"`
		SubscriptionRole string `json:"subscription_role"`
	}{SubscriptionId: m.SubscriptionId.String(), SubscriptionRole: m.SubscriptionRole.String()})
}

func (m *SubscriptionMembership) UnmarshalJSON(data []byte) error {
	var raw struct {
		SubscriptionId   string `json:"subscription_id"`
		SubscriptionRole string `json:"subscription_role"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	id, err := ids.NewSubscriptionIdFromString(raw.SubscriptionId)
	if err != nil {
		return err
	}
	role, err := NewSubscriptionRole(raw.SubscriptionRole)
	if err != nil {
		return err
	}
	*m = NewSubscriptionMembership(id, role)
	return nil
}
