package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

func TestNewSessionState_ValidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    string
		expected sessionvos.SessionState
	}{
		{"active", sessionvos.SessionStateActive},
		{"authorizing", sessionvos.Authorizing},
		{"consumed", sessionvos.SessionStateConsumed},
		{"expired", sessionvos.SessionStateExpired},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			s, err := sessionvos.NewSessionState(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, s)
			assert.Equal(t, tc.input, s.String())
		})
	}
}

func TestNewSessionState_Invalid(t *testing.T) {
	t.Parallel()
	_, err := sessionvos.NewSessionState("unknown_state")
	assert.Error(t, err)
}
