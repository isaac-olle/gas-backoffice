package user_assertion_keys

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type Repository interface {
	persistance.Repository[*UserAssertionKey, ids.UserAssertionKeyId]
}

type AssertionVerifier interface {
	Verify(ctx context.Context, assertion string, key *UserAssertionKey) (map[string]any, error)
}

type JwtDecoder interface {
	ExtractKeyId(ctx context.Context, assertion string) (string, error)
}
