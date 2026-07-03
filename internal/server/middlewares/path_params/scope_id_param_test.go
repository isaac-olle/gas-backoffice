package path_params_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestScopeIdParam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		param     string
		wantNext  bool
		wantInCtx bool
	}{
		{name: "valid_uuid_sets_id_in_context", param: "a1b2c3d4-e5f6-7890-abcd-ef1234567890", wantNext: true, wantInCtx: true},
		{name: "invalid_uuid_aborts", param: "not-a-uuid", wantNext: false, wantInCtx: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotId ids.ScopeId
			var gotInCtx bool
			nextCalled := false

			r := gin.New()
			r.GET("/scopes/:scope_id", path_params.ScopeIdParam(), func(c *gin.Context) {
				nextCalled = true
				gotId, gotInCtx = ids.ScopeIdFromContext(c.Request.Context())
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/scopes/"+tt.param, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantNext, nextCalled)
			assert.Equal(t, tt.wantInCtx, gotInCtx)
			if tt.wantInCtx {
				assert.Equal(t, tt.param, gotId.String())
			}
		})
	}
}
