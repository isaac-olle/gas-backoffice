package session_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	protocolEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/protocol"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestNewSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		clientID      string
		redirectURI   string
		scope         string
		responseType  string
		pkceChallenge string
		pkceMethod    string
		wantErr       bool
		wantErrCode   app_errors.ErrorCode
		wantGrantType string
	}{
		{
			name:          "creates code flow session",
			clientID:      uuid.New().String(),
			redirectURI:   "https://example.com/cb",
			scope:         "openid profile",
			responseType:  "code",
			wantGrantType: "authorization_code",
		},
		{
			name:          "creates authorization_code flow session when challenge is provided",
			clientID:      uuid.New().String(),
			redirectURI:   "https://example.com/cb",
			scope:         "openid",
			responseType:  "code",
			pkceChallenge: "challenge",
			pkceMethod:    "S256",
			wantGrantType: "authorization_code",
		},
		{
			name:         "rejects invalid client id",
			clientID:     "not-uuid",
			redirectURI:  "https://example.com/cb",
			scope:        "openid",
			responseType: "code",
			wantErr:      true,
			wantErrCode:  session.ErrInvalidSessionParameters,
		},
		{
			name:         "rejects non-url redirect uri",
			clientID:     uuid.New().String(),
			redirectURI:  "not-a-url",
			scope:        "openid",
			responseType: "code",
			wantErr:      true,
			wantErrCode:  session.ErrInvalidSessionParameters,
		},
		{
			name:         "rejects empty scope",
			clientID:     uuid.New().String(),
			redirectURI:  "https://example.com/cb",
			scope:        "",
			responseType: "code",
			wantErr:      true,
			wantErrCode:  session.ErrInvalidSessionParameters,
		},
		{
			name:         "rejects unsupported response type",
			clientID:     uuid.New().String(),
			redirectURI:  "https://example.com/cb",
			scope:        "openid",
			responseType: "magic",
			wantErr:      true,
			wantErrCode:  session.ErrInvalidSessionParameters,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var client *session.Client
			if !tc.wantErr {
				client = newTestClient(tc.clientID, tc.redirectURI, strings.Fields(tc.scope))
			}
			sess, err := session.NewAuthorizationCodeSession(
				ids.NewSessionId(), tc.clientID, tc.redirectURI, tc.scope,
				"state", tc.responseType, &mockCodeGenerator{code: "c"},
				tc.pkceChallenge, tc.pkceMethod,
				client, newUserIdentityMother(t),
			)
			if tc.wantErr {
				assertErrCode(t, err, tc.wantErrCode)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, sess)
			if tc.wantGrantType != "" {
				assert.Equal(t, tc.wantGrantType, sess.ToPrimitives().GrantType)
			}
		})
	}
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("round-trips through primitives", func(t *testing.T) {
		t.Parallel()
		orig := newSessionMother(t)
		prim := orig.ToPrimitives()

		restored, err := session.AuthorizationCodeSessionFromPrimitives(prim)
		require.NoError(t, err)
		assert.Equal(t, prim.SessionID, restored.ToPrimitives().SessionID)
	})

	t.Run("returns error for nil primitives", func(t *testing.T) {
		t.Parallel()
		_, err := session.AuthorizationCodeSessionFromPrimitives(nil)
		assertErrCode(t, err, session.ErrInvalidSessionParameters)
	})
}

func TestCheckExchangeParams(t *testing.T) {
	t.Parallel()

	t.Run("accepts correct client id and redirect uri", func(t *testing.T) {
		t.Parallel()
		s := newSessionMother(t)
		p := s.ToPrimitives()
		assert.NoError(t, s.CheckExchangeParams(p.ClientID, p.RedirectURI))
	})

	t.Run("rejects consumed session", func(t *testing.T) {
		t.Parallel()
		s := newSessionMother(t)
		p := s.ToPrimitives()
		s.Consume()
		assertErrCode(t, s.CheckExchangeParams(p.ClientID, p.RedirectURI), session.ErrSessionConsumed)
	})

	t.Run("rejects wrong client id", func(t *testing.T) {
		t.Parallel()
		s := newSessionMother(t)
		p := s.ToPrimitives()
		assertErrCode(t, s.CheckExchangeParams("other-client", p.RedirectURI), session.ErrCodeNotFound)
	})

	t.Run("rejects wrong redirect uri", func(t *testing.T) {
		t.Parallel()
		s := newSessionMother(t)
		p := s.ToPrimitives()
		assertErrCode(t, s.CheckExchangeParams(p.ClientID, "https://evil.com/cb"), session.ErrInvalidRedirectURI)
	})
}

