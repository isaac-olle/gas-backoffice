package authors

import (
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type UserAuthor struct {
	Id   ids2.UserId
	Role vo2.Role
}

func NewUserAuthorFromString(id string, role vo2.Role) (*UserAuthor, error) {
	userId, err := ids2.NewUserIdFromString(id)
	if err != nil {
		return nil, err
	}
	return &UserAuthor{Id: userId, Role: role}, nil
}

func (a *UserAuthor) EqualsRole(other vo2.Role) bool {
	return a.Role.Equals(other)
}

func (a *UserAuthor) ToPrimitives() *vo2.AuthorPrimitives {
	return &vo2.AuthorPrimitives{
		Id:   a.Id.String(),
		Role: a.Role.String(),
	}
}
