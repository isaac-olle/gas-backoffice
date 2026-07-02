package vos

import "fmt"

type ClientName string

func NewClientName(raw string) (ClientName, error) {
	if raw == "" {
		return "", fmt.Errorf("client name cannot be empty")
	}
	return ClientName(raw), nil
}

func (n ClientName) String() string { return string(n) }
