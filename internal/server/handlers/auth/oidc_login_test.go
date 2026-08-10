package auth_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loginParams(idpID, sessionID string) gin.Params {
	return gin.Params{
		{Key: "idp_id", Value: idpID},
		{Key: "session_id", Value: sessionID},
	}
}

func TestLoginHandler(t *testing.T) {
	t.Parallel()

	t.Run("invalid session id returns error", func(t *testing.T) {
		t.Parallel()
		h := auth.NewLoginHandler(handlerstest.NewCommandExecutor(nil), handlerstest.NewQueryBus(nil, nil))
		_, c := handlerstest.NewCtx(http.MethodGet, nil)
		c.Params = loginParams("idp-123", "not-a-uuid")
		_, err := h.Execute(c)
		assert.Error(t, err)
	})

	t.Run("create session command error returns error", func(t *testing.T) {
		t.Parallel()
		h := auth.NewLoginHandler(handlerstest.NewCommandExecutor(errors.New("cmd error")), handlerstest.NewQueryBus(nil, nil))
		_, c := handlerstest.NewCtx(http.MethodGet, nil)
		c.Params = loginParams("idp-123", ids2.NewSessionId().String())
		_, err := h.Execute(c)
		assert.Error(t, err)
	})

	t.Run("authenticate query error returns error", func(t *testing.T) {
		t.Parallel()
		h := auth.NewLoginHandler(handlerstest.NewCommandExecutor(nil), handlerstest.NewQueryBus(nil, errors.New("query error")))
		_, c := handlerstest.NewCtx(http.MethodGet, nil)
		c.Params = loginParams("idp-123", ids2.NewSessionId().String())
		_, err := h.Execute(c)
		assert.Error(t, err)
	})

	t.Run("success redirects to the authorization url", func(t *testing.T) {
		t.Parallel()
		authURL := "https://accounts.google.com/o/oauth2/v2/auth?state=x"
		h := auth.NewLoginHandler(handlerstest.NewCommandExecutor(nil), handlerstest.NewQueryBus(authURL, nil))

		rec, c := handlerstest.NewCtx(http.MethodGet, nil)
		c.Params = loginParams("idp-123", ids2.NewSessionId().String())

		result, err := h.Execute(c)
		require.NoError(t, err)
		require.NoError(t, h.Success(c, result))

		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, authURL, rec.Header().Get("Location"))
	})
}

func TestLoginHandler_Default(t *testing.T) {
	t.Parallel()
	h := auth.NewLoginHandler(handlerstest.NewCommandExecutor(nil), handlerstest.NewQueryBus(nil, nil))
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
