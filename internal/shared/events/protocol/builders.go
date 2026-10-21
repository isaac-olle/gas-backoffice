package protocol

import (
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

// ── Session ───────────────────────────────────────────────────────────────

type sessionCreatedBuilder struct{}

func NewSessionCreatedBuilder() event2.PrimitiveEventBuilder {
	return &sessionCreatedBuilder{}
}
func (b *sessionCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return SessionCreated{
		ClientID:    str(p.Data, "client_id"),
		RedirectURI: str(p.Data, "redirect_uri"),
		Scope:       str(p.Data, "scope"),
		Flow:        str(p.Data, "flow"),
		BaseEvent:   event2.FromPrimitivesBase(p),
	}, nil
}

type sessionAuthenticatedBuilder struct{}

func NewSessionAuthenticatedBuilder() event2.PrimitiveEventBuilder {
	return &sessionAuthenticatedBuilder{}
}
func (b *sessionAuthenticatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return SessionAuthenticated{
		ClientID:  str(p.Data, "client_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

func str(data map[string]any, key string) string {
	v, _ := data[key].(string)
	return v
}
