package authors

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type UserAuthor struct {
	Id      ids2.UserId
	AdminId ids2.AdminId
}

func NewUserAuthorFromString(id string) (*UserAuthor, error) {
	userId, err := ids2.NewUserIdFromString(id)
	if err != nil {
		return nil, err
	}
	return &UserAuthor{Id: userId}, nil
}

func (a *UserAuthor) EqualsRole(other vo2.Role) bool {
	return other.Equals(vo2.UserRole)
}

func (a *UserAuthor) ToPrimitives() *vo2.AuthorPrimitives {
	return &vo2.AuthorPrimitives{
		Id:   a.Id.String(),
		Role: vo2.UserRole.String(),
	}
}

func (a *UserAuthor) Ownership() users.Ownership {
	return users.NewOwnership(a.AdminId)
}
