package vos

import "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

var ErrInvalidScopes app_errors.ErrorCode = "INVALID_SCOPES"

type Scopes struct{ values []string }

func NewScopes(raw []string) (Scopes, error) {
	if len(raw) == 0 {
		return Scopes{}, app_errors.BadRequestError(ErrInvalidScopes).
			WithDescription("at least one scope is required")
	}
	for _, s := range raw {
		if s == "" {
			return Scopes{}, app_errors.BadRequestError(ErrInvalidScopes).
				WithDescription("scope cannot be empty")
		}
	}
	cp := make([]string, len(raw))
	copy(cp, raw)
	return Scopes{values: cp}, nil
}

func (s Scopes) ToSlice() []string {
	cp := make([]string, len(s.values))
	copy(cp, s.values)
	return cp
}
