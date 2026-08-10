package get_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

func aSessionWithId(t *testing.T, sessionId ids2.SessionId) *auth.Session {
	t.Helper()
	session, err := auth.NewSession(sessionId, ids2.NewUserId())
	require.NoError(t, err)
	session.PullDomainEvents()
	return session
}
