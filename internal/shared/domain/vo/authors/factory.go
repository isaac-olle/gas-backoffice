package authors

import (
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func NewAuthorForRole(role vo.Role, sub string) (vo.Author, error) {
	switch {
	case role.IsAdmin(), role.IsTenant():
		userId, err := ids.NewUserIdFromString(sub)
		if err != nil {
			return nil, err
		}
		return &UserAuthor{Id: userId, Role: role}, nil

	default:
		return nil, fmt.Errorf("unknown role: %s", role)
	}
}
