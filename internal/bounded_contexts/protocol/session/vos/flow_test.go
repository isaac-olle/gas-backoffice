package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

func TestNewFlowFromString_Valid(t *testing.T) {
	t.Parallel()
	for _, v := range []sessionvos.Flow{sessionvos.AuthorizationCodeFlow, sessionvos.ClientCredentialsFlow} {
		f, err := sessionvos.NewFlowFromString(v.String())
		require.NoError(t, err)
		assert.Equal(t, v, f)
	}
}

func TestNewFlowFromString_Invalid(t *testing.T) {
	t.Parallel()
	_, err := sessionvos.NewFlowFromString("bogus")
	assert.Error(t, err)
}
