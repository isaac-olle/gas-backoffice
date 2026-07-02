package modules

import (
	event2 "github.com/isaac-olle/gas/internal/shared/buses/event"
	eventInfra "github.com/isaac-olle/gas/internal/shared/buses/event/infrastructure"
	shared_logger "github.com/isaac-olle/gas/internal/shared/logger"
	"github.com/isaac-olle/gas/internal/shared/workers"
)

type SubscriptionManager struct {
	subscriptions map[event2.Type][]event2.Subscriber
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[event2.Type][]event2.Subscriber),
	}
}

func (sm *SubscriptionManager) AddInMemorySubscription(subscriber event2.Subscriber) {
	t := subscriber.SubscribedTo()
	sm.subscriptions[t] = append(sm.subscriptions[t], subscriber)
}

func (sm *SubscriptionManager) Build(logger shared_logger.Logger, pool *workers.Pool) event2.Bus {
	return eventInfra.NewInMemoryEventBus(logger, pool, sm.subscriptions)
}
