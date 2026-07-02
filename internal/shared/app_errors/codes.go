package app_errors

type ErrorType string

func (e ErrorType) String() string {
	return string(e)
}

const (
	BadRequestErrorType   ErrorType = "BAD_REQUEST"
	UnauthorizedErrorType ErrorType = "UNAUTHORIZED"
	NotFoundErrorType     ErrorType = "NOT_FOUND"
	ConflictErrorType     ErrorType = "CONFLICT"
	InternalErrorType     ErrorType = "INTERNAL"
	ForbiddenErrorType    ErrorType = "FORBIDDEN"
)
