package response_types_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	responsetypes "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos/response_types"
)

// --- mock code generator ---

type mockCodeGen struct{ code string }

func (m *mockCodeGen) Generate() (string, error) { return m.code, nil }

// --- tests ---

func TestNew_CodeResponse(t *testing.T) {
	t.Parallel()
	gen := &mockCodeGen{code: "auth-code-123"}
	rt, err := responsetypes.New("code", gen)
	require.NoError(t, err)
	assert.Equal(t, responsetypes.ResponseTypeCode, rt.Type())
}

func TestNew_TokenResponse(t *testing.T) {
	t.Parallel()
	_, err := responsetypes.New("token", nil)
	assert.Error(t, err, "token response_type is not yet supported via New()")
}

func TestNew_IDTokenResponse(t *testing.T) {
	t.Parallel()
	_, err := responsetypes.New("id_token", nil)
	assert.Error(t, err, "id_token response_type is not yet supported via New()")
}

func TestNew_EmptyResponseType(t *testing.T) {
	t.Parallel()
	_, err := responsetypes.New("", nil)
	assert.Error(t, err)
}

func TestNew_UnknownResponseType(t *testing.T) {
	t.Parallel()
	_, err := responsetypes.New("magic_token", nil)
	assert.Error(t, err)
}

func TestNew_StackResponseType(t *testing.T) {
	t.Parallel()
	gen := &mockCodeGen{code: "code-456"}
	_, err := responsetypes.New("code token", gen)
	assert.Error(t, err, "stack response_type is not yet supported via New()")
}

func TestFromPrimitives_Code(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.FromPrimitives([]string{"code"}, "stored-code")
	require.NoError(t, err)
	code, ok := rt.(*responsetypes.CodeResponse)
	require.True(t, ok)
	assert.Equal(t, "stored-code", code.Code())
}

func TestFromPrimitives_Token(t *testing.T) {
	t.Parallel()
	_, err := responsetypes.FromPrimitives([]string{"token"}, "")
	assert.Error(t, err, "token is not yet supported via FromPrimitives()")
}

func TestFromPrimitives_Empty(t *testing.T) {
	t.Parallel()
	_, err := responsetypes.FromPrimitives([]string{}, "")
	assert.Error(t, err)
}

func TestAuthorizeResponse_Build_QueryParam(t *testing.T) {
	t.Parallel()
	r := responsetypes.NewAuthorizeResponse("https://example.com/callback")
	r.AddQueryParameter("code", "mycode")
	r.AddState("mystate")

	url := r.Build()
	assert.Contains(t, url, "code=mycode")
	assert.Contains(t, url, "state=mystate")
}

func TestAuthorizeResponse_Build_FragmentParam(t *testing.T) {
	t.Parallel()
	r := responsetypes.NewAuthorizeResponse("https://example.com/callback")
	r.AddFragmentParameter("access_token", "tok123")
	r.AddState("s1")

	url := r.Build()
	assert.Contains(t, url, "#")
	assert.Contains(t, url, "access_token=tok123")
	assert.Contains(t, url, "state=s1")
}

func TestStackResponseType_AddToResponse(t *testing.T) {
	t.Parallel()
	gen := &mockCodeGen{code: "c1"}
	rt, err := responsetypes.NewStack("code token", gen)
	require.NoError(t, err)

	resp := responsetypes.NewAuthorizeResponse("https://example.com/cb")
	rt.AddToResponse(resp)
	url := resp.Build()
	assert.Contains(t, url, "code=c1")
}

func TestResponseTypeID_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "code", responsetypes.ResponseTypeCode.String())
	assert.Equal(t, "token", responsetypes.ResponseTypeToken.String())
	assert.Equal(t, "id_token", responsetypes.ResponseTypeIDToken.String())
}

func TestCodeResponse_String(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.New("code", &mockCodeGen{code: "abc"})
	require.NoError(t, err)
	assert.Equal(t, "code", rt.String())
}

func TestCodeResponse_AddToResponse(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.New("code", &mockCodeGen{code: "mycode"})
	require.NoError(t, err)
	resp := responsetypes.NewAuthorizeResponse("https://example.com/cb")
	rt.AddToResponse(resp)
	assert.Contains(t, resp.Build(), "code=mycode")
}

