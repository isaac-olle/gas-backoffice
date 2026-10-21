package app_errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCauseError_DelegatesAndFlags(t *testing.T) {
	t.Parallel()
	underlying := errors.New("boom")

	pub := CauseError{public: true, error: underlying}
	assert.Equal(t, "boom", pub.Error())
	assert.True(t, pub.IsPublic())
	assert.Equal(t, underlying, pub.Unwrap())

	priv := CauseError{public: false, error: underlying}
	assert.False(t, priv.IsPublic())
}

func TestCauseError_ErrorsIsThroughUnwrap(t *testing.T) {
	t.Parallel()
	sentinel := errors.New("sentinel")
	ce := CauseError{public: true, error: sentinel}
	assert.True(t, errors.Is(ce, sentinel))
}

func TestPublicCauses_OnlyPublicOnes(t *testing.T) {
	t.Parallel()
	err := InternalError(testCode).
		WithCause(errors.New("visible"), true).
		WithCause(errors.New("secret"), false)

	assert.Equal(t, []string{"visible"}, err.PublicCauses())
}

func TestPublicCauses_Empty(t *testing.T) {
	t.Parallel()
	err := InternalError(testCode).WithCause(errors.New("private"), false)
	assert.Empty(t, err.PublicCauses())
}
