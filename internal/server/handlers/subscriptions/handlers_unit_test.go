package subscriptions_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backofficeSubscriptions "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	subscriptionHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestHandler_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		body       string
		userID     string
		subID      string
		authorOK   bool
		commandErr error
		wantErr    bool
		wantStatus int
	}{
		{name: "invalid create body", method: http.MethodPost, body: `{}`, authorOK: true, wantErr: true},
		{name: "create without author", method: http.MethodPost, body: `{"name":"Team"}`, wantErr: true},
		{name: "create command error", method: http.MethodPost, body: `{"name":"Team"}`, authorOK: true, commandErr: errors.New("command failed"), wantErr: true},
		{name: "create subscription", method: http.MethodPost, body: `{"name":"Team"}`, authorOK: true, wantStatus: http.StatusCreated},
		{name: "invalid subscription id", method: http.MethodPost, userID: ids.NewUserId().String(), subID: "invalid", wantErr: true},
		{name: "invalid user id", method: http.MethodPost, userID: "invalid", subID: ids.NewSubscriptionId().String(), wantErr: true},
		{name: "add member command error", method: http.MethodPost, userID: ids.NewUserId().String(), subID: ids.NewSubscriptionId().String(), commandErr: errors.New("command failed"), wantErr: true},
		{name: "add member", method: http.MethodPost, userID: ids.NewUserId().String(), subID: ids.NewSubscriptionId().String()},
		{name: "remove member", method: http.MethodDelete, userID: ids.NewUserId().String(), subID: ids.NewSubscriptionId().String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			commands := handlerstest.NewCommandExecutor(tt.commandErr)
			author := func(_ context.Context) (*vo.UserAuthor, bool) {
				if !tt.authorOK {
					return nil, false
				}
				value, err := vo.NewUserAuthorFromString(ids.NewUserId().String(), ids.NewUserId().String(), vo.SystemAdministratorRole)
				require.NoError(t, err)
				return value, true
			}
			h := subscriptionHandlers.NewHandler(commands, author)
			_, c := handlerstest.NewCtx(tt.method, strings.NewReader(tt.body))
			c.Params = gin.Params{{Key: "subscription_id", Value: tt.subID}, {Key: "user_id", Value: tt.userID}}

			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			if tt.wantStatus != 0 {
				assert.Equal(t, tt.wantStatus, c.Writer.Status())
			} else {
				assert.Equal(t, http.StatusOK, c.Writer.Status())
			}
		})
	}
}

func TestListHandler_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		authorOK bool
		queryErr error
		wantErr  bool
	}{
		{name: "missing author", wantErr: true},
		{name: "query error", authorOK: true, queryErr: errors.New("query failed"), wantErr: true},
		{name: "success", authorOK: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			queries := handlerstest.NewQueryBus([]backofficeSubscriptions.Summary{}, tt.queryErr)
			author := func(_ context.Context) (*vo.UserAuthor, bool) {
				if !tt.authorOK {
					return nil, false
				}
				value, err := vo.NewUserAuthorFromString(ids.NewUserId().String(), ids.NewUserId().String(), vo.SystemAdministratorRole)
				require.NoError(t, err)
				return value, true
			}
			h := subscriptionHandlers.NewListHandler(queries, author)
			_, c := handlerstest.NewCtx(http.MethodGet, nil)
			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusOK, c.Writer.Status())
		})
	}
}

func TestUpdateNameHandler_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		subID    string
		body     string
		authorOK bool
		busErr   error
		wantErr  bool
	}{
		{name: "invalid subscription id", subID: "invalid", body: `{"name":"New"}`, authorOK: true, wantErr: true},
		{name: "invalid body", subID: ids.NewSubscriptionId().String(), body: `{}`, authorOK: true, wantErr: true},
		{name: "missing author", subID: ids.NewSubscriptionId().String(), body: `{"name":"New"}`, wantErr: true},
		{name: "command error", subID: ids.NewSubscriptionId().String(), body: `{"name":"New"}`, authorOK: true, busErr: errors.New("command failed"), wantErr: true},
		{name: "success", subID: ids.NewSubscriptionId().String(), body: `{"name":"New"}`, authorOK: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			commands := handlerstest.NewCommandExecutor(tt.busErr)
			author := func(_ context.Context) (*vo.UserAuthor, bool) {
				if !tt.authorOK {
					return nil, false
				}
				value, err := vo.NewUserAuthorFromString(ids.NewUserId().String(), ids.NewUserId().String(), vo.SystemAdministratorRole)
				require.NoError(t, err)
				return value, true
			}
			h := subscriptionHandlers.NewUpdateNameHandler(commands, author)
			_, c := handlerstest.NewCtx(http.MethodPatch, strings.NewReader(tt.body))
			c.Params = gin.Params{{Key: "subscription_id", Value: tt.subID}}
			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusNoContent, c.Writer.Status())
		})
	}
}
