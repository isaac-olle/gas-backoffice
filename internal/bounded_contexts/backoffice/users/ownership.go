package users

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	UserId ids.UserId
}

func NewOwnership(userId ids.UserId) Ownership {
	return Ownership{UserId: userId}
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}
