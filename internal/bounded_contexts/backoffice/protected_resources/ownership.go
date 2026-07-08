package protected_resources

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	UserId  ids2.UserId
	AdminId ids2.AdminId
}

func NewOwnership(userId ids2.UserId, adminId ids2.AdminId) Ownership {
	return Ownership{UserId: userId, AdminId: adminId}
}

func NewOwnershipFromStrings(userId string, adminId string) (Ownership, error) {
	tId, err := ids2.NewUserIdFromString(userId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid user id").
			WithCause(err)
	}
	aId, err := ids2.NewAdminIdFromString(adminId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid admin id").
			WithCause(err)
	}
	return NewOwnership(tId, aId), nil
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}
