package error_presenters

import (
	app_errors2 "github.com/isaac-olle/gas/internal/shared/app_errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpErrorPresenter struct{}

func NewHttpErrorPresenter() *HttpErrorPresenter {
	return &HttpErrorPresenter{}
}

func (p *HttpErrorPresenter) Present(c *gin.Context, err *app_errors2.AppError) {
	status := toHttpStatus(err.Type)
	c.JSON(status, gin.H{
		"error_id":    err.ErrorId,
		"code":        string(err.Code),
		"description": err.Description,
	})
}

func toHttpStatus(errType app_errors2.ErrorType) int {
	switch errType {
	case app_errors2.BadRequestErrorType:
		return http.StatusBadRequest
	case app_errors2.UnauthorizedErrorType:
		return http.StatusUnauthorized
	case app_errors2.ForbiddenErrorType:
		return http.StatusForbidden
	case app_errors2.NotFoundErrorType:
		return http.StatusNotFound
	case app_errors2.ConflictErrorType:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
