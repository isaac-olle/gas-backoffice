package client_secrets

import (
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

var clientSecretAggregateName aggregate2.Name = aggregate2.ClientSecret

type ClientSecret struct {
	*aggregate2.Root
	id         ids2.ClientSecretId
	ownership  Ownership
	secretHash string
}

func NewClientSecret(id string, clientId ids2.ClientId, secretHash string) (*ClientSecret, error) {
	errArr := app_errors.NewEmptyErrorArray()

	csId, err := ids2.NewClientSecretIdFromString(id)
	if err != nil {
		errArr.AddError(err)
	}

	if secretHash == "" {
		errArr.AddError(fmt.Errorf("secret hash cannot be empty"))
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	cs := &ClientSecret{
		Root:       aggregate2.NewRoot(clientSecretAggregateName),
		id:         csId,
		ownership:  NewOwnershipFromVos(clientId),
		secretHash: secretHash,
	}
	cs.Record(securityevents.NewClientSecretCreated(cs.id.String(), cs.ownership.ClientId.String()))
	return cs, nil
}

func FromPrimitives(p ClientSecretPrimitives) (*ClientSecret, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewClientSecretIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	ownership, err := NewOwnership(p.ClientId)
	if err != nil {
		errArr.AddError(err)
	}

	if p.SecretHash == "" {
		errArr.AddError(fmt.Errorf("secret hash cannot be empty"))
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	return &ClientSecret{
		Root:       aggregate2.NewRoot(clientSecretAggregateName),
		id:         id,
		ownership:  ownership,
		secretHash: p.SecretHash,
	}, nil
}

func (cs *ClientSecret) CheckOwnership(o Ownership) bool { return cs.ownership == o }
