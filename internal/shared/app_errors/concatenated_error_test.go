package app_errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptyErrorArray(t *testing.T) {
	t.Parallel()
	ce := NewEmptyErrorArray()
	assert.NotNil(t, ce)
	assert.True(t, ce.IsEmpty())
	assert.False(t, ce.HasErrors())
}

func TestConcatenatedError_AddError(t *testing.T) {
	t.Parallel()
	ce := NewEmptyErrorArray()
	ce.AddError(errors.New("err1"))
	ce.AddError(errors.New("err2"))

	assert.False(t, ce.IsEmpty())
	assert.True(t, ce.HasErrors())
	assert.Len(t, ce.Errors(), 2)
}

func TestConcatenatedError_Error_Empty(t *testing.T) {
	t.Parallel()
	ce := NewEmptyErrorArray()
	assert.Equal(t, "[]", ce.Error())
}

func TestConcatenatedError_Error_WithErrors(t *testing.T) {
	t.Parallel()
	ce := NewEmptyErrorArray()
	ce.AddError(errors.New("first error"))
	s := ce.Error()
	assert.NotEmpty(t, s)
	assert.NotEqual(t, "[]", s)
}

func TestConcatenatedError_Errors_Slice(t *testing.T) {
	t.Parallel()
	ce := NewEmptyErrorArray()
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	ce.AddError(e1)
	ce.AddError(e2)

	errs := ce.Errors()
	assert.Equal(t, e1, errs[0])
	assert.Equal(t, e2, errs[1])
}
