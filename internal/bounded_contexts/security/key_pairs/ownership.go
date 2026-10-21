package key_pairs

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	ProtectedResourceId ids.ProtectedResourceId
}

func newOwnership(protectedResourceId ids.ProtectedResourceId) Ownership {
	return Ownership{ProtectedResourceId: protectedResourceId}
}

func NewOwnershipFromVos(protectedResourceId ids.ProtectedResourceId) Ownership {
	return newOwnership(protectedResourceId)
}

func NewOwnership(protectedResourceId string) (Ownership, error) {
	prId, err := ids.NewProtectedResourceIdFromString(protectedResourceId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid protected resource id").WithCause(

			err, false)

	}
	return newOwnership(prId), nil
}
