package user_projection

import "context"

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*UserProjection, error)
}
