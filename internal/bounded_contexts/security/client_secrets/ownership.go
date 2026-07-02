package client_secrets

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type Ownership struct {
	ClientId ids.ClientId
}

func newOwnership(clientId ids.ClientId) Ownership {
	return Ownership{ClientId: clientId}
}

func NewOwnershipFromVos(clientId ids.ClientId) Ownership {
	return newOwnership(clientId)
}

func NewOwnership(clientId string) (Ownership, error) {
	cId, err := ids.NewClientIdFromString(clientId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid client id").
			WithCause(err)
	}
	return newOwnership(cId), nil
}
