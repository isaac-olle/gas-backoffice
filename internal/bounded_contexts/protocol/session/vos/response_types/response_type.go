package response_types

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"strings"
)

var ErrInvalidResponseType app_errors.ErrorCode = "INVALID_RESPONSE_TYPE"

type CodeGenerator interface {
	Generate() (string, error)
}

type ResponseTypeID string

func (r ResponseTypeID) String() string { return string(r) }

const (
	ResponseTypeCode    ResponseTypeID = "code"
	ResponseTypeToken   ResponseTypeID = "token"
	ResponseTypeIDToken ResponseTypeID = "id_token"
)

type ResponseType interface {
	Type() ResponseTypeID
	String() string
	AddToResponse(response *AuthorizeResponse)
}

func New(responseTypeString string, codeGenerator CodeGenerator) (ResponseType, error) {
	types := strings.Fields(strings.TrimSpace(responseTypeString))
	if len(types) == 0 {
		return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("empty response type")
	}
	for _, part := range types {
		switch part {
		case "code", "token", "id_token":
		default:
			return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("unknown response type: " + part)
		}
	}
	if len(types) == 1 {
		switch types[0] {
		case "code":
			return NewCode(types[0], codeGenerator)
		case "token":
			return NewToken(types[0])
		case "id_token":
			return NewIDToken(types[0])
		}
	}
	raw := strings.Join(types, " ")
	return NewStack(raw, codeGenerator)
}

func FromPrimitives(types []string, code string) (ResponseType, error) {
	if len(types) == 0 {
		return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("empty response types")
	}
	for _, part := range types {
		switch part {
		case "code", "token", "id_token":
		default:
			return nil, app_errors.BadRequestError(ErrInvalidResponseType).WithDescription("unknown response type: " + part)
		}
	}
	if len(types) == 1 {
		switch types[0] {
		case "code":
			return CodeFromPrimitives(code)
		case "token":
			return NewToken(types[0])
		case "id_token":
			return NewIDToken(types[0])
		}
	}
	raw := strings.Join(types, " ")
	return StackFromPrimitives(raw, code)
}
