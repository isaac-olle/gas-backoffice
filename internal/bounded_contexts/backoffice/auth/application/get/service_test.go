package get_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/domain_services"
	authmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns session when found", func(t *testing.T) {
		t.Parallel()
		repo := authmocks.NewMockRepository(t)
		sessionId := ids.NewSessionId()
		session := aSessionWithId(t, sessionId)

		repo.EXPECT().SearchById(mock.Anything, sessionId).Return(session, nil)

		svc := get.NewService(domain_services.NewSessionFinder(repo))
		result, err := svc.Get(context.Background(), get.NewQuery(sessionId))
		require.NoError(t, err)

		got, ok := result.(*auth.Session)
		require.True(t, ok)
		assert.Equal(t, sessionId.String(), got.ToPrimitives().Id)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := authmocks.NewMockRepository(t)
		sessionId := ids.NewSessionId()

		repo.EXPECT().SearchById(mock.Anything, sessionId).Return(nil, fmt.Errorf("db error"))

		svc := get.NewService(domain_services.NewSessionFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(sessionId))
		assert.Error(t, err)
	})

	t.Run("returns error when session not found", func(t *testing.T) {
		t.Parallel()
		repo := authmocks.NewMockRepository(t)
		sessionId := ids.NewSessionId()

		repo.EXPECT().SearchById(mock.Anything, sessionId).Return(nil, nil)

		svc := get.NewService(domain_services.NewSessionFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(sessionId))
		assert.Error(t, err)
	})
}
