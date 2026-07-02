package event

import "context"

type Subscriber interface {
	On(context.Context, Event) error
	SubscribedTo() Type
}

type BaseSubscriber[T Event] struct {
	handler      func(ctx context.Context, event T) error
	subscribedTo Type
}

func NewBaseSubscriber[T Event](handler func(ctx context.Context, event T) error) *BaseSubscriber[T] {
	var zero T
	return &BaseSubscriber[T]{
		handler:      handler,
		subscribedTo: zero.Type(),
	}
}

func (s *BaseSubscriber[T]) On(ctx context.Context, evt Event) error {
	typedEvt, ok := evt.(T)
	if !ok {
		return nil
	}
	return s.handler(ctx, typedEvt)
}

func (s *BaseSubscriber[T]) SubscribedTo() Type {
	return s.subscribedTo
}

type InMemorySubscriber[T Event] struct {
	*BaseSubscriber[T]
}

func NewInMemorySubscriber[T Event](handler func(ctx context.Context, event T) error) *InMemorySubscriber[T] {
	return &InMemorySubscriber[T]{NewBaseSubscriber[T](handler)}
}
