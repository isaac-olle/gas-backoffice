package connections_provider

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionsProvider struct {
	postgresPool *pgxpool.Pool
}

func NewConnectionsProvider(cfg *config.Config) (*ConnectionsProvider, error) {
	pool, err := pgxpool.New(context.Background(), cfg.Databases.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	return &ConnectionsProvider{postgresPool: pool}, nil
}

func (p *ConnectionsProvider) PostgresPool() *pgxpool.Pool {
	return p.postgresPool
}

func (p *ConnectionsProvider) Close() {
	if p.postgresPool != nil {
		p.postgresPool.Close()
	}
}
