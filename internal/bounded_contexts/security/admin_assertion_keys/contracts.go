package admin_assertion_keys

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type Repository interface {
	persistance.Repository[*AdminAssertionKey, ids.AdminAssertionKeyId]
}

type AssertionVerifier interface {
	Verify(ctx context.Context, assertion string, key *AdminAssertionKey) (map[string]any, error)
}

type JwtDecoder interface {
	ExtractKeyId(ctx context.Context, assertion string) (string, error)
}
