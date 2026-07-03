package key_pairs

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrKeyPairNotFound   app_errors.ErrorCode = "KEY_PAIR_NOT_FOUND"
	ErrKeyPairNotCreated app_errors.ErrorCode = "KEY_PAIR_NOT_CREATED"
	ErrKeyPairNotSaved   app_errors.ErrorCode = "KEY_PAIR_NOT_SAVED"
	ErrInvalidOwnership  app_errors.ErrorCode = "KEY_PAIR_INVALID_OWNERSHIP"
	ErrInvalidToken      app_errors.ErrorCode = "INVALID_TOKEN"
)
