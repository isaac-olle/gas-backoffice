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
	id        ids.SubscriptionId
	name      string
	createdBy ids.UserId
	memberIds *vo.SetVo[ids.UserId]
	version   int
}

func NewSubscription(id ids.SubscriptionId, name string, createdBy ids.UserId) (*Subscription, error) {
	if name == "" {
		return nil, app_errors.BadRequestError(ErrInvalidSubscriptionName)
	}
	s := &Subscription{Root: aggregate.NewRoot(subscriptionAggregateName), id: id, name: name, createdBy: createdBy, memberIds: vo.NewEmptySetVo[ids.UserId](), version: 1}
	s.Record(backoffice.NewSubscriptionCreated(id.String(), name, createdBy.String()))
	return s, nil
}

func FromPrimitives(p SubscriptionPrimitives) (*Subscription, error) {
	id, err := ids.NewSubscriptionIdFromString(p.Id)
	if err != nil {
		return nil, err
	}
	creator, err := ids.NewUserIdFromString(p.CreatedBy)
	if err != nil {
		return nil, err
	}
	members := vo.NewEmptySetVo[ids.UserId]()
	for _, raw := range p.MemberIds {
		member, e := ids.NewUserIdFromString(raw)
		if e != nil {
			return nil, e
		}
		members.Add(member)
	}
	version := p.Version
	if version == 0 {
		version = 1
	}
	return &Subscription{Root: aggregate.NewRoot(subscriptionAggregateName), id: id, name: p.Name, createdBy: creator, memberIds: members, version: version}, nil
}

func (s *Subscription) AddMember(userId ids.UserId) {
	if s.memberIds.Contains(userId) {
		return
	}
	s.memberIds.Add(userId)
	s.Record(backoffice.NewSubscriptionMemberAdded(s.id.String(), userId.String()))
}

func (s *Subscription) RemoveMember(userId ids.UserId) {
	if !s.memberIds.Contains(userId) {
		return
	}
	s.memberIds.Remove(userId)
	s.Record(backoffice.NewSubscriptionMemberRemoved(s.id.String(), userId.String()))
}

func (s *Subscription) HasMember(userId ids.UserId) bool { return s.memberIds.Contains(userId) }
