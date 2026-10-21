package client_projection

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Repository interface {
	Find(ctx context.Context, id ids.ClientId) (*ClientProjection, error)
}
