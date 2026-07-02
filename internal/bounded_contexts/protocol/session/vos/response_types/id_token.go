package response_types

import "fmt"

type IDTokenResponse struct{}

func NewIDToken(raw string) (*IDTokenResponse, error) {
	if raw != string(ResponseTypeIDToken) {
		return nil, fmt.Errorf("invalid response_type: %s", raw)
	}
	return &IDTokenResponse{}, nil
}

func (i *IDTokenResponse) Type() ResponseTypeID { return ResponseTypeIDToken }
func (i *IDTokenResponse) String() string       { return string(ResponseTypeIDToken) }
func (i *IDTokenResponse) AddToResponse(response *AuthorizeResponse) {
	response.AddFragmentParameter("id_token", "")
}
