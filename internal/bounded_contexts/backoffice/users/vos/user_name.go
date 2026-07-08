package vos

import "fmt"

type UserName string

var forbiddenUserNames = map[string]struct{}{
	"Gas-User": {},
}

func NewUserName(raw string) (UserName, error) {
	if raw == "" {
		return "", fmt.Errorf("user name cannot be empty")
	}
	if _, ok := forbiddenUserNames[raw]; ok {
		return "", fmt.Errorf("user name '%s' is reserved", raw)
	}
	return UserName(raw), nil
}

func UserNameFromPrimitives(raw string) UserName {
	return UserName(raw)
}

func (n UserName) String() string { return string(n) }
