package create

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type createIdpRepository struct {
	saved *identity_providers.IdentityProvider
	err   error
}

func (r *createIdpRepository) Save(_ context.Context, idp *identity_providers.IdentityProvider) error {
	r.saved = idp
	return r.err
}
func (r *createIdpRepository) FindById(context.Context, ids.IdentityProviderId) (*identity_providers.IdentityProvider, error) {
	return nil, nil
}
func (r *createIdpRepository) FindByReferenceName(context.Context, string) (*identity_providers.IdentityProvider, error) {
	return nil, nil
}

type idpCreateBus struct{ events int }

func (b *idpCreateBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

type idpTestEncrypter struct {
	err   error
	plain string
}

func (e *idpTestEncrypter) Encrypt(s string) (string, error) {
	e.plain = s
	if e.err != nil {
		return "", e.err
	}
	return "encrypted:" + s, nil
}
func (*idpTestEncrypter) Decrypt(s string) (string, error) { return s, nil }
func (*idpTestEncrypter) Seal(b []byte) ([]byte, error)    { return b, nil }
func (*idpTestEncrypter) Open(b []byte) ([]byte, error)    { return b, nil }

func validCreateIdpCommand() *Command {
	return NewCommand(ids.NewIdentityProviderId(), "google", "Google", "https://issuer.example", "client-id", "secret", "https://issuer.example/authorize", "https://issuer.example/token", "", "https://issuer.example/jwks", []string{"openid"}, "https://app.example/callback")
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name                string
		invalidConfig       bool
		encryptErr, saveErr error
		wantErr             bool
		wantEvents          int
	}{
		{name: "encrypts and persists provider credentials", wantEvents: 1},
		{name: "wraps encryption errors", encryptErr: errors.New("encryption failed"), wantErr: true},
		{name: "rejects invalid provider configuration", invalidConfig: true, wantErr: true},
		{name: "wraps repository errors", saveErr: errors.New("write failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus, encryptor := &createIdpRepository{err: tt.saveErr}, &idpCreateBus{}, &idpTestEncrypter{err: tt.encryptErr}
			cmd := validCreateIdpCommand()
			if tt.invalidConfig {
				cmd = NewCommand(ids.NewIdentityProviderId(), "google", "", "https://issuer.example", "client-id", "secret", "https://issuer.example/authorize", "https://issuer.example/token", "", "https://issuer.example/jwks", []string{"openid"}, "https://app.example/callback")
			}
			err := NewService(repo, bus, encryptor).Create(context.Background(), cmd)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.wantEvents, bus.events)
			if !tt.wantErr {
				require.NotNil(t, repo.saved)
				require.Equal(t, "encrypted:secret", repo.saved.ToPrimitives().ClientSecret)
				require.Equal(t, "secret", encryptor.plain)
			}
		})
	}
}
