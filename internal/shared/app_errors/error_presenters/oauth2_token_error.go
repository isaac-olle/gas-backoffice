package error_presenters

import (
	app_errors2 "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuth2TokenErrorPresenter struct{}

func NewOAuth2TokenErrorPresenter() *OAuth2TokenErrorPresenter {
	return &OAuth2TokenErrorPresenter{}
}

func (p *OAuth2TokenErrorPresenter) Present(c *gin.Context, err *app_errors2.AppError) {
	oauthCode, status := toOAuth2TokenError(err)
	c.JSON(status, gin.H{
		"error":             oauthCode,
		"error_description": err.Description,
	})
}

func toOAuth2TokenError(err *app_errors2.AppError) (string, int) {
	switch err.Type {
	case app_errors2.UnauthorizedErrorType:
		return "invalid_client", http.StatusUnauthorized
	case app_errors2.InternalErrorType:
		return "server_error", http.StatusInternalServerError
	default:
		switch err.Code {
		case "UNSUPPORTED_GRANT_TYPE":
			return "unsupported_grant_type", http.StatusBadRequest
		case "INVALID_SCOPE":
			return "invalid_scope", http.StatusBadRequest
		case "CODE_NOT_FOUND", "SESSION_ALREADY_CONSUMED", "INVALID_CODE_VERIFIER", "INVALID_CLIENT_SECRET":
			return "invalid_grant", http.StatusBadRequest
		default:
			return "invalid_request", http.StatusBadRequest
		}
	}
}
