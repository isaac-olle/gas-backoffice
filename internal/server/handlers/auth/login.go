package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/initiate_login"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type LoginHandler struct {
	queryBus query.Bus
}

func NewLoginHandler(queryBus query.Bus) *LoginHandler {
	return &LoginHandler{queryBus: queryBus}
}

func (h *LoginHandler) Execute(c *gin.Context) (any, error) {
	qry := initiate_login.NewQuery(c.Query("idp_reference_name"))
	return h.queryBus.Dispatch(c.Request.Context(), qry)
}

func (h *LoginHandler) Success(c *gin.Context, result any) error {
	res := result.(initiate_login.Result)
	c.Redirect(http.StatusFound, res.AuthorizationURL)
	return nil
}

func (h *LoginHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
