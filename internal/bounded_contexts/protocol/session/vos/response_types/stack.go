package response_types

import (
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"strings"
)

type StackResponseType struct {
	responseTypes []ResponseTypeID
	raw           string
	code          string
}

func NewStack(raw string, codeGenerator CodeGenerator) (*StackResponseType, error) {
	parts := strings.Fields(strings.TrimSpace(raw))
	if len(parts) == 0 {
		return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("empty stack response type")
	}
	stack := &StackResponseType{responseTypes: make([]ResponseTypeID, 0, len(parts)), raw: raw}
	seen := make(map[ResponseTypeID]bool)
	hasCode := false
	for _, part := range parts {
		var rt ResponseTypeID
		switch part {
		case "code":
			rt = ResponseTypeCode
			hasCode = true
		case "token":
			rt = ResponseTypeToken
		case "id_token":
			rt = ResponseTypeIDToken
		default:
			return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("unknown response type: " + part)
		}
		if seen[rt] {
			return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("duplicate response type: " + part)
		}
		seen[rt] = true
		stack.responseTypes = append(stack.responseTypes, rt)
	}
	if hasCode {
		code, err := codeGenerator.Generate()
		if err != nil {
			return nil, fmt.Errorf("failed to generate code: %w", err)
		}
		stack.code = code
	}
	return stack, nil
}

func StackFromPrimitives(raw string, codeValue string) (*StackResponseType, error) {
	parts := strings.Fields(strings.TrimSpace(raw))
	if len(parts) == 0 {
		return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("empty stack response type primitives")
	}
	stack := &StackResponseType{responseTypes: make([]ResponseTypeID, 0, len(parts)), raw: raw, code: codeValue}
	seen := make(map[ResponseTypeID]bool)
	for _, part := range parts {
		var rt ResponseTypeID
		switch part {
		case "code":
			rt = ResponseTypeCode
		case "token":
			rt = ResponseTypeToken
		case "id_token":
			rt = ResponseTypeIDToken
		default:
			return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("unknown response type: " + part)
		}
		if seen[rt] {
			return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("duplicate response type: " + part)
		}
		seen[rt] = true
		stack.responseTypes = append(stack.responseTypes, rt)
	}
	return stack, nil
}

func (s *StackResponseType) Type() ResponseTypeID {
	if len(s.responseTypes) > 0 {
		return s.responseTypes[0]
	}
	return ""
}

func (s *StackResponseType) String() string {
	if s.raw != "" {
		return s.raw
	}
	types := make([]string, len(s.responseTypes))
	for i, rt := range s.responseTypes {
		types[i] = string(rt)
	}
	return strings.Join(types, " ")
}

func (s *StackResponseType) ResponseTypes() []ResponseTypeID { return s.responseTypes }
func (s *StackResponseType) Code() string                    { return s.code }

func (s *StackResponseType) AddToResponse(response *AuthorizeResponse) {
	for _, rt := range s.responseTypes {
		switch rt {
		case ResponseTypeCode:
			response.AddFragmentParameter("code", s.code)
		case ResponseTypeToken:
			response.AddFragmentParameter("access_token", "")
			response.AddFragmentParameter("token_type", "Bearer")
		case ResponseTypeIDToken:
			response.AddFragmentParameter("id_token", "")
		}
	}
}
