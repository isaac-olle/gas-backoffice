package authors

import (
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ClientAuthor struct {
	Id ids.ClientId
}

func NewClientAuthorFromString(id string) (*ClientAuthor, error) {
	clientId, err := ids.NewClientIdFromString(id)
	if err != nil {
		return nil, err
	}
	return &ClientAuthor{Id: clientId}, nil
}

func (a *ClientAuthor) EqualsRole(other vo2.Role) bool {
	return other.Equals(vo2.ClientRole)
}

func (a *ClientAuthor) ToPrimitives() *vo2.AuthorPrimitives {
	return &vo2.AuthorPrimitives{
		Id:   a.Id.String(),
		Role: vo2.ClientRole.String(),
	}
}
