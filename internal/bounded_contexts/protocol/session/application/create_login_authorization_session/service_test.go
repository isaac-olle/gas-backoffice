package create_login_authorization_session_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	create "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/create_login_authorization_session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const (
	injectedRedirectURI = "https://example.com/callback"
	injectedScope       = "openid"
)

func newService(sessionRepo *mockSessionRepo, backofficeClientId string, bus *mockEventBus) *create.Service {
	return create.NewService(sessionRepo, bus, &mockCodeGen{code: "authcode"}, backofficeClientId, injectedRedirectURI, injectedScope, uuid.New().String())
}

func TestCreateLoginAuthorizationSession_HappyPath(t *testing.T) {
	t.Parallel()

	backofficeClientId := uuid.New().String()
	sessionID := ids.NewSessionId()
	sessionRepo := &mockSessionRepo{}
	bus := &mockEventBus{}

	svc := newService(sessionRepo, backofficeClientId, bus)

	require.NoError(t, svc.Create(context.Background(), create.NewCommand(sessionID)))

	require.Len(t, sessionRepo.saved, 1)
	p := sessionRepo.saved[0].ToPrimitives()
	assert.Equal(t, sessionID.String(), p.SessionID)
	assert.Equal(t, backofficeClientId, p.ClientID, "client id must be the injected backoffice client")
	assert.Equal(t, injectedScope, p.Scope, "scope must be injected")
	assert.Equal(t, injectedRedirectURI, p.RedirectURI, "redirect uri must be injected")
	assert.NotEmpty(t, p.State, "state must be computed")
	assert.True(t, bus.publishCalled)
}

func TestCreateLoginAuthorizationSession_ComputesDistinctState(t *testing.T) {
	t.Parallel()

	backofficeClientId := uuid.New().String()
	sessionRepo := &mockSessionRepo{}
	svc := newService(sessionRepo, backofficeClientId, &mockEventBus{})

	require.NoError(t, svc.Create(context.Background(), create.NewCommand(ids.NewSessionId())))
	require.NoError(t, svc.Create(context.Background(), create.NewCommand(ids.NewSessionId())))

	require.Len(t, sessionRepo.saved, 2)
	assert.NotEqual(t, sessionRepo.saved[0].ToPrimitives().State, sessionRepo.saved[1].ToPrimitives().State)
}

func TestCreateLoginAuthorizationSession_WrapsSaveError(t *testing.T) {
	t.Parallel()

	sessionRepo := &mockSessionRepo{err: errors.New("db error")}
	bus := &mockEventBus{}

	svc := newService(sessionRepo, uuid.New().String(), bus)

	err := svc.Create(context.Background(), create.NewCommand(ids.NewSessionId()))

	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, session.ErrCannotSaveSession, appErr.Code)
	assert.False(t, bus.publishCalled)
}
