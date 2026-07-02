package exchange_code_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/exchange_code"
	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestExchangeCodeService(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		sess := newPKCESessionMother(t, clientID)
		prim := sess.ToPrimitives()
		client := newPKCEClientMother(clientID, prim.RedirectURI)

		sessionRepo := &mockSessionRepo{byCode: sess}
		clientRepo := &mockClientRepo{client: client}
		tokenGen := &mockTokenGen{token: "signed-token"}
		sessionFinder := domain_services.NewSessionFinder(sessionRepo)
		clientFinder := domain_services.NewClientFinder(clientRepo)
		svc := exchange_code.NewService(sessionFinder, clientFinder, sessionRepo, tokenGen)

		qry := exchange_code.NewQuery(prim.Code, clientID, prim.RedirectURI, "my-challenge")
		result, err := svc.ExchangeCode(context.Background(), qry)
		require.NoError(t, err)

		tokenResp, ok := result.(*session.TokenResponse)
		require.True(t, ok)
		assert.Equal(t, "signed-token", tokenResp.Token)
	})

	t.Run("returns code not found when session does not exist", func(t *testing.T) {
		t.Parallel()
		sessionRepo := &mockSessionRepo{byCode: nil}
		clientRepo := &mockClientRepo{}
		tokenGen := &mockTokenGen{token: "tok"}
		sessionFinder := domain_services.NewSessionFinder(sessionRepo)
		clientFinder := domain_services.NewClientFinder(clientRepo)
		svc := exchange_code.NewService(sessionFinder, clientFinder, sessionRepo, tokenGen)

		qry := exchange_code.NewQuery("missing-code", uuid.New().String(), "https://example.com/cb", "")
		_, err := svc.ExchangeCode(context.Background(), qry)
		assertErrCode(t, err, session.ErrCodeNotFound)
	})

	t.Run("returns code not found for wrong client id", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		sess := newPKCESessionMother(t, clientID)

		sessionRepo := &mockSessionRepo{byCode: sess}
		clientRepo := &mockClientRepo{}
		tokenGen := &mockTokenGen{}
		sessionFinder := domain_services.NewSessionFinder(sessionRepo)
		clientFinder := domain_services.NewClientFinder(clientRepo)
		svc := exchange_code.NewService(sessionFinder, clientFinder, sessionRepo, tokenGen)

		qry := exchange_code.NewQuery("auth-code", "wrong-client-id", "https://example.com/cb", "my-challenge")
		_, err := svc.ExchangeCode(context.Background(), qry)
		assertErrCode(t, err, session.ErrCodeNotFound)
	})

	t.Run("returns session consumed error for already consumed session", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		sess := newPKCESessionMother(t, clientID)
		prim := sess.ToPrimitives()
		sess.Consume()

		sessionRepo := &mockSessionRepo{byCode: sess}
		clientRepo := &mockClientRepo{}
		tokenGen := &mockTokenGen{}
		sessionFinder := domain_services.NewSessionFinder(sessionRepo)
		clientFinder := domain_services.NewClientFinder(clientRepo)
		svc := exchange_code.NewService(sessionFinder, clientFinder, sessionRepo, tokenGen)

		qry := exchange_code.NewQuery(prim.Code, clientID, prim.RedirectURI, "my-challenge")
		_, err := svc.ExchangeCode(context.Background(), qry)
		assertErrCode(t, err, session.ErrSessionConsumed)
	})

	t.Run("returns error when PKCE verifier is empty", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		sess := newPKCESessionMother(t, clientID)
		prim := sess.ToPrimitives()

		sessionRepo := &mockSessionRepo{byCode: sess}
		clientRepo := &mockClientRepo{}
		tokenGen := &mockTokenGen{}
		sessionFinder := domain_services.NewSessionFinder(sessionRepo)
		clientFinder := domain_services.NewClientFinder(clientRepo)
		svc := exchange_code.NewService(sessionFinder, clientFinder, sessionRepo, tokenGen)

		qry := exchange_code.NewQuery(prim.Code, clientID, prim.RedirectURI, "")
		_, err := svc.ExchangeCode(context.Background(), qry)
		assertErrCode(t, err, session.ErrInvalidCodeVerifier)
	})

	t.Run("returns error when client finder fails", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		sess := newPKCESessionMother(t, clientID)
		prim := sess.ToPrimitives()

		sessionRepo := &mockSessionRepo{byCode: sess}
		clientRepo := &mockClientRepo{err: errors.New("db error")}
		tokenGen := &mockTokenGen{}
		sessionFinder := domain_services.NewSessionFinder(sessionRepo)
		clientFinder := domain_services.NewClientFinder(clientRepo)
		svc := exchange_code.NewService(sessionFinder, clientFinder, sessionRepo, tokenGen)

		qry := exchange_code.NewQuery(prim.Code, clientID, prim.RedirectURI, "my-challenge")
		_, err := svc.ExchangeCode(context.Background(), qry)
		assert.Error(t, err)
	})

	t.Run("returns internal error when session save fails after consume", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		sess := newPKCESessionMother(t, clientID)
		prim := sess.ToPrimitives()
		client := newPKCEClientMother(clientID, prim.RedirectURI)

		sessionRepo := &mockSessionRepo{byCode: sess, saveErr: errors.New("save error")}
		clientRepo := &mockClientRepo{client: client}
		tokenGen := &mockTokenGen{}
		sessionFinder := domain_services.NewSessionFinder(sessionRepo)
		clientFinder := domain_services.NewClientFinder(clientRepo)
		svc := exchange_code.NewService(sessionFinder, clientFinder, sessionRepo, tokenGen)

		qry := exchange_code.NewQuery(prim.Code, clientID, prim.RedirectURI, "my-challenge")
		_, err := svc.ExchangeCode(context.Background(), qry)
		assertErrCode(t, err, session.ErrCannotSaveSession)
	})

	t.Run("returns internal error when token generator fails", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		sess := newPKCESessionMother(t, clientID)
		prim := sess.ToPrimitives()
		client := newPKCEClientMother(clientID, prim.RedirectURI)

		sessionRepo := &mockSessionRepo{byCode: sess}
		clientRepo := &mockClientRepo{client: client}
		tokenGen := &mockTokenGen{err: errors.New("gen error")}
		sessionFinder := domain_services.NewSessionFinder(sessionRepo)
		clientFinder := domain_services.NewClientFinder(clientRepo)
		svc := exchange_code.NewService(sessionFinder, clientFinder, sessionRepo, tokenGen)

		qry := exchange_code.NewQuery(prim.Code, clientID, prim.RedirectURI, "my-challenge")
		_, err := svc.ExchangeCode(context.Background(), qry)
		assertErrCode(t, err, session.ErrCannotSignToken)
	})
}
