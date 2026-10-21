package vos

import "fmt"

type UserName string

func NewUserName(raw string) (UserName, error) {
	if raw == "" {
		return "", fmt.Errorf("admin name cannot be empty")
	}
	return UserName(raw), nil
}

func UserNameFromPrimitives(raw string) UserName {
	return UserName(raw)
}

func (n UserName) String() string { return string(n) }
