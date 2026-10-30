package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	verifyJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type RefreshTokenVerifier struct {
	queryBus query.Bus
	issuer   string
}

func NewRefreshTokenVerifier(queryBus query.Bus, issuer string) *RefreshTokenVerifier {
	return &RefreshTokenVerifier{queryBus: queryBus, issuer: issuer}
}

func (v *RefreshTokenVerifier) Verify(ctx context.Context, token string) (*auth.RefreshTokenClaims, error) {
	claims, err := query.DispatchTyped[map[string]any](ctx, v.queryBus,
		verifyJwt.NewQuery(token).
			WithExpectedClaim("iss", v.issuer).
			WithExpectedClaim("aud", v.issuer).
			WithExpectedClaim("typ", "backoffice_refresh").
			WithVerifyExpiration().
			WithClaimsToReturn(map[string]any{
				"sub":         nil,
				"sid":         nil,
				"system_role": nil,
				"picture":     nil,
				"memberships": nil,
			}),
	)
	if err != nil {
		return nil, err
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("refresh token has an invalid sub claim")
	}
	role, ok := claims["system_role"].(string)
	if !ok || role == "" {
		return nil, fmt.Errorf("refresh token has an invalid system_role claim")
	}
	sessionIDRaw, ok := claims["sid"].(string)
	if !ok || sessionIDRaw == "" {
		return nil, fmt.Errorf("refresh token has an invalid sid claim")
	}
	sessionID, err := ids.NewSessionIdFromString(sessionIDRaw)
	if err != nil {
		return nil, fmt.Errorf("refresh token has an invalid sid claim: %w", err)
	}

	memberships := make([]vo.SubscriptionMembership, 0)
	if raw, exists := claims["memberships"]; exists && raw != nil {
		encoded, err := json.Marshal(raw)
		if err != nil {
			return nil, fmt.Errorf("marshal refresh token memberships: %w", err)
		}
		if err = json.Unmarshal(encoded, &memberships); err != nil {
			return nil, fmt.Errorf("unmarshal refresh token memberships: %w", err)
		}
	}

	picture, _ := claims["picture"].(string)
	return &auth.RefreshTokenClaims{SessionID: sessionID, UserID: userID, SystemRole: role, Picture: picture, Memberships: memberships}, nil
}

var _ auth.RefreshTokenVerifier = (*RefreshTokenVerifier)(nil)
