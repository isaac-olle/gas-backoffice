package clients

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	AdminId ids2.AdminId
	UserId  ids2.UserId
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}

func NewOwnership(adminId ids2.AdminId, userId ids2.UserId) Ownership {
	return Ownership{AdminId: adminId, UserId: userId}
}
