package clients

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	UserId ids2.UserId
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}

func NewOwnership(userId ids2.UserId) Ownership {
	return Ownership{UserId: userId}
}
