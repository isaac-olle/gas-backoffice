package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	adminGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	userGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/get"
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
	AdminId     string `json:"admin_id,omitempty"`
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

// Execute looks the caller up in whichever bounded context matches their
// role: an Admin and a Tenant are now separate aggregates (see
// backoffice/admins and backoffice/tenants), so "my profile" means something
// different depending on who is asking.
func (h *GetCurrentUserHandler) Execute(c *gin.Context) (any, error) {
	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}

	if author.Role.IsAdmin() {
		admin, err := query.DispatchTyped[*admins.Admin](c.Request.Context(), h.queryBus, adminGet.NewQuery(author.Id, author.Id))
		if err != nil {
			return nil, err
		}
		p := admin.ToPrimitives()
		return CurrentUserResponse{
			Id:          p.Id,
			Email:       p.Email,
			DisplayName: p.Name,
			Role:        author.Role.String(),
			AdminId:     p.Id,
		}, nil
	}

	tenant, err := query.DispatchTyped[*tenants.Tenant](c.Request.Context(), h.queryBus, userGet.NewQuery(author.Id, author))
	if err != nil {
		return nil, err
	}
	p := tenant.ToPrimitives()
	return CurrentUserResponse{
		Id:          p.Id,
		Email:       p.Email,
		DisplayName: p.Name,
		Role:        author.Role.String(),
		AdminId:     p.AdminId,
	}, nil
}

func (h *GetCurrentUserHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetCurrentUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
