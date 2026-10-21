package error_logs

import "context"

type Repository interface {
	Save(ctx context.Context, errorLog ErrorLog) error
	Latest(ctx context.Context, limit int) ([]ErrorLog, error)
}
