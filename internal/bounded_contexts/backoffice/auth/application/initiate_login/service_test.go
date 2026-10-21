package initiate_login_test

import (
	"context"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/initiate_login"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitiate(t *testing.T) {
	t.Parallel()

	t.Run("returns the authorization url built by the authorizer", func(t *testing.T) {
		t.Parallel()
		authorizer := &stubAuthorizer{url: "https://accounts.google.com/o/oauth2/v2/auth?state=abc"}
		service := initiate_login.NewService(
			&stubIdpConfigRepo{idp: anIdp()},
			authorizer,
			stubCodeVerifierGenerator{value: aValidVerifier()},
		)

		res, err := service.Initiate(context.Background(), initiate_login.NewQuery("google"))
		require.NoError(t, err)

		result, ok := res.(initiate_login.Result)
		require.True(t, ok)
		assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth?state=abc", result.AuthorizationURL)
	})

	t.Run("looks the idp up by its reference name", func(t *testing.T) {
		t.Parallel()
		repo := &stubIdpConfigRepo{idp: anIdp()}
		service := initiate_login.NewService(
			repo, &stubAuthorizer{url: "https://idp.local/authorize"},
			stubCodeVerifierGenerator{value: aValidVerifier()},
		)

		_, err := service.Initiate(context.Background(), initiate_login.NewQuery("google"))
		require.NoError(t, err)

		assert.Equal(t, "google", repo.gotReferenceName.String())
	})

	t.Run("hands the authorizer a login request carrying the generated verifier", func(t *testing.T) {
		t.Parallel()
		authorizer := &stubAuthorizer{url: "https://idp.local/authorize"}
		service := initiate_login.NewService(
			&stubIdpConfigRepo{idp: anIdp()},
			authorizer,
			stubCodeVerifierGenerator{value: aValidVerifier()},
		)

		_, err := service.Initiate(context.Background(), initiate_login.NewQuery("google"))
		require.NoError(t, err)

		require.NotNil(t, authorizer.gotLoginRequest)
		prims := authorizer.gotLoginRequest.ToPrimitives()
		assert.Equal(t, aValidVerifier(), prims.CodeVerifier)
		assert.Equal(t, "google", prims.IdpReferenceName)
		assert.NotEmpty(t, prims.Id)
	})

	t.Run("rejects an unknown idp reference name", func(t *testing.T) {
		t.Parallel()
		service := initiate_login.NewService(
			&stubIdpConfigRepo{idp: anIdp()},
			&stubAuthorizer{},
			stubCodeVerifierGenerator{value: aValidVerifier()},
		)

		_, err := service.Initiate(context.Background(), initiate_login.NewQuery("not-an-idp"))
		assertErrCode(t, err, auth.ErrInvalidIdpReferenceName)
	})

	t.Run("wraps idp lookup failure", func(t *testing.T) {
		t.Parallel()
		service := initiate_login.NewService(
			&stubIdpConfigRepo{err: errBoom},
			&stubAuthorizer{},
			stubCodeVerifierGenerator{value: aValidVerifier()},
		)

		_, err := service.Initiate(context.Background(), initiate_login.NewQuery("google"))
		assertErrCode(t, err, auth.ErrInitiateLoginFailed)
	})

	t.Run("wraps authorizer failure", func(t *testing.T) {
		t.Parallel()
		service := initiate_login.NewService(
			&stubIdpConfigRepo{idp: anIdp()},
			&stubAuthorizer{err: errBoom},
			stubCodeVerifierGenerator{value: aValidVerifier()},
		)

		_, err := service.Initiate(context.Background(), initiate_login.NewQuery("google"))
		assertErrCode(t, err, auth.ErrInitiateLoginFailed)
	})
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
