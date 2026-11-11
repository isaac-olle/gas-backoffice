package error_logs

import "time"

// ErrorLog is a plain record over public.errors: the error manager creates one
// per managed AppError via the create use case, and list reads them back. It
// is not an aggregate — there is nothing to validate and no domain events.
type ErrorLog struct {
	AuthorId    string         `json:"author_id"`
	AuthorRole  string         `json:"author_role"`
	ErrorId     string         `json:"error_id"`
	Type        string         `json:"type"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	Details     map[string]any `json:"details"`
	Cause       []string       `json:"cause"`
	StackTrace  any            `json:"stack_trace"`
	RequestId   string         `json:"request_id"`
	OccurredAt  *time.Time     `json:"occurred_at"`
}
