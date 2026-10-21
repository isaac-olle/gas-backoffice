package vos_test

import (
	"strings"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCodeVerifierFromString(t *testing.T) {
	t.Parallel()

	t.Run("accepts values within the RFC 7636 shape", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			name string
			raw  string
		}{
			{name: "minimum length", raw: strings.Repeat("a", 43)},
			{name: "maximum length", raw: strings.Repeat("a", 128)},
			{name: "every unreserved character", raw: strings.Repeat("aA0-._~", 7)},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				verifier, err := vos.NewCodeVerifierFromString(tc.raw)
				require.NoError(t, err)
				assert.Equal(t, tc.raw, verifier.String())
			})
		}
	})

	t.Run("rejects values outside the RFC 7636 shape", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			name string
			raw  string
		}{
			{name: "empty", raw: ""},
			{name: "one char short", raw: strings.Repeat("a", 42)},
			{name: "a bare uuid", raw: "0f8fad5b-d9cb-469f-a165-70867728950e"},
			{name: "one char too long", raw: strings.Repeat("a", 129)},
			{name: "reserved character", raw: strings.Repeat("a", 42) + "/"},
			{name: "whitespace", raw: strings.Repeat("a", 42) + " "},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				_, err := vos.NewCodeVerifierFromString(tc.raw)

				var appErr *app_errors.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, vos.ErrInvalidCodeVerifier, appErr.Code)
			})
		}
	})
}

func TestCodeVerifierEquals(t *testing.T) {
	t.Parallel()

	a, err := vos.NewCodeVerifierFromString(strings.Repeat("a", 43))
	require.NoError(t, err)
	sameAsA, err := vos.NewCodeVerifierFromString(strings.Repeat("a", 43))
	require.NoError(t, err)
	b, err := vos.NewCodeVerifierFromString(strings.Repeat("b", 43))
	require.NoError(t, err)

	assert.True(t, a.Equals(sameAsA))
	assert.False(t, a.Equals(b))
}
