package vos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

func TestNewGrantTypeFromString_Valid(t *testing.T) {
	t.Parallel()
	for _, v := range []sessionvos.GrantType{sessionvos.AuthorizationCodeGrantType, sessionvos.ClientCredentialsGrantType} {
		g, err := sessionvos.NewGrantTypeFromString(v.String())
		require.NoError(t, err)
		assert.Equal(t, v, g)
	}
}

func TestNewGrantTypeFromString_Invalid(t *testing.T) {
	t.Parallel()
	_, err := sessionvos.NewGrantTypeFromString("bogus")
	assert.Error(t, err)
}
