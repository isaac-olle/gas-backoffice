package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var errInvalidErrorType app_errors.ErrorCode = "INVALID_ERROR_TYPE"

func HttpError(
	presenter app_errors.ErrorPresenter,
	errorManagers ...app_errors.ErrorManager,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		var appError *app_errors.AppError
		if !errors.As(err, &appError) {
			appError = app_errors.InternalError(errInvalidErrorType).
				WithDescription("error type is not 'app_error'").
				WithCause(err)
		}

		for _, manager := range errorManagers {
			manager.Manage(appError)
		}
		presenter.Present(c, appError)
	}
}
