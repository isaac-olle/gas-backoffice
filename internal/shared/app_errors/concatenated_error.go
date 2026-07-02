package app_errors

import (
	"encoding/json"
	"errors"
)

var ErrSeveralErrorsBuildingEntity = errors.New("several errors building entity")

type ConcatenatedError []error

func NewEmptyErrorArray() *ConcatenatedError {
	return &ConcatenatedError{}
}

func (e *ConcatenatedError) Error() string {
	if len(*e) == 0 {
		return "[]"
	}
	errorsSlice := make([]interface{}, 0, len(*e))
	for _, err := range *e {
		errorsSlice = append(errorsSlice, err)
	}
	b, err := json.Marshal(errorsSlice)
	if err != nil {
		return `["failed to marshal errors"]`
	}
	return string(b)
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
