package protocol_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas/internal/server/handlers/protocol"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthorizeHandler(t *testing.T) {
	t.Parallel()

	validQuery := "response_type=code&client_id=abc&redirect_uri=https://example.com/cb&scope=openid&state=xyz"

	tests := []struct {
		name         string
		query        string
		cmdErr       error
		qryResult    any
		qryErr       error
		wantErr      bool
		wantStatus   int
		wantLocation string
	}{
		{
			name:    "missing_required_params",
			query:   "",
			wantErr: true,
		},
		{
			name:    "command_bus_error",
			query:   validQuery,
			cmdErr:  assert.AnError,
			wantErr: true,
		},
		{
			name:    "query_bus_error",
			query:   validQuery,
			qryErr:  assert.AnError,
			wantErr: true,
		},
		{
			name:         "success",
			query:        validQuery,
			qryResult:    "https://example.com/cb?code=abc123&state=xyz",
			wantStatus:   http.StatusFound,
			wantLocation: "https://example.com/cb?code=abc123&state=xyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := protocol.NewAuthorizeHandler(
				handlerstest.NewCommandBus(tt.cmdErr),
				handlerstest.NewQueryBus(tt.qryResult, tt.qryErr),
			)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/?"+tt.query, nil)

			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantLocation, w.Header().Get("Location"))
		})
	}
}

func TestAuthorizeHandler_Default(t *testing.T) {
	t.Parallel()

	h := protocol.NewAuthorizeHandler(
		handlerstest.NewCommandBus(nil),
		handlerstest.NewQueryBus(nil, nil),
	)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	h.Default(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
