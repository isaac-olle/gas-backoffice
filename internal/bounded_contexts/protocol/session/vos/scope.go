package vos

import "errors"

type Scope struct{ value string }

func NewScope(v string) (Scope, error) {
	if v == "" {
		return Scope{}, errors.New("scope cannot be empty")
	}
	return Scope{value: v}, nil
}

func (s Scope) Value() string  { return s.value }
func (s Scope) String() string { return s.value }
