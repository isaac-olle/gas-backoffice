package middlewares_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
)

func TestHandlerChain(t *testing.T) {
	t.Parallel()

	respond201 := func(c *gin.Context) error {
		c.JSON(http.StatusCreated, gin.H{})
		return nil
	}

	tests := []struct {
		name      string
		beginErr  error
		commitErr error
		h1Err     error
		h2Err     error
		skipH2    bool
		assert    func(t *testing.T, w *httptest.ResponseRecorder, pool *fakeTxPool, h1, h2 *stubHandler)
	}{
		{
			name: "all_handlers_succeed_commits_and_calls_onSuccess",
			assert: func(t *testing.T, w *httptest.ResponseRecorder, pool *fakeTxPool, h1, h2 *stubHandler) {
				assert.Equal(t, http.StatusCreated, w.Code)
				assert.True(t, pool.committed_())
				assert.False(t, pool.rolledBack_())
			},
		},
		{
			name:  "first_handler_fails_rolls_back_and_aborts",
			h1Err: errors.New("h1 fail"),
			assert: func(t *testing.T, w *httptest.ResponseRecorder, pool *fakeTxPool, h1, h2 *stubHandler) {
				assert.False(t, h2.executeCalled, "second handler should not run")
				assert.True(t, pool.rolledBack_())
				assert.False(t, pool.committed_())
			},
		},
		{
			name:  "second_handler_fails_rolls_back",
			h2Err: errors.New("h2 fail"),
			assert: func(t *testing.T, w *httptest.ResponseRecorder, pool *fakeTxPool, h1, h2 *stubHandler) {
				assert.True(t, h1.executeCalled)
				assert.True(t, pool.rolledBack_())
				assert.False(t, pool.committed_())
			},
		},
		{
			name:     "begin_tx_fails_aborts",
			beginErr: errors.New("db down"),
			skipH2:   true,
			assert: func(t *testing.T, w *httptest.ResponseRecorder, pool *fakeTxPool, h1, h2 *stubHandler) {
				assert.False(t, h1.executeCalled)
			},
		},
		{
			name:      "commit_fails_rolls_back",
			commitErr: errors.New("commit fail"),
			skipH2:    true,
			assert: func(t *testing.T, w *httptest.ResponseRecorder, pool *fakeTxPool, h1, h2 *stubHandler) {
				assert.True(t, pool.rolledBack_())
				assert.NotEqual(t, http.StatusCreated, w.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pool := newFakeTxPool(tt.beginErr, tt.commitErr)
			h1 := &stubHandler{executeErr: tt.h1Err}
			h2 := &stubHandler{executeErr: tt.h2Err}

			var chain *middlewares.HandlerChain
			if tt.skipH2 {
				chain = middlewares.NewHandlerChain(pool, respond201, h1)
			} else {
				chain = middlewares.NewHandlerChain(pool, respond201, h1, h2)
			}
			m := middlewares.NewHandlerMiddleware()

			r := gin.New()
			r.GET("/test", m.Run(chain))

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			tt.assert(t, w, pool, h1, h2)
		})
	}
}
