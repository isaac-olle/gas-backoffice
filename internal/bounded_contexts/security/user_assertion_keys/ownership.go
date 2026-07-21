package user_assertion_keys

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	UserId ids.UserId
}

func newOwnership(userId ids.UserId) Ownership {
	return Ownership{UserId: userId}
}

func NewOwnershipFromVos(userId ids.UserId) Ownership {
	return newOwnership(userId)
}

func NewOwnership(userId string) (Ownership, error) {
	uId, err := ids.NewUserIdFromString(userId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid user id").
			WithCause(err)
	}
	return newOwnership(uId), nil
}
