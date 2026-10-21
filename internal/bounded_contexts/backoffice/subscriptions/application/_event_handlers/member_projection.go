package _event_handlers

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	adminInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	tenantInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberProjectionHandler struct {
	admins  admins.Repository
	tenants tenants.Repository
}

func NewMemberProjectionHandler(db *pgxpool.Pool) *MemberProjectionHandler {
	return &MemberProjectionHandler{admins: adminInfra.NewPostgresAdminRepository(db), tenants: tenantInfra.NewPostgresTenantRepository(db)}
}
func (h *MemberProjectionHandler) Add(ctx context.Context, e backoffice.SubscriptionMemberAdded) error {
	return h.mutate(ctx, e.AggregateID(), e.UserId, true)
}
func (h *MemberProjectionHandler) Remove(ctx context.Context, e backoffice.SubscriptionMemberRemoved) error {
	return h.mutate(ctx, e.AggregateID(), e.UserId, false)
}
func (h *MemberProjectionHandler) mutate(ctx context.Context, subscriptionId, userId string, add bool) error {
	sid, err := ids.NewSubscriptionIdFromString(subscriptionId)
	if err != nil {
		return err
	}
	uid, err := ids.NewUserIdFromString(userId)
	if err != nil {
		return err
	}
	admin, err := h.admins.SearchById(ctx, uid)
	if err != nil {
		return err
	}
	if admin != nil {
		if add {
			admin.AddSubscriptionId(sid)
		} else {
			admin.RemoveSubscriptionId(sid)
		}
		return h.admins.Save(ctx, admin)
	}
	tenant, err := h.tenants.SearchById(ctx, uid)
	if err != nil {
		return err
	}
	if tenant == nil {
		return fmt.Errorf("subscription member user not found: %s", userId)
	}
	if add {
		tenant.AddSubscriptionId(sid)
	} else {
		tenant.RemoveSubscriptionId(sid)
	}
	return h.tenants.Save(ctx, tenant)
}
