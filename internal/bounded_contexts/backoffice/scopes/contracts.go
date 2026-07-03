package scopes

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type Repository interface {
	persistance.Repository[*Scope, ids.ScopeId]
}
