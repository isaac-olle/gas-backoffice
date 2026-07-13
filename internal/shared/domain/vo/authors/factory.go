package authors

import (
	"fmt"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func NewAuthorForRole(role vo2.Role, sub string) (vo2.Author, error) {
	switch {
	case role.IsAdmin():
		adminId, err := ids2.NewAdminIdFromString(sub)
		if err != nil {
			return nil, err
		}
		return &AdminAuthor{Id: adminId}, nil

	case role.IsUser():
		userId, err := ids2.NewUserIdFromString(sub)
		if err != nil {
			return nil, err
		}
		return &UserAuthor{Id: userId}, nil

	case role.IsClient():
		clientId, err := ids2.NewClientIdFromString(sub)
		if err != nil {
			return nil, err
		}
		return &ClientAuthor{Id: clientId}, nil

	default:
		return nil, fmt.Errorf("unknown role: %s", role)
	}
}
