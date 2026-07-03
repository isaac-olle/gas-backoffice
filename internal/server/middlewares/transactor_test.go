package middlewares_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
)

type stubTx struct {
	pgx.Tx
	commitErr      error
	rollbackErr    error
	commitCalled   bool
	rollbackCalled bool
}

func (s *stubTx) Commit(_ context.Context) error {
	s.commitCalled = true
	return s.commitErr
}

func (s *stubTx) Rollback(_ context.Context) error {
	s.rollbackCalled = true
	return s.rollbackErr
}

func (s *stubTx) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

type stubTxStarter struct {
	tx  pgx.Tx
	err error
}

func (s *stubTxStarter) Begin(_ context.Context) (pgx.Tx, error) {
	return s.tx, s.err
}

func TestTransactor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		beginErr     error
		handlerErr   error
		commitErr    error
		wantCommit   bool
		wantRollback bool
		wantAbort    bool
		wantCtxHasTx bool
	}{
		{
			name:         "begin_fails_aborts_request",
			beginErr:     errors.New("connection error"),
			wantAbort:    true,
			wantCommit:   false,
			wantRollback: false,
			wantCtxHasTx: false,
		},
		{
			name:         "handler_success_commits",
			wantAbort:    false,
			wantCommit:   true,
			wantRollback: false,
			wantCtxHasTx: true,
		},
		{
			name:         "handler_error_rolls_back",
			handlerErr:   errors.New("handler failed"),
			wantAbort:    false,
			wantCommit:   false,
			wantRollback: true,
			wantCtxHasTx: true,
		},
		{
			name:         "commit_error_rolls_back",
			commitErr:    errors.New("commit failed"),
			wantAbort:    false,
			wantCommit:   true,
			wantRollback: true,
			wantCtxHasTx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &stubTx{commitErr: tt.commitErr}
			var pool middlewares.TxStarter
			if tt.beginErr != nil {
				pool = &stubTxStarter{err: tt.beginErr}
			} else {
				pool = &stubTxStarter{tx: tx}
			}

			nextCalled := false
			r := gin.New()
			r.GET("/test", middlewares.Transactor(pool), func(c *gin.Context) {
				nextCalled = true
				if tt.handlerErr != nil {
					_ = c.Error(tt.handlerErr)
				}
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			if tt.beginErr != nil {
				assert.False(t, nextCalled, "next should not be called when begin fails")
				return
			}

			assert.Equal(t, tt.wantCommit, tx.commitCalled)
			assert.Equal(t, tt.wantRollback, tx.rollbackCalled)
		})
	}
}
