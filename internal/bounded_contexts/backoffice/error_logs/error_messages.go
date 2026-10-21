package error_logs

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrErrorLogsNotListed app_errors.ErrorCode = "ERROR_LOGS_NOT_LISTED"
	ErrErrorLogNotSaved   app_errors.ErrorCode = "ERROR_LOG_NOT_SAVED"
)
