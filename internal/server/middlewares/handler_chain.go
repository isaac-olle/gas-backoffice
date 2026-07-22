package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	persistance "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
)

type TxStarter interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

var errBeginTx app_errors.ErrorCode = "TRANSACTION_BEGIN_FAILED"
var errCommitTx app_errors.ErrorCode = "TRANSACTION_COMMIT_FAILED"

type HandlerChain struct {
	pool      TxStarter
	handlers  []Handler
	onSuccess func(*gin.Context) error
}

func NewHandlerChain(pool TxStarter, onSuccess func(*gin.Context) error, handlers ...Handler) *HandlerChain {
	return &HandlerChain{pool: pool, handlers: handlers, onSuccess: onSuccess}
}

func (ch *HandlerChain) Execute(c *gin.Context) (any, error) {
	tx, err := ch.pool.Begin(c.Request.Context())
	if err != nil {
		return nil, app_errors.InternalError(errBeginTx).WithCause(err)
	}

	ctx := persistance.ContextWithTx(c.Request.Context(), tx)
	c.Request = c.Request.WithContext(ctx)

	for _, h := range ch.handlers {
		if _, err := h.Execute(c); err != nil {
			_ = tx.Rollback(c.Request.Context())
			return nil, err
		}
	}

	if err = tx.Commit(c.Request.Context()); err != nil {
		_ = tx.Rollback(c.Request.Context())
		return nil, app_errors.InternalError(errCommitTx).WithCause(err)
	}

	return nil, nil
}

func (ch *HandlerChain) Success(c *gin.Context, result any) error {
	return ch.onSuccess(c)
}

func (ch *HandlerChain) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
