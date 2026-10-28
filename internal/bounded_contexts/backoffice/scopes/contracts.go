package scopes

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type Repository interface {
	persistance.Repository[*Scope, ids.ScopeId]
	// FindByAllowedClientId returns every non-deleted scope whose allowedClientIds
	// contains the given client. Used to clean up a client's grants when the
	// client itself is deleted.
	FindByAllowedClientId(ctx context.Context, clientId ids.ClientId) ([]*Scope, error)
}

type ClientProjectionRepository interface {
	FindClientProjection(context.Context, ids.ClientId) (*ClientProjection, error)
}
