package authenticate_client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func TestService_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("delegates to the first method that can process", func(t *testing.T) {
		t.Parallel()
		skipped := &stubMethod{canProcess: false}
		matched := &stubMethod{canProcess: true, result: authenticate_client.Result{ClientID: "client-b"}}

		svc := authenticate_client.NewService(skipped, matched)

		result, err := svc.Authenticate(context.Background(), authenticate_client.Credentials{})

		require.NoError(t, err)
		assert.Equal(t, "client-b", result.ClientID)
	})

	t.Run("propagates the matched method's error", func(t *testing.T) {
		t.Parallel()
		matched := &stubMethod{canProcess: true, err: assert.AnError}

		svc := authenticate_client.NewService(matched)

		_, err := svc.Authenticate(context.Background(), authenticate_client.Credentials{})

		assert.Error(t, err)
	})

	t.Run("returns ErrNoAuthenticationMethod when no method matches", func(t *testing.T) {
		t.Parallel()
		svc := authenticate_client.NewService(&stubMethod{canProcess: false})

		_, err := svc.Authenticate(context.Background(), authenticate_client.Credentials{})

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, authenticate_client.ErrNoAuthenticationMethod, appErr.Code)
	})

	t.Run("returns ErrNoAuthenticationMethod with zero methods registered", func(t *testing.T) {
		t.Parallel()
		svc := authenticate_client.NewService()

		_, err := svc.Authenticate(context.Background(), authenticate_client.Credentials{})

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, authenticate_client.ErrNoAuthenticationMethod, appErr.Code)
	})
}
