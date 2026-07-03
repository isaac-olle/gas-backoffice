package error_managers

import (
	"encoding/json"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	shared_logger "github.com/isaac-olle/gas-backoffice/internal/shared/logger"
)

type LogErrorPresenter struct {
	logger shared_logger.Logger
}

func NewLogErrorPresenter(logger shared_logger.Logger) *LogErrorPresenter {
	return &LogErrorPresenter{logger: logger}
}

func (p *LogErrorPresenter) Manage(err *app_errors.AppError) {
	trace := any(nil)
	if err.StackTrace != nil {
		trace = err.StackTrace.ToJSON()
	}

	if p.logger != nil {
		traceJSON, _ := json.Marshal(trace)
		p.logger.Error(
			"HTTP Error",
			shared_logger.Field{Key: "error_id", Value: err.ErrorId},
			shared_logger.Field{Key: "error_type", Value: err.Type.String()},
			shared_logger.Field{Key: "error_code", Value: string(err.Code)},
			shared_logger.Field{Key: "description", Value: err.Description},
			shared_logger.Field{Key: "details", Value: func() string {
				if err.Details == nil {
					return ""
				}
				b, _ := json.Marshal(err.Details)
				return string(b)
			}()},
			shared_logger.Field{Key: "cause", Value: func() string {
				if err.Cause != nil {
					return err.Cause.Error()
				}
				return ""
			}()},
			shared_logger.Field{Key: "trace", Value: string(traceJSON)},
		)
	}
}
