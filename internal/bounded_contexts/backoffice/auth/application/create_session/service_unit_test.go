package create_session_test

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create_session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/domain/user_projection"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type createSessionDecoder struct {
	verifier string
	err      error
}

func (d *createSessionDecoder) DecodeState(string) (string, error) { return d.verifier, d.err }

type createSessionTokenClient struct {
	result         auth.TokenExchangeResult
	err            error
	code, verifier string
}

func (c *createSessionTokenClient) ExchangeCode(_ context.Context, code, verifier string) (auth.TokenExchangeResult, error) {
	c.code, c.verifier = code, verifier
	return c.result, c.err
}

type createSessionProjectionRepo struct {
	projection *user_projection.UserProjection
	err        error
	email      string
}

func (r *createSessionProjectionRepo) FindByEmail(_ context.Context, email string) (*user_projection.UserProjection, error) {
	r.email = email
	return r.projection, r.err
}

type createSessionTokenGenerator struct {
	token string
	err   error
}

func (g *createSessionTokenGenerator) Generate(context.Context, *auth.Session) (string, error) {
	return g.token, g.err
}

type createSessionRefreshGenerator struct {
	token string
	err   error
}

func (g *createSessionRefreshGenerator) GenerateRefreshToken(context.Context, *auth.Session) (string, error) {
	return g.token, g.err
}

type createSessionRepository struct {
	saveErr error
	saved   int
}

func (r *createSessionRepository) SaveSession(context.Context, *auth.Session) error {
	r.saved++
	return r.saveErr
}
func (*createSessionRepository) IsActive(context.Context, ids.SessionId) (bool, error) {
	return true, nil
}
func (*createSessionRepository) Logout(context.Context, ids.SessionId) error { return nil }

type createSessionBus struct{ events int }

func (b *createSessionBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

func TestExchange(t *testing.T) {
	validProjection := &user_projection.UserProjection{UserId: ids.NewUserId().String(), Email: "ada@example.com", SystemRole: "system_administrator", AdminId: ids.NewUserId().String()}
	tests := []struct {
		name                                                                   string
		decoderErr, exchangeErr, projectionErr, saveErr, accessErr, refreshErr error
		projection                                                             *user_projection.UserProjection
		wantErr                                                                bool
	}{
		{name: "creates a session and returns both tokens", projection: validProjection},
		{name: "rejects undecodable state", decoderErr: errors.New("invalid state"), projection: validProjection, wantErr: true},
		{name: "wraps token exchange errors", exchangeErr: errors.New("exchange failed"), projection: validProjection, wantErr: true},
		{name: "rejects users without a projection", projectionErr: errors.New("not a member"), projection: validProjection, wantErr: true},
		{name: "rejects invalid projection identifiers and roles", projection: &user_projection.UserProjection{UserId: "bad-id", SystemRole: "bad-role", AdminId: "bad-id"}, wantErr: true},
		{name: "wraps session persistence errors", projection: validProjection, saveErr: errors.New("save failed"), wantErr: true},
		{name: "wraps access token generation errors", projection: validProjection, accessErr: errors.New("sign failed"), wantErr: true},
		{name: "wraps refresh token generation errors", projection: validProjection, refreshErr: errors.New("sign failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := &createSessionDecoder{verifier: "pkce-verifier", err: tt.decoderErr}
			client := &createSessionTokenClient{result: auth.TokenExchangeResult{Email: "ada@example.com", Picture: "picture-url"}, err: tt.exchangeErr}
			projection := &createSessionProjectionRepo{projection: tt.projection, err: tt.projectionErr}
			access := &createSessionTokenGenerator{token: "access-token", err: tt.accessErr}
			refresh := &createSessionRefreshGenerator{token: "refresh-token", err: tt.refreshErr}
			sessions, bus := &createSessionRepository{saveErr: tt.saveErr}, &createSessionBus{}
			service := create_session.NewService(client, decoder, projection, access, refresh, sessions, bus)
			result, err := service.Exchange(context.Background(), create_session.NewQuery("auth-code", "encoded-state"))
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				require.Zero(t, bus.events)
				return
			}
			got, ok := result.(create_session.Result)
			require.True(t, ok)
			require.Equal(t, create_session.Result{AccessToken: "access-token", RefreshToken: "refresh-token"}, got)
			require.Equal(t, "pkce-verifier", client.verifier)
			require.Equal(t, "ada@example.com", projection.email)
			require.Equal(t, 1, sessions.saved)
			require.Equal(t, 1, bus.events)
		})
	}
}
