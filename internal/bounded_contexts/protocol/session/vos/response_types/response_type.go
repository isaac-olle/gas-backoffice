package response_types

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
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

// New only supports "code" for now — the authorization_code flow is the only
// one implemented. Reject anything else here instead of spreading the check
// across token/id_token/stack machinery that has no caller yet.
func New(responseTypeString string, codeGenerator CodeGenerator) (ResponseType, error) {
	if responseTypeString != string(ResponseTypeCode) {
		return nil, app_errors.BadRequestError(ErrInvalidResponseType).
			WithDescription("unsupported response_type: " + responseTypeString)
	}
	return NewCode(responseTypeString, codeGenerator)
}

func FromPrimitives(types []string, code string) (ResponseType, error) {
	if len(types) != 1 || types[0] != string(ResponseTypeCode) {
		return nil, app_errors.BadRequestError(ErrInvalidResponseType).
			WithDescription("unsupported response_type")
	}
	return CodeFromPrimitives(code)
}
