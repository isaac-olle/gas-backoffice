package infrastructure_test

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type mockQueryBus struct {
	result any
	err    error
}

func (m *mockQueryBus) Dispatch(_ context.Context, _ query.Query) (any, error) {
	return m.result, m.err
}

func (m *mockQueryBus) Register(_ query.Type, _ query.Handler) {}
