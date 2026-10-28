package _event_handlers

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	userInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberProjectionHandler struct {
	users users.Repository
}

func NewMemberProjectionHandler(db *pgxpool.Pool) *MemberProjectionHandler {
	return &MemberProjectionHandler{users: userInfra.NewPostgresUserRepository(db)}
}
func (h *MemberProjectionHandler) Add(ctx context.Context, e backoffice.SubscriptionMemberAdded) error {
	role, err := vo.NewSubscriptionRole(e.SubscriptionRole)
	if err != nil {
		return err
	}
	return h.mutate(ctx, e.AggregateID(), e.UserId, &role)
}
func (h *MemberProjectionHandler) Remove(ctx context.Context, e backoffice.SubscriptionMemberRemoved) error {
	return h.mutate(ctx, e.AggregateID(), e.UserId, nil)
}
func (h *MemberProjectionHandler) mutate(ctx context.Context, subscriptionId, userId string, role *vo.SubscriptionRole) error {
	sid, err := ids.NewSubscriptionIdFromString(subscriptionId)
	if err != nil {
		return err
	}
	uid, err := ids.NewUserIdFromString(userId)
	if err != nil {
		return err
	}
	user, err := h.users.SearchById(ctx, uid)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("subscription member user not found: %s", userId)
	}
	if role != nil {
		user.AddMembership(vo.NewSubscriptionMembership(sid, *role))
	} else {
		user.RemoveMembership(sid)
	}
	return h.users.Save(ctx, user)
}
