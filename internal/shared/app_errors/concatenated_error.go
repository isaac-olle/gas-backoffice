package app_errors

import (
	"errors"
	"strings"
)

var ErrSeveralErrorsBuildingEntity = errors.New("several errors building entity")

type ConcatenatedError []error

func NewEmptyErrorArray() *ConcatenatedError {
	return &ConcatenatedError{}
}

func (e *ConcatenatedError) Error() string {
	messages := make([]string, 0, len(*e))
	for _, err := range *e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func (e *ConcatenatedError) AddError(err error) {
	*e = append(*e, err)
}

func (e *ConcatenatedError) HasErrors() bool {
	return len(*e) > 0
}

func (e *ConcatenatedError) IsEmpty() bool {
	return len(*e) == 0
}

func (e *ConcatenatedError) Errors() []error {
	return *e
}
