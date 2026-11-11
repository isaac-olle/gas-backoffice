package error_presenters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/stretchr/testify/require"
)

func TestOAuth2TokenErrorPresenterRedirectDetails(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	details := map[string]any{
		"expected_redirect_uri": "https://app.example/callback/google",
		"received_redirect_uri": "https://app.example/callback/:idp_reference_name",
	}
	err := app_errors.BadRequestError("INVALID_REDIRECT_URI").
		WithDescription("redirect_uri does not match the one used in the authorization request").
		WithDetails(details)

	NewOAuth2TokenErrorPresenter().Present(ctx, err)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	require.Equal(t, "invalid_request", body["error"])
	require.Equal(t, details, body["details"])
}
