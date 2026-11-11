package identity_providers_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/identity_providers"
)

func TestCreateIdentityProviderHandler_Execute(t *testing.T) {
	t.Parallel()
	validBody := `{"reference_name":"google","name":"Google","issuer":"https://accounts.google.com","client_id":"client","client_secret":"secret","authorization_endpoint":"https://accounts.google.com/auth","token_endpoint":"https://accounts.google.com/token","jwks_uri":"https://accounts.google.com/keys","scopes":["openid","email"]}`
	tests := []struct {
		name    string
		body    string
		busErr  error
		wantErr bool
	}{
		{name: "invalid json", body: `{`, wantErr: true},
		{name: "missing required fields", body: `{"reference_name":"google"}`, wantErr: true},
		{name: "command error", body: validBody, busErr: errors.New("command failed"), wantErr: true},
		{name: "success", body: validBody},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			commands := handlerstest.NewCommandExecutor(tt.busErr)
			h := identity_providers.NewCreateIdentityProviderHandler(commands)
			_, c := handlerstest.NewCtx(http.MethodPost, strings.NewReader(tt.body))
			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusCreated, c.Writer.Status())
		})
	}
}
