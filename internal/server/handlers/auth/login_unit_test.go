package auth_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/initiate_login"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestLoginHandler_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		result  any
		busErr  error
		wantErr bool
	}{
		{name: "query error", busErr: errors.New("query failed"), wantErr: true},
		{name: "success", result: initiate_login.Result{AuthorizationURL: "https://idp.example/authorize"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			queries := handlerstest.NewQueryBus(tt.result, tt.busErr)
			h := auth.NewLoginHandler(queries)
			_, c := handlerstest.NewCtx(http.MethodGet, nil)
			c.Request.URL.RawQuery = "idp_reference_name=google"
			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusFound, c.Writer.Status())
			assert.Equal(t, "https://idp.example/authorize", c.Writer.Header().Get("Location"))
		})
	}
}
