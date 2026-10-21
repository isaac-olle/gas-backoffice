package protocol

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JwksExecutor is the narrow view of get_jwks.Service the HTTP layer needs,
// so tests can stub it without wiring the whole dependency graph.
type JwksExecutor interface {
	Get(ctx context.Context) ([]map[string]any, error)
}

type jwksResponse struct {
	Keys []map[string]any `json:"keys"`
}

type JwksHandler struct {
	jwksService JwksExecutor
}

func NewJwksHandler(jwksService JwksExecutor) *JwksHandler {
	return &JwksHandler{jwksService: jwksService}
}

func (h *JwksHandler) Execute(c *gin.Context) (any, error) {
	keys, err := h.jwksService.Get(c.Request.Context())
	if err != nil {
		return nil, err
	}
	return jwksResponse{Keys: keys}, nil
}

func (h *JwksHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *JwksHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
