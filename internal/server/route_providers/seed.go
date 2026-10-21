package route_providers

import (
	"github.com/gin-gonic/gin"
	seedHandler "github.com/isaac-olle/gas-backoffice/internal/server/handlers/seed"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	sharedSeed "github.com/isaac-olle/gas-backoffice/internal/shared/seed"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Seed struct {
	pool      *pgxpool.Pool
	executor  command.CommandExecutor
	seedToken string
}

func NewSeed(pool *pgxpool.Pool, executor command.CommandExecutor, seedToken string) *Seed {
	return &Seed{pool: pool, executor: executor, seedToken: seedToken}
}

func (p *Seed) Register(root *gin.RouterGroup) {
	g := root.Group("/v1/seed")
	g.Use(middlewares.SeedToken(p.seedToken))
	h := seedHandler.NewHandler(sharedSeed.NewRepository(p.pool), p.executor)
	g.POST("", h.Execute)
}
