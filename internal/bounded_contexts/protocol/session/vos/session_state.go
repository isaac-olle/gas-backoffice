package vos

import "errors"

type SessionState int

const (
	SessionStateActive SessionState = iota
	Authorizing
	SessionStateConsumed
	SessionStateExpired
)

func NewSessionState(s string) (SessionState, error) {
	switch s {
	case "active":
		return SessionStateActive, nil
	case "authorizing":
		return Authorizing, nil
	case "consumed":
		return SessionStateConsumed, nil
	case "expired":
		return SessionStateExpired, nil
	default:
		return 0, errors.New("invalid session state")
	}
}

func (s SessionState) String() string {
	switch s {
	case SessionStateActive:
		return "active"
	case Authorizing:
		return "authorizing"
	case SessionStateConsumed:
		return "consumed"
	case SessionStateExpired:
		return "expired"
	default:
		return "unknown"
	}
}
