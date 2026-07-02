package app_errors

import "github.com/gin-gonic/gin"

type ErrorManager interface {
	Manage(err *AppError)
}

type ErrorPresenter interface {
	Present(c *gin.Context, err *AppError)
}
