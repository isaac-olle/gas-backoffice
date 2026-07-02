package response_types

import "fmt"

type CodeResponse struct{ code string }

func NewCode(responseType string, codeGenerator CodeGenerator) (*CodeResponse, error) {
	if responseType != "code" {
		return nil, fmt.Errorf("invalid response_type")
	}
	code, err := codeGenerator.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}
	return &CodeResponse{code: code}, nil
}

func CodeFromPrimitives(codeValue string) (*CodeResponse, error) {
	if codeValue == "" {
		return nil, fmt.Errorf("code cannot be empty")
	}
	return &CodeResponse{code: codeValue}, nil
}

func (c *CodeResponse) Type() ResponseTypeID { return ResponseTypeCode }
func (c *CodeResponse) String() string       { return string(ResponseTypeCode) }
func (c *CodeResponse) Code() string         { return c.code }
func (c *CodeResponse) AddToResponse(response *AuthorizeResponse) {
	response.AddQueryParameter("code", c.code)
}
