package auth

import "time"

import "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"

type SessionPrimitives struct {
	SessionId       string
	UserId          string
	AdminId         string
	Role            string
	ExpiresAt       time.Time
	SubscriptionIds []string
	Memberships     []vo.SubscriptionMembership
	Picture         string
}

func (s *Session) ToPrimitives() *SessionPrimitives {
	subscriptionIds := make([]string, 0, len(s.subscriptionIds))
	for _, id := range s.subscriptionIds {
		subscriptionIds = append(subscriptionIds, id.String())
	}
	return &SessionPrimitives{
		SessionId:       s.sessionId.String(),
		UserId:          s.userId.String(),
		AdminId:         s.adminId.String(),
		Role:            s.role.String(),
		ExpiresAt:       s.expiresAt,
		SubscriptionIds: subscriptionIds,
		Memberships:     s.memberships,
		Picture:         s.picture,
	}
}
