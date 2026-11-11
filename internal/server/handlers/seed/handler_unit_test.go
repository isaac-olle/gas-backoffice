package seed_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type seedRepoStub struct {
	seeded     bool
	checkErr   error
	markErr    error
	checkCalls int
	markCalls  int
}

func (r *seedRepoStub) IsSeeded(context.Context) (bool, error) {
	r.checkCalls++
	return r.seeded, r.checkErr
}
func (r *seedRepoStub) MarkSeeded(context.Context) error {
	r.markCalls++
	return r.markErr
}

type seedExecutorStub struct {
	err      error
	ctx      context.Context
	commands []command.Command
}

func (*seedExecutorStub) SingleExecution(context.Context, command.Command) error { return nil }
func (e *seedExecutorStub) MultipleExecution(ctx context.Context, commands ...command.Command) error {
	e.ctx = ctx
	e.commands = commands
	return e.err
}

const validSeedBody = `{"admin":{"name":"Admin","email":"admin@example.com"},"subscription":{"name":"Main"},"protected_resource":{"id":"d1b2c3d4-e5f6-7890-abcd-ef1234567890","name":"API"},"scopes":[{"id":"e1b2c3d4-e5f6-7890-abcd-ef1234567890","name":"openid"}],"client":{"id":"c1b2c3d4-e5f6-7890-abcd-ef1234567890","name":"Backoffice","redirect_uris":["https://app.example/callback"],"secret":"client-secret"},"identity_providers":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","reference_name":"google","name":"Google","issuer":"https://issuer.example","client_id":"client-id","authorization_endpoint":"https://issuer.example/auth","token_endpoint":"https://issuer.example/token","jwks_uri":"https://issuer.example/jwks"}]}`

func TestHandler_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		body         string
		seeded       bool
		checkErr     error
		markErr      error
		executorErr  error
		wantStatus   int
		wantCommands int
		wantCheck    int
	}{
		{name: "invalid json", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "already seeded", body: validSeedBody, seeded: true, wantStatus: http.StatusConflict, wantCheck: 1},
		{name: "seed state query failed", body: validSeedBody, checkErr: errors.New("query failed"), wantStatus: http.StatusInternalServerError, wantCheck: 1},
		{name: "invalid protected resource id", body: strings.Replace(validSeedBody, "d1b2c3d4-e5f6-7890-abcd-ef1234567890", "invalid", 1), wantStatus: http.StatusBadRequest, wantCheck: 1},
		{name: "invalid client id", body: strings.Replace(validSeedBody, "c1b2c3d4-e5f6-7890-abcd-ef1234567890", "invalid", 1), wantStatus: http.StatusBadRequest, wantCheck: 1},
		{name: "invalid scope id", body: strings.Replace(validSeedBody, "e1b2c3d4-e5f6-7890-abcd-ef1234567890", "invalid", 1), wantStatus: http.StatusBadRequest, wantCheck: 1},
		{name: "invalid identity provider id", body: strings.Replace(validSeedBody, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", "invalid", 1), wantStatus: http.StatusBadRequest, wantCheck: 1},
		{name: "application command failed", body: validSeedBody, executorErr: app_errors.BadRequestError("SEED_REJECTED"), wantStatus: http.StatusUnprocessableEntity, wantCheck: 1, wantCommands: 8},
		{name: "unexpected command failure", body: validSeedBody, executorErr: errors.New("command failed"), wantStatus: http.StatusInternalServerError, wantCheck: 1, wantCommands: 8},
		{name: "marking seed failed", body: validSeedBody, markErr: errors.New("write failed"), wantStatus: http.StatusInternalServerError, wantCheck: 1, wantCommands: 8},
		{name: "success", body: validSeedBody, wantStatus: http.StatusCreated, wantCheck: 1, wantCommands: 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &seedRepoStub{seeded: tt.seeded, checkErr: tt.checkErr, markErr: tt.markErr}
			executor := &seedExecutorStub{err: tt.executorErr}
			h := NewHandler(repo, executor)
			w, c := handlerstest.NewCtx(http.MethodPost, strings.NewReader(tt.body))
			h.Execute(c)

			assert.Equal(t, tt.wantStatus, c.Writer.Status())
			assert.Equal(t, tt.wantCheck, repo.checkCalls)
			assert.Equal(t, tt.wantCommands, len(executor.commands))
			if tt.wantCommands == 0 {
				assert.Zero(t, repo.markCalls)
				return
			}
			if tt.executorErr != nil {
				assert.Zero(t, repo.markCalls)
				return
			}
			assert.Equal(t, 1, repo.markCalls)
			if tt.wantStatus == http.StatusCreated {
				var response struct {
					AdminId string `json:"admin_id"`
				}
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
				assert.NotEmpty(t, response.AdminId)
				author, ok := vo.AuthorFromContext(executor.ctx)
				assert.True(t, ok)
				if ok {
					assert.Equal(t, vo.SystemAdministratorRole, author.Role)
					subscriptionID, found := ids.SubscriptionIdFromContext(executor.ctx)
					assert.True(t, found)
					assert.NotEmpty(t, subscriptionID.String())
				}
			}
		})
	}
}
