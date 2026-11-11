package app_errors

import (
	"context"
	"github.com/gin-gonic/gin"
)

type ErrorManager interface {
	Manage(ctx context.Context, err *AppError)
}

type ErrorPresenter interface {
	Present(c *gin.Context, err *AppError)
}
