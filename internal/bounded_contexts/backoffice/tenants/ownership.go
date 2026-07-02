package tenants

import (
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type Ownership struct {
	AdminId ids.AdminId
}

func NewOwnership(adminId ids.AdminId) Ownership {
	return Ownership{AdminId: adminId}
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}
