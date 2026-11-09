package oidc_test

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure/helpers/oidc"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCodeVerifier(t *testing.T) {
	t.Parallel()

	t.Run("produces a verifier the value object accepts", func(t *testing.T) {
		t.Parallel()
		verifier, err := oidc.NewRandomCodeVerifierGenerator().GenerateCodeVerifier()
		require.NoError(t, err)

		// 32 bytes en base64url són exactament el mínim de l'RFC.
		assert.Len(t, verifier.String(), vos.MinCodeVerifierLength)

		roundTripped, err := vos.NewCodeVerifierFromString(verifier.String())
		require.NoError(t, err)
		assert.True(t, verifier.Equals(roundTripped))
	})

	t.Run("does not repeat itself across calls", func(t *testing.T) {
		t.Parallel()
		generator := oidc.NewRandomCodeVerifierGenerator()

		first, err := generator.GenerateCodeVerifier()
		require.NoError(t, err)
		second, err := generator.GenerateCodeVerifier()
		require.NoError(t, err)

		assert.False(t, first.Equals(second))
	})
}
