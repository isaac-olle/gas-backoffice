package client_assertion_keys

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type Repository interface {
	persistance.Repository[*ClientAssertionKey, ids.ClientAssertionKeyId]
}

type AssertionVerifier interface {
	Verify(ctx context.Context, assertion string) (map[string]any, error)
}
