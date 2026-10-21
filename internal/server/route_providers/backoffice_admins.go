package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	userHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// registerAdmins registers admin-only CRUD routes for managing users.
// The target user identity comes from the :user_id path param.
func registerAdmins(g *gin.RouterGroup, p *Backoffice) {
	a := g.Group("/users/:user_id")
	a.Use(
		middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins),
		middlewares.RoleMiddleware(vo.AdminRole),
		path_params.UserIdParam(),
	)

	a.PUT("",
		p.executor.Run(userHandlers.NewCreateUserHandler(
			p.commandExecutor,
			ids.UserIdFromContext,
			helpers.UserIdFromAuthor,
		)),
	)
	a.GET("",
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, ids.UserIdFromContext, helpers.AuthorFromContext)),
	)
	a.DELETE("",
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandExecutor, ids.UserIdFromContext, helpers.AuthorFromContext)),
	)
	a.PATCH("/name",
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandExecutor, ids.UserIdFromContext, helpers.AuthorFromContext)),
	)
}
