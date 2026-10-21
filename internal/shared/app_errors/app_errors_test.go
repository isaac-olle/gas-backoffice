package app_errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCode    ErrorCode = "TEST_CODE"
	anotherCode ErrorCode = "ANOTHER_CODE"
)

func TestBadRequestError(t *testing.T) {
	t.Parallel()
	err := BadRequestError(testCode)
	require.NotNil(t, err)
	assert.Equal(t, BadRequestErrorType, err.Type)
	assert.Equal(t, testCode, err.Code)
	assert.NotEmpty(t, err.ErrorId)
	assert.NotNil(t, err.StackTrace)
}

func TestInternalError(t *testing.T) {
	t.Parallel()
	err := InternalError(testCode)
	assert.Equal(t, InternalErrorType, err.Type)
	assert.Equal(t, testCode, err.Code)
}

func TestUnauthorizedError(t *testing.T) {
	t.Parallel()
	err := UnauthorizedError(testCode)
	assert.Equal(t, UnauthorizedErrorType, err.Type)
}

func TestForbiddenError(t *testing.T) {
	t.Parallel()
	err := ForbiddenError(testCode)
	assert.Equal(t, ForbiddenErrorType, err.Type)
}

func TestNotFoundError(t *testing.T) {
	t.Parallel()
	err := NotFoundError(testCode)
	assert.Equal(t, NotFoundErrorType, err.Type)
}

func TestConflictError(t *testing.T) {
	t.Parallel()
	err := ConflictError(testCode)
	assert.Equal(t, ConflictErrorType, err.Type)
}

func TestWithDescription(t *testing.T) {
	t.Parallel()
	err := BadRequestError(testCode).WithDescription("something went wrong")
	assert.Equal(t, "something went wrong", err.Description)
}

func TestWithCause(t *testing.T) {
	t.Parallel()
	cause := fmt.Errorf("underlying error")
	err := InternalError(testCode).WithCause(cause, false)

	assert.Len(t, err.Cause, 1)
	assert.Equal(t, "underlying error", err.Cause[0].Error())
	assert.False(t, err.Cause[0].IsPublic())
}

func TestWithCause_AppendsWithVisibility(t *testing.T) {
	t.Parallel()
	err := InternalError(testCode).
		WithCause(fmt.Errorf("a"), true).
		WithCause(fmt.Errorf("b"), false)

	assert.Len(t, err.Cause, 2)
	assert.True(t, err.Cause[0].IsPublic())
	assert.False(t, err.Cause[1].IsPublic())
}

func TestWithCause_FlattensConcatenatedError(t *testing.T) {
	t.Parallel()
	ce := NewEmptyErrorArray()
	ce.AddError(fmt.Errorf("first"))
	ce.AddError(fmt.Errorf("second"))

	err := InternalError(testCode).WithCause(ce, true)

	// The concatenated error is flattened into its individual errors, each
	// carrying the visibility passed to WithCause.
	assert.Len(t, err.Cause, 2)
	assert.Equal(t, []string{"first", "second"}, err.PublicCauses())
}

func TestWithDetails(t *testing.T) {
	t.Parallel()
	details := map[string]any{"field": "value", "count": 5}
	err := BadRequestError(testCode).WithDetails(details)
	assert.Equal(t, details, err.Details)
}

func TestAppError_Is(t *testing.T) {
	t.Parallel()
	err := BadRequestError(testCode)
	assert.True(t, err.Is(BadRequestError(testCode)))
	assert.False(t, err.Is(BadRequestError(anotherCode)))
}

func TestAppError_Error_ReturnsJSON(t *testing.T) {
	t.Parallel()
	err := BadRequestError(testCode).WithDescription("desc")
	jsonStr := err.Error()
	assert.Contains(t, jsonStr, "TEST_CODE")
	assert.Contains(t, jsonStr, "desc")
}

func TestAddTrace_NilError(t *testing.T) {
	t.Parallel()
	result := AddTrace(nil)
	assert.Nil(t, result)
}

func TestAddTrace_NonAppError(t *testing.T) {
	t.Parallel()
	plain := errors.New("plain error")
	result := AddTrace(plain)
	assert.Equal(t, plain, result)
}

func TestAddTrace_AppError(t *testing.T) {
	t.Parallel()
	appErr := BadRequestError(testCode)
	originalTrace := appErr.StackTrace

	// Wrap the call so AddTrace picks up a proper caller frame
	wrapped := func() error { return AddTrace(appErr) }()
	var ae *AppError
	require.True(t, errors.As(wrapped, &ae))
	// After AddTrace the caller is chained
	assert.Equal(t, originalTrace, ae.StackTrace.Caller)
}

func TestErrorType_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "BAD_REQUEST", BadRequestErrorType.String())
	assert.Equal(t, "NOT_FOUND", NotFoundErrorType.String())
	assert.Equal(t, "INTERNAL", InternalErrorType.String())
}

func TestErrorTrace_ToJSON_Nil(t *testing.T) {
	t.Parallel()
	var trace *ErrorTrace
	assert.Nil(t, trace.ToJSON())
}

func TestErrorTrace_ToJSON_NonNil(t *testing.T) {
	t.Parallel()
	trace := &ErrorTrace{Origin: "file.go:fn", Line: 42}
	j := trace.ToJSON()
	assert.NotNil(t, j)
}

func TestAppError_ToPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("flattens type, code, description, details, and stack trace", func(t *testing.T) {
		t.Parallel()
		err := BadRequestError(testCode).
			WithDescription("bad input").
			WithDetails(map[string]any{"field": "name"})

		p := err.ToPrimitives()
		assert.Equal(t, err.ErrorId, p.ErrorId)
		assert.Equal(t, "BAD_REQUEST", p.Type)
		assert.Equal(t, "TEST_CODE", p.Code)
		assert.Equal(t, "bad input", p.Description)
		assert.Equal(t, map[string]any{"field": "name"}, p.Details)
		assert.NotNil(t, p.StackTrace)
	})

	t.Run("flattens every cause into its error message", func(t *testing.T) {
		t.Parallel()
		err := BadRequestError(testCode).
			WithCause(fmt.Errorf("public cause"), true).
			WithCause(fmt.Errorf("private cause"), false)

		p := err.ToPrimitives()
		assert.Equal(t, []string{"public cause", "private cause"}, p.Cause)
	})

	t.Run("returns empty cause slice and nil stack trace when absent", func(t *testing.T) {
		t.Parallel()
		err := &AppError{ErrorId: "id", Type: InternalErrorType, Code: testCode}

		p := err.ToPrimitives()
		assert.Empty(t, p.Cause)
		assert.Nil(t, p.StackTrace)
	})
}
