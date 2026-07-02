package app_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"runtime"

	"github.com/google/uuid"
)

type ErrorCode string

type ErrorTrace struct {
	Origin string
	Line   int         `json:"line"`
	Caller *ErrorTrace `json:"caller,omitempty"`
}

type AppError struct {
	ErrorId     string
	Type        ErrorType
	Code        ErrorCode
	Description string
	Details     map[string]any
	Cause       error
	StackTrace  *ErrorTrace
}

func newAppError(errType ErrorType, code ErrorCode) *AppError {
	e := &AppError{
		ErrorId: uuid.New().String(),
		Type:    errType,
		Code:    code,
	}

	pc, file, line, ok := runtime.Caller(2)
	origin := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		origin = fmt.Sprintf("%s:%s", path.Base(file), fn.Name())
	}

	e.StackTrace = &ErrorTrace{
		Origin: origin,
		Line:   line,
	}

	return e
}

func (e *AppError) Error() string {
	jsonError, _ := json.Marshal(e)
	return string(jsonError)
}

func BadRequestError(code ErrorCode) *AppError {
	return newAppError(BadRequestErrorType, code)
}

func InternalError(code ErrorCode) *AppError {
	return newAppError(InternalErrorType, code)
}

func UnauthorizedError(code ErrorCode) *AppError {
	return newAppError(UnauthorizedErrorType, code)
}

func ForbiddenError(code ErrorCode) *AppError {
	return newAppError(ForbiddenErrorType, code)
}

func NotFoundError(code ErrorCode) *AppError {
	return newAppError(NotFoundErrorType, code)
}

func ConflictError(code ErrorCode) *AppError {
	return newAppError(ConflictErrorType, code)
}

func (e *AppError) WithDescription(desc string) *AppError {
	e.Description = desc
	return e
}

func (e *AppError) WithDetails(details map[string]any) *AppError {
	e.Details = details
	return e
}

func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

func (e *AppError) Is(target ErrorCode) bool {
	return e.Code == target
}

func AddTrace(err error) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	ok := errors.As(err, &appErr)
	if !ok {
		return err
	}

	pc, file, line, ok := runtime.Caller(2)
	origin := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		origin = fmt.Sprintf("%s:%s", path.Base(file), fn.Name())
	}

	appErr.StackTrace = &ErrorTrace{
		Origin: origin,
		Line:   line,
		Caller: appErr.StackTrace,
	}

	return appErr
}

func (t *ErrorTrace) ToJSON() any {
	if t == nil {
		return nil
	}
	return struct {
		Origin string `json:"origin"`
		Line   int    `json:"line"`
		Caller any    `json:"trace,omitempty"`
	}{
		Origin: t.Origin,
		Line:   t.Line,
		Caller: t.Caller.ToJSON(),
	}
}
