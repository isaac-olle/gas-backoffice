package session_test

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
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
		wantFlow      string
	}{
		{
			name:         "creates code flow session",
			clientID:     uuid.New().String(),
			redirectURI:  "https://example.com/cb",
			scope:        "openid profile",
			responseType: "code",
			wantFlow:     "authorization_code",
		},
		{
			name:          "creates authorization_code flow session when challenge is provided",
			clientID:      uuid.New().String(),
			redirectURI:   "https://example.com/cb",
			scope:         "openid",
			responseType:  "code",
			pkceChallenge: "challenge",
			pkceMethod:    "S256",
			wantFlow:      "authorization_code",
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
			sess, err := session.NewAuthorizationCodeSession(
				ids.NewSessionId(), tc.clientID, tc.redirectURI, tc.scope,
				"state", tc.responseType, &mockCodeGenerator{code: "c"},
				tc.pkceChallenge, tc.pkceMethod,
			)
			if tc.wantErr {
				assertErrCode(t, err, tc.wantErrCode)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, sess)
			if tc.wantFlow != "" {
				assert.Equal(t, tc.wantFlow, sess.ToPrimitives().Flow)
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
		sess, err := session.NewAuthorizationCodeSession(
			ids.NewSessionId(), uuid.New().String(),
			"https://example.com/cb", "openid", "s",
			"code", &mockCodeGenerator{code: "c"}, "challenge", "plain",
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

func TestCheckAuthorizationRequest(t *testing.T) {
	t.Parallel()

	const (
		responseType = "code"
		redirectURI  = "https://example.com/callback"
		scope        = "openid"
	)

	tests := []struct {
		name        string
		flows       []string
		uris        []string
		scopes      []string
		wantErr     bool
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:   "accepts client with matching flow, uri, and scope",
			flows:  []string{"authorization_code"},
			uris:   []string{redirectURI},
			scopes: []string{scope},
		},
		{
			name:        "rejects client without the authorization_code flow",
			flows:       []string{"client_credentials"},
			uris:        []string{redirectURI},
			scopes:      []string{scope},
			wantErr:     true,
			wantErrCode: session.ErrFlowNotAllowed,
		},
		{
			name:        "rejects client with wrong redirect uri",
			flows:       []string{"authorization_code"},
			uris:        []string{"https://other.com/cb"},
			scopes:      []string{scope},
			wantErr:     true,
			wantErrCode: session.ErrInvalidRedirectURI,
		},
		{
			name:        "rejects client without required scope",
			flows:       []string{"authorization_code"},
			uris:        []string{redirectURI},
			scopes:      []string{"profile"},
			wantErr:     true,
			wantErrCode: session.ErrInvalidScope,
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
				tc.flows, tc.uris, scopeObjs,
			)

			_, err := client.CheckAuthorizationRequest(responseType, redirectURI, scope)
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

func TestSession_Flow(t *testing.T) {
	t.Parallel()
	sess := newSessionMother(t)
	assert.Equal(t, "authorization_code", sess.ToPrimitives().Flow)
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
