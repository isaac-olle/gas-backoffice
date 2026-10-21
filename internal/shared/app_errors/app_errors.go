package app_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"runtime"
	"time"

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
	Cause       []CauseError
	StackTrace  *ErrorTrace
	RequestId   string
	OccurredAt  time.Time
}

func newAppError(errType ErrorType, code ErrorCode) *AppError {
	e := &AppError{
		ErrorId:    uuid.New().String(),
		Type:       errType,
		Code:       code,
		OccurredAt: time.Now().UTC(),
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

type AppErrorPrimitives struct {
	ErrorId     string
	Type        string
	Code        string
	Description string
	Details     map[string]any
	Cause       []string
	StackTrace  any
	RequestId   string
	OccurredAt  string
}

func (e *AppError) ToPrimitives() AppErrorPrimitives {
	causes := make([]string, 0, len(e.Cause))
	for _, c := range e.Cause {
		causes = append(causes, c.Error())
	}

	var stackTrace any
	if e.StackTrace != nil {
		stackTrace = e.StackTrace.ToJSON()
	}

	return AppErrorPrimitives{
		ErrorId:     e.ErrorId,
		Type:        e.Type.String(),
		Code:        string(e.Code),
		Description: e.Description,
		Details:     e.Details,
		Cause:       causes,
		StackTrace:  stackTrace,
		RequestId:   e.RequestId,
		OccurredAt:  e.OccurredAt.Format(time.RFC3339),
	}
}

var ErrInvalidRequestBody ErrorCode = "INVALID_REQUEST_BODY"
var ErrQueryingDatabase ErrorCode = "ERROR_QUERYING_DATABASE"

func BadRequestError(code ErrorCode) *AppError {
	return newAppError(BadRequestErrorType, code)
}

func InvalidRequestBodyError(cause error) *AppError {
	return newAppError(BadRequestErrorType, ErrInvalidRequestBody).WithCause(cause, false)
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

func (e *AppError) WithRequestId(requestId string) *AppError {
	e.RequestId = requestId
	return e
}

// WithCause attaches an underlying cause with a visibility flag. If public, the
// HTTP presenter may expose it to the client; otherwise it is only logged. A
// *ConcatenatedError is flattened into its individual errors (each wrapped with
// the same visibility) so every validation error surfaces as its own entry.
func (e *AppError) WithCause(cause error, public bool) *AppError {
	if ce, ok := cause.(*ConcatenatedError); ok {
		for _, sub := range ce.Errors() {
			e.Cause = append(e.Cause, CauseError{public: public, error: sub})
		}
		return e
	}
	e.Cause = append(e.Cause, CauseError{public: public, error: cause})
	return e
}

// PublicCauses returns the messages of the causes explicitly marked public
// (wrapped in a public *CauseError). Plain errors and private causes are omitted,
// so internal details are never leaked to the API client.
func (e *AppError) PublicCauses() []string {
	var out []string
	for _, c := range e.Cause {
		if msg := c.StringIfPublic(); msg != "" {
			out = append(out, msg)
		}
	}
	return out
}

func (e *AppError) Is(target error) bool {
	var t *AppError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
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
