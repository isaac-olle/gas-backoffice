package authentication_methods_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol/authentication_methods"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockMethod struct {
	canProcess bool
	canErr     error
	clientID   string
	authErr    error
}

func (m *mockMethod) CanProcess(_ *gin.Context) (bool, error)     { return m.canProcess, m.canErr }
func (m *mockMethod) Authenticate(_ *gin.Context) (string, error) { return m.clientID, m.authErr }

func newCtx() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
	return c
}

func TestClientAuthMethodFactory_Process(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		methods     []*mockMethod
		wantErr     bool
		wantErrCode string
		wantClient  string
	}{
		{
			name:        "no_methods_registered",
			methods:     nil,
			wantErr:     true,
			wantErrCode: "NO_AUTHENTICATION_METHOD",
		},
		{
			name:        "no_method_matches",
			methods:     []*mockMethod{{canProcess: false}},
			wantErr:     true,
			wantErrCode: "NO_AUTHENTICATION_METHOD",
		},
		{
			name:    "can_process_returns_error",
			methods: []*mockMethod{{canProcess: true, canErr: assert.AnError}},
			wantErr: true,
		},
		{
			name:    "authenticate_returns_error",
			methods: []*mockMethod{{canProcess: true, authErr: assert.AnError}},
			wantErr: true,
		},
		{
			name:       "first_skipped_second_matches",
			methods:    []*mockMethod{{canProcess: false}, {canProcess: true, clientID: "client-b"}},
			wantClient: "client-b",
		},
		{
			name:       "success",
			methods:    []*mockMethod{{canProcess: true, clientID: "client-123"}},
			wantClient: "client-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var iface []authentication_methods.AuthenticationMethod
			for _, m := range tt.methods {
				iface = append(iface, m)
			}
			factory := authentication_methods.NewClientAuthMethodFactory(iface...)
			c := newCtx()

			clientID, err := factory.Process(c)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantClient, clientID)
		})
	}
}
