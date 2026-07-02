package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	persistance "github.com/isaac-olle/gas/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
)

type TxStarter interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

var errBeginTx app_errors.ErrorCode = "TRANSACTION_BEGIN_FAILED"
var errCommitTx app_errors.ErrorCode = "TRANSACTION_COMMIT_FAILED"

// Transactor starts a Postgres transaction and puts it in the context.
// On success it commits. On any c.Error it rolls back.
// Place it before the handlers that need transactional consistency.
func Transactor(pool TxStarter) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx, err := pool.Begin(c.Request.Context())
		if err != nil {
			_ = c.Error(app_errors.InternalError(errBeginTx).WithCause(err))
			c.Abort()
			return
		}

		ctx := persistance.ContextWithTx(c.Request.Context(), tx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if len(c.Errors) > 0 {
			_ = tx.Rollback(c.Request.Context())
			return
		}

		if err = tx.Commit(c.Request.Context()); err != nil {
			_ = tx.Rollback(c.Request.Context())
			_ = c.Error(app_errors.InternalError(errCommitTx).WithCause(err))
		}
	}
}
