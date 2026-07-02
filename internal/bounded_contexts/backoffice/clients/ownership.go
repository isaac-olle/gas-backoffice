package clients

import (
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type Ownership struct {
	AdminId  ids2.AdminId
	TenantId ids2.TenantId
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}

func NewOwnership(adminId ids2.AdminId, tenantId ids2.TenantId) Ownership {
	return Ownership{AdminId: adminId, TenantId: tenantId}
}
