package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	userGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

// CurrentUserResponse is the profile of the authenticated user, shaped for the
// frontend session/profile views.
type CurrentUserResponse struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	AuthorId    string `json:"author_id,omitempty"`
}

type GetCurrentUserHandler struct {
	queryBus query.Bus
	author   func(context.Context) (*authors.UserAuthor, bool)
}

func NewGetCurrentUserHandler(
	queryBus query.Bus,
	author func(context.Context) (*authors.UserAuthor, bool),
) *GetCurrentUserHandler {
	return &GetCurrentUserHandler{queryBus: queryBus, author: author}
}

// Execute resolves every authenticated identity through the single User
// aggregate. System administrators differ by system role, not by aggregate.
func (h *GetCurrentUserHandler) Execute(c *gin.Context) (any, error) {
	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}

	user, err := query.DispatchTyped[*users.User](c.Request.Context(), h.queryBus, userGet.NewQuery(author.Id, author))
	if err != nil {
		return nil, err
	}
	p := user.ToPrimitives()
	return CurrentUserResponse{
		Id:          p.Id,
		Email:       p.Email,
		DisplayName: p.Name,
		Role:        author.Role.String(),
	}, nil
}

func (h *GetCurrentUserHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetCurrentUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
