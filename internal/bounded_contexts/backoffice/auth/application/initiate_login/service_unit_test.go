package initiate_login_test

import (
	"context"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/initiate_login"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/stretchr/testify/require"
)

func TestInitiate(t *testing.T) {
	tests := []struct {
		name          string
		referenceName string
		repoErr       error
		connectorErr  error
		wantURL       string
		wantCode      app_errors.ErrorCode
	}{
		{name: "connects to the configured provider", referenceName: "google", wantURL: "https://idp.example/authorize"},
		{name: "rejects an invalid provider name", referenceName: "not valid", wantCode: auth.ErrInvalidIdpReferenceName},
		{name: "wraps provider lookup failures", referenceName: "google", repoErr: errInitiateLogin, wantCode: auth.ErrInitiateLoginFailed},
		{name: "wraps connector failures", referenceName: "google", connectorErr: errInitiateLogin, wantCode: auth.ErrInitiateLoginFailed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &idpConfigRepositoryStub{idp: validIdp(), err: tt.repoErr}
			connector := &idpConnectorStub{url: "https://idp.example/authorize", err: tt.connectorErr}
			got, err := initiate_login.NewService(repo, connector).Initiate(context.Background(), initiate_login.NewQuery(tt.referenceName))
			if tt.wantCode != "" {
				require.Error(t, err)
				var appErr *app_errors.AppError
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, tt.wantCode, appErr.Code)
				return
			}
			require.NoError(t, err)
			result, ok := got.(initiate_login.Result)
			require.True(t, ok)
			require.Equal(t, tt.wantURL, result.AuthorizationURL)
			require.Same(t, repo.idp, connector.got)
			require.Equal(t, "google", repo.gotName.String())
		})
	}
}
