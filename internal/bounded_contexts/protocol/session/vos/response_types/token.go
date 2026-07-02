package response_types

import "fmt"

type TokenResponse struct{}

func NewToken(raw string) (*TokenResponse, error) {
	if raw != string(ResponseTypeToken) {
		return nil, fmt.Errorf("invalid response_type: %s", raw)
	}
	return &TokenResponse{}, nil
}

func (t *TokenResponse) Type() ResponseTypeID { return ResponseTypeToken }
func (t *TokenResponse) String() string       { return string(ResponseTypeToken) }
func (t *TokenResponse) AddToResponse(response *AuthorizeResponse) {
	response.AddFragmentParameter("access_token", "")
	response.AddFragmentParameter("token_type", "Bearer")
}
