package error_managers

import (
	"context"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	shared_logger "github.com/isaac-olle/gas-backoffice/internal/shared/logger"
)

type LogErrorPresenter struct {
	logger shared_logger.Logger
}

func NewLogErrorPresenter(logger shared_logger.Logger) *LogErrorPresenter {
	return &LogErrorPresenter{logger: logger}
}

func (p *LogErrorPresenter) Manage(_ context.Context, err *app_errors.AppError) {
	trace := any(nil)
	if err.StackTrace != nil {
		trace = err.StackTrace.ToJSON()
	}

	if p.logger != nil {
		p.logger.Error(
			"HTTP Error",
			shared_logger.Field{Key: "error_id", Value: err.ErrorId},
			shared_logger.Field{Key: "error_type", Value: err.Type.String()},
			shared_logger.Field{Key: "error_code", Value: string(err.Code)},
			shared_logger.Field{Key: "description", Value: err.Description},
			shared_logger.Field{Key: "request_id", Value: err.RequestId},
			shared_logger.Field{Key: "occurred_at", Value: err.OccurredAt.Format(time.RFC3339)},
			// Pass structured values directly; zap.Any encodes them as nested
			// JSON. Pre-marshalling to a string would double-encode (escaped \").
			shared_logger.Field{Key: "details", Value: err.Details},
			shared_logger.Field{Key: "cause", Value: func() []string {
				causes := make([]string, 0, len(err.Cause))
				for _, c := range err.Cause {
					causes = append(causes, c.Error())
				}
				return causes
			}()},
			shared_logger.Field{Key: "trace", Value: trace},
		)
	}
}
