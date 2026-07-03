package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

func TestDetermineFlow_Code(t *testing.T) {
	t.Parallel()
	f, err := sessionvos.DetermineFlow("code")
	require.NoError(t, err)
	assert.Equal(t, sessionvos.AuthorizationCodeFlow, f.String())
}

func TestDetermineFlow_Token(t *testing.T) {
	t.Parallel()
	f, err := sessionvos.DetermineFlow("token")
	require.NoError(t, err)
	assert.Equal(t, sessionvos.ClientCredentialsFlow, f.String())
}

func TestDetermineFlow_Unknown(t *testing.T) {
	t.Parallel()
	_, err := sessionvos.DetermineFlow("unknown_type")
	assert.Error(t, err)
}

func TestNewFlowFromString_Valid(t *testing.T) {
	t.Parallel()
	for _, v := range []string{sessionvos.AuthorizationCodeFlow, sessionvos.ClientCredentialsFlow} {
		f, err := sessionvos.NewFlowFromString(v)
		require.NoError(t, err)
		assert.Equal(t, v, f.String())
	}
}

func TestNewFlowFromString_Invalid(t *testing.T) {
	t.Parallel()
	_, err := sessionvos.NewFlowFromString("bogus")
	assert.Error(t, err)
}