func TestTokenResponse_String(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.NewToken("token")
	require.NoError(t, err)
	assert.Equal(t, "token", rt.String())
}

func TestTokenResponse_AddToResponse(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.NewToken("token")
	require.NoError(t, err)
	resp := responsetypes.NewAuthorizeResponse("https://example.com/cb")
	rt.AddToResponse(resp)
	// access_token has empty value so is omitted; token_type=Bearer is always present
	assert.Contains(t, resp.Build(), "token_type=Bearer")
}

func TestIDTokenResponse_String(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.NewIDToken("id_token")
	require.NoError(t, err)
	assert.Equal(t, "id_token", rt.String())
}

func TestIDTokenResponse_AddToResponse(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.NewIDToken("id_token")
	require.NoError(t, err)
	resp := responsetypes.NewAuthorizeResponse("https://example.com/cb")
	rt.AddToResponse(resp)
	// id_token value is empty so omitted from output, but fragment mode is triggered
	assert.Contains(t, resp.Build(), "#")
}

func TestStackFromPrimitives(t *testing.T) {
	t.Parallel()
	t.Run("rejects stack types from primitives (not yet supported)", func(t *testing.T) {
		_, err := responsetypes.FromPrimitives([]string{"code", "token"}, "stored-code")
		assert.Error(t, err, "stack not yet supported via FromPrimitives()")
	})

	t.Run("rejects unknown type in primitives", func(t *testing.T) {
		_, err := responsetypes.FromPrimitives([]string{"code", "magic"}, "c")
		assert.Error(t, err)
	})

	t.Run("rejects empty primitives list", func(t *testing.T) {
		_, err := responsetypes.FromPrimitives([]string{}, "")
		assert.Error(t, err)
	})
}

func TestStackResponseType_Type(t *testing.T) {
	t.Parallel()
	gen := &mockCodeGen{code: "c"}
	rt, err := responsetypes.NewStack("token code", gen)
	require.NoError(t, err)
	assert.Equal(t, responsetypes.ResponseTypeToken, rt.Type())
}

func TestStackResponseType_String(t *testing.T) {
	t.Parallel()
	gen := &mockCodeGen{code: "c"}
	rt, err := responsetypes.NewStack("code id_token", gen)
	require.NoError(t, err)
	assert.Contains(t, rt.String(), "code")
}

func TestStackResponseType_ResponseTypes_And_Code(t *testing.T) {
	t.Parallel()
	gen := &mockCodeGen{code: "xyz"}
	rt, err := responsetypes.NewStack("code token", gen)
	require.NoError(t, err)
	assert.Len(t, rt.ResponseTypes(), 2)
	assert.Equal(t, "xyz", rt.Code())
}

func TestFromPrimitives_IDToken(t *testing.T) {
	t.Parallel()
	_, err := responsetypes.FromPrimitives([]string{"id_token"}, "")
	assert.Error(t, err, "id_token is not yet supported via FromPrimitives()")
}

// Regressió: els paràmetres es concatenaven crus, així que un valor amb
// caràcters reservats arribava corromput al destí. El cas que ho va destapar
// era un state en base64 estàndard: el '+' es llegia com un espai a l'altra
// banda i el desxifrat petava amb "illegal base64 data".
func TestAuthorizeResponse_Build_EscapesReservedCharacters(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		value string
	}{
		{name: "plus from standard base64", value: "qrvMA+ARIjNE"},
		{name: "slash from standard base64", value: "Pvv/Gis8TV5v"},
		{name: "base64 padding", value: "qrvMA+ARIjNE=="},
		{name: "ampersand", value: "a&b=c"},
		{name: "space", value: "a b"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := responsetypes.NewAuthorizeResponse("https://example.com/callback")
			r.AddQueryParameter("code", "mycode")
			r.AddState(tc.value)

			parsed, err := url.Parse(r.Build())
			require.NoError(t, err)

			assert.Equal(t, tc.value, parsed.Query().Get("state"))
			assert.Equal(t, "mycode", parsed.Query().Get("code"))
		})
	}
}
