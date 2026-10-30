package subscriptions

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var subscriptionAggregateName aggregate.Name = aggregate.Name("backoffice.subscription")

type Subscription struct {
	*aggregate.Root
	id       ids.SubscriptionId
	name     string
	ownerId  ids.UserId
	isSystem bool
	members  map[ids.UserId]vo.SubscriptionRole
	version  int
}

func NewSubscription(id ids.SubscriptionId, name string, ownerId ids.UserId, isSystem bool) (*Subscription, error) {
	if name == "" {
		return nil, app_errors.BadRequestError(ErrInvalidSubscriptionName)
	}
	s := &Subscription{Root: aggregate.NewRoot(subscriptionAggregateName), id: id, name: name, ownerId: ownerId, isSystem: isSystem, members: make(map[ids.UserId]vo.SubscriptionRole), version: 1}
	s.Record(backoffice.NewSubscriptionCreated(id.String(), name, ownerId.String()))
	return s, nil
}

func (s *Subscription) IsSystemResource() bool { return s.isSystem }

func (s *Subscription) ChangeName(name string) error {
	if name == "" {
		return app_errors.BadRequestError(ErrInvalidSubscriptionName)
	}
	s.name = name
	return nil
}

func FromPrimitives(p SubscriptionPrimitives) (*Subscription, error) {
	id, err := ids.NewSubscriptionIdFromString(p.Id)
	if err != nil {
		return nil, err
	}
	creator, err := ids.NewUserIdFromString(p.OwnerId)
	if err != nil {
		return nil, err
	}
	members := make(map[ids.UserId]vo.SubscriptionRole, len(p.Members))
	for _, raw := range p.Members {
		member, e := ids.NewUserIdFromString(raw.UserId)
		if e != nil {
			return nil, e
		}
		role, e := vo.NewSubscriptionRole(raw.SubscriptionRole)
		if e != nil {
			return nil, e
		}
		members[member] = role
	}
	version := p.Version
	if version == 0 {
		version = 1
	}
	return &Subscription{Root: aggregate.NewRoot(subscriptionAggregateName), id: id, name: p.Name, ownerId: creator, isSystem: p.IsSystem, members: members, version: version}, nil
}

func (s *Subscription) AddMember(userId ids.UserId, role vo.SubscriptionRole) {
	if _, exists := s.members[userId]; exists {
		return
	}
	s.members[userId] = role
	s.Record(backoffice.NewSubscriptionMemberAdded(s.id.String(), userId.String(), role.String()))
}

func (s *Subscription) RemoveMember(userId ids.UserId) {
	if _, exists := s.members[userId]; !exists {
		return
	}
	delete(s.members, userId)
	s.Record(backoffice.NewSubscriptionMemberRemoved(s.id.String(), userId.String()))
}

func (s *Subscription) HasMember(userId ids.UserId) bool {
	_, exists := s.members[userId]
	return exists
}
