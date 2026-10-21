package error_presenters

import (
	app_errors2 "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpErrorPresenter struct{}

func NewHttpErrorPresenter() *HttpErrorPresenter {
	return &HttpErrorPresenter{}
}

func (p *HttpErrorPresenter) Present(c *gin.Context, err *app_errors2.AppError) {
	status := toHttpStatus(err.Type)
	body := gin.H{
		"error_id":    err.ErrorId,
		"code":        string(err.Code),
		"description": err.Description,
		"request_id":  err.RequestId,
		"occurred_at": err.OccurredAt.Format(time.RFC3339),
	}
	if causes := err.PublicCauses(); len(causes) > 0 {
		body["causes"] = causes
	}
	c.JSON(status, body)
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
