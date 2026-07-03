package security

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const AdminAssertionKeyCreatedType event.Type = "security.admin_assertion_key.created"

type AdminAssertionKeyCreated struct {
	AdminId string
	event.BaseEvent
}

func (a AdminAssertionKeyCreated) Type() event.Type {
	return AdminAssertionKeyCreatedType
}

func NewAdminAssertionKeyCreated(aggregateId, adminId string) AdminAssertionKeyCreated {
	return AdminAssertionKeyCreated{
		AdminId: adminId,
		BaseEvent: event.NewEvent(
			aggregateId,
			AdminAssertionKeyCreatedType,
			false,
			map[string]any{"admin_id": adminId},
			aggregate.AdminAssertionKey,
		),
	}
}
