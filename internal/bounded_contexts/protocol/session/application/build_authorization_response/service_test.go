package build_authorization_response_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/build_authorization_response"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestBuildAuthorizationResponseService(t *testing.T) {
	t.Parallel()

	existingSess := newSessionMother(t)

	tests := []struct {
		name        string
		repoResult  *session.AuthorizationCodeSession
		repoErr     error
		wantErr     bool
		wantErrCode app_errors.ErrorCode
		checkResult func(t *testing.T, result any)
	}{
		{
			name:       "returns redirect uri with code and state for existing session",
			repoResult: existingSess,
			checkResult: func(t *testing.T, result any) {
				t.Helper()
				redirectURI, ok := result.(string)
				require.True(t, ok)
				assert.Contains(t, redirectURI, "https://example.com/cb")
				assert.Contains(t, redirectURI, "code=")
				assert.Contains(t, redirectURI, "state=mystate")
			},
		},
		{
			name:        "returns session not found error when session does not exist",
			repoResult:  nil,
			wantErr:     true,
			wantErrCode: session.ErrSessionNotFound,
		},
		{
			name:        "returns session not found error when repo fails",
			repoErr:     errors.New("timeout"),
			wantErr:     true,
			wantErrCode: session.ErrSessionNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &mockSessionRepo{byId: tc.repoResult, err: tc.repoErr}
			finder := domain_services.NewSessionFinder(repo)
			svc := build_authorization_response.NewService(finder)

			result, err := svc.BuildAuthorizationResponse(context.Background(), build_authorization_response.NewQuery(ids.NewSessionId()))

			if tc.wantErr {
				assertErrCode(t, err, tc.wantErrCode)
				return
			}
			require.NoError(t, err)
			if tc.checkResult != nil {
				tc.checkResult(t, result)
			}
		})
	}
}
