package session

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type Repository interface {
	Save(ctx context.Context, sess *AuthorizationCodeSession) error
	FindById(ctx context.Context, id ids.SessionId) (*AuthorizationCodeSession, error)
	FindByCode(ctx context.Context, code string) (*AuthorizationCodeSession, error)
	SoftDelete(ctx context.Context, id ids.SessionId) error
	HardDelete(ctx context.Context, id ids.SessionId) error
}

type ClientFinderRepository interface {
	FindById(ctx context.Context, clientID string) (*Client, error)
}

type TokenDecoder interface {
	DecodeKid(token string) (string, error)
}

type ClientCredentialsSessionRepository interface {
	Save(ctx context.Context, sess *ClientCredentialsSession) error
}

type TokenGenerator interface {
	GenerateFromSession(ctx context.Context, sess *AuthorizationCodeSession, client *Client) (string, error)
	GenerateFromClientCredentials(ctx context.Context, sess *ClientCredentialsSession) (string, error)
}
