package scopes

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	ProtectedResourceId ids2.ProtectedResourceId
	UserId              ids2.UserId
}

func newOwnership(protectedResourceId ids2.ProtectedResourceId, userId ids2.UserId) Ownership {
	return Ownership{ProtectedResourceId: protectedResourceId, UserId: userId}
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}

func NewOwnershipFromVos(protectedResourceId ids2.ProtectedResourceId, userId ids2.UserId) Ownership {
	return newOwnership(protectedResourceId, userId)
}

func NewOwnership(protectedResourceId, userId string) (Ownership, error) {
	prId, err := ids2.NewProtectedResourceIdFromString(protectedResourceId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid protected resource id").
			WithCause(err)
	}
	tId, err := ids2.NewUserIdFromString(userId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid user id").
			WithCause(err)
	}
	return newOwnership(prId, tId), nil
}
