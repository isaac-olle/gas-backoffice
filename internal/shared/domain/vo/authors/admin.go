package authors

import (
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type AdminAuthor struct {
	Id ids.AdminId
}

func NewAdminAuthorFromString(id string) (*AdminAuthor, error) {
	adminId, err := ids.NewAdminIdFromString(id)
	if err != nil {
		return nil, err
	}
	return &AdminAuthor{Id: adminId}, nil
}

func (a *AdminAuthor) EqualsRole(other vo2.Role) bool {
	return other.Equals(vo2.AdminRole)
}

func (a *AdminAuthor) ToPrimitives() *vo2.AuthorPrimitives {
	return &vo2.AuthorPrimitives{
		Id:   a.Id.String(),
		Role: vo2.AdminRole.String(),
	}
}
