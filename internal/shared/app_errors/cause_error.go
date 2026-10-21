package app_errors

// CauseError wraps an underlying error together with a visibility flag. A public
// cause may be exposed to the API client (the HTTP presenter includes it in the
// response); a private cause is only ever logged internally, never returned.
//
// The error interface is embedded, so CauseError.Error() delegates to the
// wrapped error; Unwrap keeps errors.Is/As working through the chain.
type CauseError struct {
	public bool
	error
}

func (c CauseError) IsPublic() bool {
	return c.public
}

// StringIfPublic returns the cause message when public, or an empty string when
// private — so callers can map over causes without leaking internal details.
func (c CauseError) StringIfPublic() string {
	if c.public {
		return c.Error()
	}
	return ""
}

func (c CauseError) Unwrap() error {
	return c.error
}