func TestCheckPkce(t *testing.T) {
	t.Parallel()

	newPkceSession := func(t *testing.T) *session.AuthorizationCodeSession {
		t.Helper()
		clientID := uuid.New().String()
		redirectURI := "https://example.com/cb"
		client := newTestClient(clientID, redirectURI, []string{"openid"})
		sess, err := session.NewAuthorizationCodeSession(
			ids.NewSessionId(), clientID,
			redirectURI, "openid", "s",
			"code", &mockCodeGenerator{code: "c"}, "challenge", "plain",
			client, newUserIdentityMother(t),
		)
		require.NoError(t, err)
		return sess
	}

	tests := []struct {
		name        string
		pkce        bool
		verifier    string
		wantIsPkce  bool
		wantErr     bool
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:       "no pkce session returns false and no error",
			pkce:       false,
			wantIsPkce: false,
		},
		{
			name:        "pkce session with missing verifier returns error",
			pkce:        true,
			verifier:    "",
			wantIsPkce:  true,
			wantErr:     true,
			wantErrCode: session.ErrInvalidCodeVerifier,
		},
		{
			name:        "pkce session with wrong verifier returns error",
			pkce:        true,
			verifier:    "wrong-verifier",
			wantIsPkce:  true,
			wantErr:     true,
			wantErrCode: session.ErrInvalidCodeVerifier,
		},
		{
			name:       "pkce session with correct verifier returns true and no error",
			pkce:       true,
			verifier:   "challenge",
			wantIsPkce: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var sess *session.AuthorizationCodeSession
			if tc.pkce {
				sess = newPkceSession(t)
			} else {
				sess = newSessionMother(t)
			}
			isPkce, err := sess.CheckPkce(tc.verifier)
			assert.Equal(t, tc.wantIsPkce, isPkce)
			if tc.wantErr {
				assertErrCode(t, err, tc.wantErrCode)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckAgainstClient(t *testing.T) {
	t.Parallel()

	const (
		responseType = "code"
		redirectURI  = "https://example.com/callback"
		scope        = "openid"
	)

	tests := []struct {
		name          string
		grantTypes    []string
		uris          []string
		scopes        []string
		public        bool
		codeChallenge string
		wantErr       bool
		wantErrCode   app_errors.ErrorCode
	}{
		{
			name:       "accepts confidential client with matching grant type, uri, and scope",
			grantTypes: []string{"authorization_code"},
			uris:       []string{redirectURI},
			scopes:     []string{scope},
		},
		{
			name:        "rejects client without the authorization_code grant type",
			grantTypes:  []string{"client_credentials"},
			uris:        []string{redirectURI},
			scopes:      []string{scope},
			wantErr:     true,
			wantErrCode: session.ErrGrantTypeNotAllowed,
		},
		{
			name:        "rejects client with wrong redirect uri",
			grantTypes:  []string{"authorization_code"},
			uris:        []string{"https://other.com/cb"},
			scopes:      []string{scope},
			wantErr:     true,
			wantErrCode: session.ErrInvalidRedirectURI,
		},
		{
			name:        "rejects client without required scope",
			grantTypes:  []string{"authorization_code"},
			uris:        []string{redirectURI},
			scopes:      []string{"profile"},
			wantErr:     true,
			wantErrCode: session.ErrInvalidScope,
		},
		{
			name:          "accepts public client when code_challenge is present",
			grantTypes:    []string{"authorization_code"},
			uris:          []string{redirectURI},
			scopes:        []string{scope},
			public:        true,
			codeChallenge: "challenge-value",
		},
		{
			name:        "rejects public client when code_challenge is missing",
			grantTypes:  []string{"authorization_code"},
			uris:        []string{redirectURI},
			scopes:      []string{scope},
			public:      true,
			wantErr:     true,
			wantErrCode: session.ErrPkceRequired,
		},
		{
			name:          "accepts confidential client that also uses code_challenge",
			grantTypes:    []string{"authorization_code"},
			uris:          []string{redirectURI},
			scopes:        []string{scope},
			codeChallenge: "challenge-value",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			prId := uuid.New().String()
			scopeObjs := make([]*session.Scope, len(tc.scopes))
			for i, s := range tc.scopes {
				scopeObjs[i] = session.NewScope(uuid.New().String(), s, prId, "Resource A")
			}
			client := session.NewClient(
				uuid.New().String(), "TC",
				uuid.New().String(),
				tc.grantTypes, tc.uris, scopeObjs,
				tc.public,
			)

			codeChallengeMethod := ""
			if tc.codeChallenge != "" {
				codeChallengeMethod = "S256"
			}
			_, err := session.NewAuthorizationCodeSession(
				ids.NewSessionId(), client.ClientID, redirectURI, scope,
				"state", responseType, &mockCodeGenerator{code: "c"},
				tc.codeChallenge, codeChallengeMethod,
				client, newUserIdentityMother(t),
			)
			if tc.wantErr {
				assertErrCode(t, err, tc.wantErrCode)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestSession_Scope(t *testing.T) {
	t.Parallel()
	sess := newSessionMother(t)
	assert.Equal(t, "openid", sess.ToPrimitives().Scope)
}

func TestSession_GrantType(t *testing.T) {
	t.Parallel()
	sess := newSessionMother(t)
	assert.Equal(t, "authorization_code", sess.ToPrimitives().GrantType)
}

func TestBuildRedirectURI(t *testing.T) {
	t.Parallel()

	t.Run("returns redirect uri with code and state params", func(t *testing.T) {
		t.Parallel()
		sess := newSessionMother(t)
		uri := sess.BuildRedirectURI()
		assert.Contains(t, uri, "https://example.com/callback")
		assert.Contains(t, uri, "code=")
	})
}

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	t.Run("moves an authorizing session to active and records an event", func(t *testing.T) {
		t.Parallel()
		sess := newSessionMother(t)
		sess.PullDomainEvents()

		require.NoError(t, sess.Authenticate())
		assert.Equal(t, "active", sess.ToPrimitives().SessionState)
		events := sess.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, protocolEvents.SessionAuthenticatedType, events[0].Type())
	})

	t.Run("fails when the session is not authorizing", func(t *testing.T) {
		t.Parallel()
		sess := newSessionMother(t)
		require.NoError(t, sess.Authenticate())
		sess.PullDomainEvents()

		assertErrCode(t, sess.Authenticate(), session.ErrSessionNotAuthorizing)
		assert.Empty(t, sess.PullDomainEvents())
	})
}
