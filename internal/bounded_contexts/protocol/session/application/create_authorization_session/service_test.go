package create_authorization_session_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session"
	create_authorization_session "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/create_authorization_session"
	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/domain_services"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestCreateAuthorizationSession_HappyPath(t *testing.T) {
	t.Parallel()

	clientID := uuid.New().String()
	sessionID := ids.NewSessionId()

	sessionRepo := &mockSessionRepo{}
	clientRepo := &mockClientRepo{client: newClientMother(clientID)}
	clientFinder := domain_services.NewClientFinder(clientRepo)
	bus := &mockEventBus{}
	gen := &mockCodeGen{code: "authcode"}

	svc := create_authorization_session.NewService(sessionRepo, clientFinder, bus, gen)

	cmd := create_authorization_session.NewCommand(
		sessionID, "code",
		clientID, "https://example.com/callback", "openid", "state123",
		"", "",
	)

	require.NoError(t, svc.Create(context.Background(), cmd))

	require.Len(t, sessionRepo.saved, 1)
	p := sessionRepo.saved[0].ToPrimitives()
	assert.Equal(t, clientID, p.ClientID)
	assert.Equal(t, sessionID.String(), p.SessionID)
	assert.Equal(t, "openid", p.Scope)
	assert.Equal(t, "authcode", p.Code)
}

func TestCreateAuthorizationSession_ErrorCases(t *testing.T) {
	t.Parallel()

	clientID := uuid.New().String()

	tests := []struct {
		name         string
		clientID     string
		redirectURI  string
		scope        string
		clientResult *session.Client
		setupRepo    func(repo *mockSessionRepo)
		wantErrCode  app_errors.ErrorCode
	}{
		{
			name:         "returns not found error when client does not exist",
			clientID:     uuid.New().String(),
			redirectURI:  "https://example.com/callback",
			scope:        "openid",
			clientResult: nil,
			wantErrCode:  session.ErrClientNotFound,
		},
		{
			name:         "returns invalid parameters error for empty scope",
			clientID:     clientID,
			redirectURI:  "https://example.com/callback",
			scope:        "",
			clientResult: newClientMother(clientID),
			wantErrCode:  session.ErrInvalidSessionParameters,
		},
		{
			name:        "returns invalid redirect uri error for mismatched uri",
			clientID:    clientID,
			redirectURI: "https://example.com/callback",
			scope:       "openid",
			clientResult: session.NewClient(
				clientID, "TC",
				uuid.New().String(), uuid.New().String(),
				[]string{"authorization_code"},
				[]string{"https://other.com/callback"},
				[]*session.Scope{{ScopeName: "openid"}},
			),
			wantErrCode: session.ErrInvalidRedirectURI,
		},
		{
			name:         "wraps repository save error",
			clientID:     clientID,
			redirectURI:  "https://example.com/callback",
			scope:        "openid",
			clientResult: newClientMother(clientID),
			setupRepo: func(repo *mockSessionRepo) {
				repo.err = errors.New("db error")
			},
			wantErrCode: session.ErrCannotSaveSession,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			sessionRepo := &mockSessionRepo{}
			if tc.setupRepo != nil {
				tc.setupRepo(sessionRepo)
			}
			clientRepo := &mockClientRepo{client: tc.clientResult}
			clientFinder := domain_services.NewClientFinder(clientRepo)
			bus := &mockEventBus{}
			gen := &mockCodeGen{code: "authcode"}

			svc := create_authorization_session.NewService(sessionRepo, clientFinder, bus, gen)

			cmd := create_authorization_session.NewCommand(
				ids.NewSessionId(), "code",
				tc.clientID, tc.redirectURI, tc.scope, "state123",
				"", "",
			)

			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
			assert.False(t, bus.publishCalled)
		})
	}
}
