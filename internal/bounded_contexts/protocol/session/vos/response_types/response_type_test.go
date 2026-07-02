package response_types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	responsetypes "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/vos/response_types"
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
	rt, err := responsetypes.New("token", nil)
	require.NoError(t, err)
	assert.Equal(t, responsetypes.ResponseTypeToken, rt.Type())
}

func TestNew_IDTokenResponse(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.New("id_token", nil)
	require.NoError(t, err)
	assert.Equal(t, responsetypes.ResponseTypeIDToken, rt.Type())
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
	rt, err := responsetypes.New("code token", gen)
	require.NoError(t, err)
	assert.NotNil(t, rt)
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
	rt, err := responsetypes.FromPrimitives([]string{"token"}, "")
	require.NoError(t, err)
	assert.Equal(t, responsetypes.ResponseTypeToken, rt.Type())
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
	rt, err := responsetypes.New("token", nil)
	require.NoError(t, err)
	assert.Equal(t, "token", rt.String())
}

func TestTokenResponse_AddToResponse(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.New("token", nil)
	require.NoError(t, err)
	resp := responsetypes.NewAuthorizeResponse("https://example.com/cb")
	rt.AddToResponse(resp)
	// access_token has empty value so is omitted; token_type=Bearer is always present
	assert.Contains(t, resp.Build(), "token_type=Bearer")
}

func TestIDTokenResponse_String(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.New("id_token", nil)
	require.NoError(t, err)
	assert.Equal(t, "id_token", rt.String())
}

func TestIDTokenResponse_AddToResponse(t *testing.T) {
	t.Parallel()
	rt, err := responsetypes.New("id_token", nil)
	require.NoError(t, err)
	resp := responsetypes.NewAuthorizeResponse("https://example.com/cb")
	rt.AddToResponse(resp)
	// id_token value is empty so omitted from output, but fragment mode is triggered
	assert.Contains(t, resp.Build(), "#")
}

func TestStackFromPrimitives(t *testing.T) {
	t.Parallel()
	t.Run("reconstructs stack from stored types and code", func(t *testing.T) {
		rt, err := responsetypes.FromPrimitives([]string{"code", "token"}, "stored-code")
		require.NoError(t, err)
		stack, ok := rt.(*responsetypes.StackResponseType)
		require.True(t, ok)
		assert.Equal(t, "stored-code", stack.Code())
		assert.Len(t, stack.ResponseTypes(), 2)
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
	rt, err := responsetypes.FromPrimitives([]string{"id_token"}, "")
	require.NoError(t, err)
	assert.Equal(t, responsetypes.ResponseTypeIDToken, rt.Type())
}
