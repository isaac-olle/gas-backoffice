package protected_resources

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	UserId ids2.UserId
}

func NewOwnership(userId ids2.UserId) Ownership {
	return Ownership{UserId: userId}
}

func NewOwnershipFromStrings(userId string) (Ownership, error) {
	tId, err := ids2.NewUserIdFromString(userId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid user id").
			WithCause(err)
	}
	return NewOwnership(tId), nil
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}
