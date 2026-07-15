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

	t.Run("all_handlers_succeed_commits_and_calls_onSuccess", func(t *testing.T) {
		t.Parallel()
		pool := newFakeTxPool(nil, nil)
		h1 := &stubHandler{}
		h2 := &stubHandler{}

		chain := middlewares.NewHandlerChain(pool, respond201, h1, h2)
		m := middlewares.NewHandlerMiddleware()

		r := gin.New()
		r.GET("/test", m.Run(chain))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.True(t, pool.committed_())
		assert.False(t, pool.rolledBack_())
	})

	t.Run("first_handler_fails_rolls_back_and_aborts", func(t *testing.T) {
		t.Parallel()
		pool := newFakeTxPool(nil, nil)
		h1 := &stubHandler{executeErr: errors.New("h1 fail")}
		h2 := &stubHandler{}

		chain := middlewares.NewHandlerChain(pool, respond201, h1, h2)
		m := middlewares.NewHandlerMiddleware()

		r := gin.New()
		r.GET("/test", m.Run(chain))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		assert.False(t, h2.executeCalled, "second handler should not run")
		assert.True(t, pool.rolledBack_())
		assert.False(t, pool.committed_())
	})

	t.Run("second_handler_fails_rolls_back", func(t *testing.T) {
		t.Parallel()
		pool := newFakeTxPool(nil, nil)
		h1 := &stubHandler{}
		h2 := &stubHandler{executeErr: errors.New("h2 fail")}

		chain := middlewares.NewHandlerChain(pool, respond201, h1, h2)
		m := middlewares.NewHandlerMiddleware()

		r := gin.New()
		r.GET("/test", m.Run(chain))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		assert.True(t, h1.executeCalled)
		assert.True(t, pool.rolledBack_())
		assert.False(t, pool.committed_())
	})

	t.Run("begin_tx_fails_aborts", func(t *testing.T) {
		t.Parallel()
		pool := newFakeTxPool(errors.New("db down"), nil)
		h1 := &stubHandler{}

		chain := middlewares.NewHandlerChain(pool, respond201, h1)
		m := middlewares.NewHandlerMiddleware()

		r := gin.New()
		r.GET("/test", m.Run(chain))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		assert.False(t, h1.executeCalled)
	})

	t.Run("commit_fails_rolls_back", func(t *testing.T) {
		t.Parallel()
		pool := newFakeTxPool(nil, errors.New("commit fail"))
		h1 := &stubHandler{}

		chain := middlewares.NewHandlerChain(pool, respond201, h1)
		m := middlewares.NewHandlerMiddleware()

		r := gin.New()
		r.GET("/test", m.Run(chain))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		assert.True(t, pool.rolledBack_())
		assert.NotEqual(t, http.StatusCreated, w.Code)
	})
}
