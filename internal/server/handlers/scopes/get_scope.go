package scopes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	scopeGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type GetScopeHandler struct {
	queryBus query.Bus
	scopeId  func(context.Context) (ids.ScopeId, bool)
	author   func(context.Context) (*vo.UserAuthor, bool)
}

func NewGetScopeHandler(
	queryBus query.Bus,
	scopeId func(context.Context) (ids.ScopeId, bool),
	author func(context.Context) (*vo.UserAuthor, bool),
) *GetScopeHandler {
	return &GetScopeHandler{queryBus: queryBus, scopeId: scopeId, author: author}
}

func (h *GetScopeHandler) Execute(c *gin.Context) (any, error) {
	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}
	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	scope, err := query.DispatchTyped[*scopes.Scope](c.Request.Context(), h.queryBus, scopeGet.NewQuery(scopeId, author))
	if err != nil {
		return nil, err
	}
	return scope.ToPrimitives(), nil
}

func (h *GetScopeHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetScopeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
