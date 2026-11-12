package event_logs

import "context"

type Repository interface {
	Latest(context.Context, int) ([]EventLog, error)
}
