package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type UserIdentityId struct{ Uuid }

func NewUserIdentityId() UserIdentityId { return UserIdentityId{GenerateUuid()} }

func NewUserIdentityIdFromString(raw string) (UserIdentityId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return UserIdentityId{}, fmt.Errorf("invalid user identity id: %w", err)
	}
	return UserIdentityId{NewUuid(parsed)}, nil
}
