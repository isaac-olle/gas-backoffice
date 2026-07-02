package server

import "github.com/gin-gonic/gin"

type RouteProvider interface {
	Register(g *gin.RouterGroup)
}
