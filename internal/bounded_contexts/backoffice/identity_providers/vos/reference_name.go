package vos

import (
	"regexp"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrInvalidReferenceName app_errors.ErrorCode = "INVALID_REFERENCE_NAME"
	referenceNamePattern                          = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

type ReferenceName struct{ value string }

func NewReferenceName(raw string) (ReferenceName, error) {
	if !referenceNamePattern.MatchString(raw) {
		return ReferenceName{}, app_errors.BadRequestError(ErrInvalidReferenceName).
			WithDescription("reference_name must be a lowercase slug (letters, digits, hyphens)")
	}
	return ReferenceName{value: raw}, nil
}

func (r ReferenceName) String() string { return r.value }
