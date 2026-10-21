package vos

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

const stateByteLength = 32

type State struct{ value string }

func NewState(v string) (State, error) { return State{value: v}, nil }

// GenerateState produces a cryptographically random, URL-safe state value.
func GenerateState() (State, error) {
	b := make([]byte, stateByteLength)
	if _, err := rand.Read(b); err != nil {
		return State{}, app_errors.InternalError("STATE_GENERATION_FAILED").
			WithDescription("failed to generate secure state").WithCause(

			err, false)

	}
	return State{value: base64.RawURLEncoding.EncodeToString(b)}, nil
}

func (s State) Value() string  { return s.value }
func (s State) String() string { return s.value }
