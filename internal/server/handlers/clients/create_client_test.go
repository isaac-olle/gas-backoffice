package clients_test

import (
	"bytes"
	"context"
	"errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestCreateClientHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		body        io.Reader
		clientIdFn  func(context.Context) (ids2.ClientId, bool)
		userIdFn    func(context.Context) (ids2.UserId, bool)
		ownershipFn func(context.Context) (users.Ownership, bool)
		busErr      error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "invalid_body_returns_error",
			body:        bytes.NewBufferString(`{}`),
			clientIdFn:  handlerstest.ValidClientIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_client_id_returns_error",
			body:        handlerstest.JsonBody(map[string]any{"name": "My Client", "allowed_flows": []string{"authorization_code"}, "redirect_uris": []string{"https://example.com/callback"}}),
			clientIdFn:  handlerstest.MissingClientIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_user_id_returns_error",
			body:        handlerstest.JsonBody(map[string]any{"name": "My Client", "allowed_flows": []string{"authorization_code"}, "redirect_uris": []string{"https://example.com/callback"}}),
			clientIdFn:  handlerstest.ValidClientIdExtractor(),
			userIdFn:    handlerstest.MissingUserIdExtractor(),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_ownership_returns_error",
			body:        handlerstest.JsonBody(map[string]any{"name": "My Client", "allowed_flows": []string{"authorization_code"}, "redirect_uris": []string{"https://example.com/callback"}}),
			clientIdFn:  handlerstest.ValidClientIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			ownershipFn: handlerstest.MissingUserOwnershipExtractor(),
			wantErr:     true,
		},
		{
			name:        "bus_error_returns_error",
			body:        handlerstest.JsonBody(map[string]any{"name": "My Client", "allowed_flows": []string{"authorization_code"}, "redirect_uris": []string{"https://example.com/callback"}}),
			clientIdFn:  handlerstest.ValidClientIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			busErr:      errors.New("bus error"),
			wantErr:     true,
		},
		{
			name:        "success_returns_201",
			body:        handlerstest.JsonBody(map[string]any{"name": "My Client", "allowed_flows": []string{"authorization_code"}, "redirect_uris": []string{"https://example.com/callback"}}),
			clientIdFn:  handlerstest.ValidClientIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			wantErr:     false,
			wantStatus:  http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := clients.NewCreateClientHandler(bus, tt.clientIdFn, tt.userIdFn, tt.ownershipFn)

			_, c := handlerstest.NewCtx(http.MethodPost, tt.body)
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, tt.wantStatus, c.Writer.Status())
		})
	}
}

func TestCreateClientHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := clients.NewCreateClientHandler(bus, handlerstest.ValidClientIdExtractor(), handlerstest.ValidUserIdExtractor(), handlerstest.ValidUserOwnershipExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
