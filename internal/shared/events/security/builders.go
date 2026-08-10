package security

import (
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type clientAssertionKeyCreatedBuilder struct{}

func NewClientAssertionKeyCreatedBuilder() event2.PrimitiveEventBuilder {
	return &clientAssertionKeyCreatedBuilder{}
}
func (b *clientAssertionKeyCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientAssertionKeyCreated{
		ClientId:  str(p.Data, "client_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type clientSecretCreatedBuilder struct{}

func NewClientSecretCreatedBuilder() event2.PrimitiveEventBuilder {
	return &clientSecretCreatedBuilder{}
}
func (b *clientSecretCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientSecretCreated{
		ClientId:  str(p.Data, "client_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type keyPairCreatedBuilder struct{}

func NewKeyPairCreatedBuilder() event2.PrimitiveEventBuilder { return &keyPairCreatedBuilder{} }
func (b *keyPairCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return KeyPairCreated{
		Algorithm: str(p.Data, "algorithm"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

func str(data map[string]any, key string) string {
	v, _ := data[key].(string)
	return v
}
