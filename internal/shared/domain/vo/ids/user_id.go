package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type UserId struct{ Uuid }

func NewUserId() UserId { return UserId{GenerateUuid()} }

func NewUserIdFromString(raw string) (UserId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return UserId{}, fmt.Errorf("invalid user id: %w", err)
	}
	return UserId{NewUuid(parsed)}, nil
}
