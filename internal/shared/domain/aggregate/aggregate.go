package aggregate

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Aggregate interface {
	Record(event event.Event)
	PullDomainEvents() []event.Event
	Clear()
	Name() Name
}

type Root struct {
	name         Name
	domainEvents []event.Event
}

func NewRoot(name Name) *Root {
	return &Root{name: name, domainEvents: make([]event.Event, 0)}
}

func (a *Root) Record(evt event.Event) {
	if a.domainEvents == nil {
		a.domainEvents = make([]event.Event, 0)
	}
	a.domainEvents = append(a.domainEvents, evt)
}

func (a *Root) PullDomainEvents() []event.Event {
	events := a.domainEvents
	a.Clear()
	return events
}

func (a *Root) Clear() {
	a.domainEvents = make([]event.Event, 0)
}

func (a *Root) Name() Name {
	return a.name
}
