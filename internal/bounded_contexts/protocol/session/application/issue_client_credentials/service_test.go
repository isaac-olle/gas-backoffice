package issue_client_credentials_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/issue_client_credentials"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestIssueClientCredentials_HappyPath(t *testing.T) {
	t.Parallel()

	clientID := uuid.New().String()
	client := newClientWithCCMother(clientID)

	clientRepo := &mockClientRepo{client: client}
	tokenGen := &mockTokenGen{token: "access-token"}
	bus := &mockEventBus{}

	clientFinder := domain_services.NewClientFinder(clientRepo)
	svc := issue_client_credentials.NewService(clientFinder, tokenGen, bus)

	qry := issue_client_credentials.NewQuery(clientID, "openid")
	result, err := svc.IssueClientCredentials(context.Background(), qry)
	require.NoError(t, err)

	tokenResp, ok := result.(*session.TokenResponse)
	require.True(t, ok)
	assert.Equal(t, "access-token", tokenResp.Token)

	require.Len(t, bus.published, 1)
}

func TestIssueClientCredentials_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		buildSetup  func(clientID string) (*mockClientRepo, *mockTokenGen)
		scope       string
		wantErrCode app_errors.ErrorCode
	}{
		{
			name: "returns not found error when client does not exist",
			buildSetup: func(clientID string) (*mockClientRepo, *mockTokenGen) {
				return &mockClientRepo{client: nil}, &mockTokenGen{}
			},
			scope:       "openid",
			wantErrCode: session.ErrClientNotFound,
		},
		{
			name: "returns flow not allowed error for client without client_credentials",
			buildSetup: func(clientID string) (*mockClientRepo, *mockTokenGen) {
				c := session.NewClient(
					clientID, "TC",
					uuid.New().String(), uuid.New().String(),
					[]string{"authorization_code"},
					[]string{},
					[]*session.Scope{{ScopeName: "openid"}},
				)
				return &mockClientRepo{client: c}, &mockTokenGen{}
			},
			scope:       "openid",
			wantErrCode: session.ErrFlowNotAllowed,
		},
		{
			name: "returns invalid scope error when scope is not registered",
			buildSetup: func(clientID string) (*mockClientRepo, *mockTokenGen) {
				return &mockClientRepo{client: newClientWithCCMother(clientID)}, &mockTokenGen{}
			},
			scope:       "admin",
			wantErrCode: session.ErrInvalidScope,
		},
		{
			name: "returns error when token generator fails",
			buildSetup: func(clientID string) (*mockClientRepo, *mockTokenGen) {
				return &mockClientRepo{client: newClientWithCCMother(clientID)},
					&mockTokenGen{err: errors.New("sign error")}
			},
			scope:       "openid",
			wantErrCode: session.ErrCannotSignToken,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clientID := uuid.New().String()
			clientRepo, tokenGen := tc.buildSetup(clientID)
			bus := &mockEventBus{}

			clientFinder := domain_services.NewClientFinder(clientRepo)
			svc := issue_client_credentials.NewService(clientFinder, tokenGen, bus)

			qry := issue_client_credentials.NewQuery(clientID, tc.scope)
			_, err := svc.IssueClientCredentials(context.Background(), qry)

			assertErrCode(t, err, tc.wantErrCode)
			assert.Equal(t, false, bus.publishCalled)
		})
	}
}

func TestIssueClientCredentials_NewSessionError(t *testing.T) {
	t.Parallel()

	clientID := uuid.New().String()
	client := session.NewClient(
		clientID, "TC",
		"invalid-user-id", "invalid-admin-id",
		[]string{"client_credentials"}, []string{},
		[]*session.Scope{{ScopeName: "openid"}},
	)
	clientRepo := &mockClientRepo{client: client}
	tokenGen := &mockTokenGen{}
	bus := &mockEventBus{}

	clientFinder := domain_services.NewClientFinder(clientRepo)
	svc := issue_client_credentials.NewService(clientFinder, tokenGen, bus)

	qry := issue_client_credentials.NewQuery(clientID, "openid")
	_, err := svc.IssueClientCredentials(context.Background(), qry)
	assert.Error(t, err)
}
