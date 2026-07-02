package infrastructure

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/alexedwards/argon2id"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets"
)

type Argon2SecretVerifier struct {
	repository client_secrets.Repository
}

func NewArgon2SecretVerifier(repository client_secrets.Repository) *Argon2SecretVerifier {
	return &Argon2SecretVerifier{repository: repository}
}

func (v *Argon2SecretVerifier) Verify(ctx context.Context, clientID, plainSecret string) error {
	cid, err := ids.NewClientIdFromString(clientID)
	if err != nil {
		return app_errors.UnauthorizedError(client_secrets.ErrInvalidClientSecret).
			WithDescription("invalid client_id")
	}

	secret, err := v.repository.FindByClientId(ctx, cid)
	if err != nil {
		return err
	}

	p := secret.ToPrimitives()
	match, err := argon2id.ComparePasswordAndHash(plainSecret, p.SecretHash)
	if err != nil {
		return app_errors.InternalError(client_secrets.ErrInvalidClientSecret).
			WithDescription("failed to verify client secret").
			WithCause(err)
	}
	if !match {
		return app_errors.UnauthorizedError(client_secrets.ErrInvalidClientSecret).
			WithDescription("invalid client_secret")
	}
	return nil
}
