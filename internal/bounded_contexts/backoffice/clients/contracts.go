package clients

import (
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/persistance"
)

type Repository interface {
	persistance.Repository[*Client, ids.ClientId]
}
