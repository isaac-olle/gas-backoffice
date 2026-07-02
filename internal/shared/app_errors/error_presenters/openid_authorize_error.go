package error_presenters

import (
	app_errors2 "github.com/isaac-olle/gas/internal/shared/app_errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type OpenIDAuthorizeErrorPresenter struct {
	noRedirectCodes map[app_errors2.ErrorCode]struct{}
}

func NewOpenIDAuthorizeErrorPresenter(noRedirectCodes ...app_errors2.ErrorCode) *OpenIDAuthorizeErrorPresenter {
	m := make(map[app_errors2.ErrorCode]struct{}, len(noRedirectCodes))
	for _, c := range noRedirectCodes {
		m[c] = struct{}{}
	}
	return &OpenIDAuthorizeErrorPresenter{noRedirectCodes: m}
}

func (p *OpenIDAuthorizeErrorPresenter) Present(c *gin.Context, err *app_errors2.AppError) {
	oidcCode := toOIDCErrorCode(err)

	if _, noRedirect := p.noRedirectCodes[err.Code]; noRedirect {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             oidcCode,
			"error_description": err.Description,
		})
		return
	}

	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             oidcCode,
			"error_description": err.Description,
		})
		return
	}

	state := c.Query("state")
	params := url.Values{}
	params.Set("error", oidcCode)
	params.Set("error_description", err.Description)
	if state != "" {
		params.Set("state", state)
	}

	c.Redirect(http.StatusFound, redirectURI+"?"+params.Encode())
}

func toOIDCErrorCode(err *app_errors2.AppError) string {
	switch err.Type {
	case app_errors2.UnauthorizedErrorType:
		return "unauthorized_client"
	case app_errors2.ForbiddenErrorType:
		return "access_denied"
	case app_errors2.InternalErrorType:
		return "server_error"
	default:
		return "invalid_request"
	}
}
