package vos_test

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/vos"
)

func TestNewPkce_NilWhenBothEmpty(t *testing.T) {
	t.Parallel()
	p, err := sessionvos.NewPkce("", "")
	require.NoError(t, err)
	assert.Nil(t, p)
}

func TestNewPkce_PlainMethod(t *testing.T) {
	t.Parallel()
	p, err := sessionvos.NewPkce("mychallenge", "plain")
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, "mychallenge", p.Challenge())
	assert.Equal(t, sessionvos.CodeChallengePlain, p.Method())
}

func TestNewPkce_S256Method(t *testing.T) {
	t.Parallel()
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	h := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(h[:])

	p, err := sessionvos.NewPkce(challenge, "S256")
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.True(t, p.Verify(verifier))
}

func TestNewPkce_InvalidMethod(t *testing.T) {
	t.Parallel()
	_, err := sessionvos.NewPkce("challenge", "MD5")
	assert.Error(t, err)
}

func TestNewPkce_ChallengeMissingMethod(t *testing.T) {
	t.Parallel()
	_, err := sessionvos.NewPkce("challenge", "")
	assert.Error(t, err)
}

func TestNewPkce_MethodMissingChallenge(t *testing.T) {
	t.Parallel()
	_, err := sessionvos.NewPkce("", "S256")
	assert.Error(t, err)
}

func TestPkce_Verify_Plain_Match(t *testing.T) {
	t.Parallel()
	p, _ := sessionvos.NewPkce("secret", "plain")
	assert.True(t, p.Verify("secret"))
	assert.False(t, p.Verify("wrong"))
}

func TestPkce_Verify_S256_NoMatch(t *testing.T) {
	t.Parallel()
	p, _ := sessionvos.NewPkce("challenge", "S256")
	assert.False(t, p.Verify("wrong-verifier"))
}

func TestNewCodeChallengeMethod(t *testing.T) {
	t.Parallel()
	t.Run("accepts plain", func(t *testing.T) {
		m, err := sessionvos.NewCodeChallengeMethod("plain")
		require.NoError(t, err)
		assert.Equal(t, sessionvos.CodeChallengePlain, m)
	})
	t.Run("accepts S256", func(t *testing.T) {
		m, err := sessionvos.NewCodeChallengeMethod("S256")
		require.NoError(t, err)
		assert.Equal(t, sessionvos.CodeChallengeS256, m)
	})
	t.Run("rejects unknown method", func(t *testing.T) {
		_, err := sessionvos.NewCodeChallengeMethod("MD5")
		assert.Error(t, err)
	})
}

func TestCodeChallengeMethod_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "S256", sessionvos.CodeChallengeS256.String())
	assert.Equal(t, "plain", sessionvos.CodeChallengePlain.String())
}
