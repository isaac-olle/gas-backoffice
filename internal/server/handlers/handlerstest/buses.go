package handlerstest

import (
	"github.com/stretchr/testify/mock"

	commandmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/command/mocks"
	querymocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/query/mocks"
)

func NewCommandBus(err error) *commandmocks.MockCommandBus {
	b := &commandmocks.MockCommandBus{}
	b.On("Dispatch", mock.Anything, mock.Anything).Return(err)
	return b
}

func NewQueryBus(result any, err error) *querymocks.MockQueryBus {
	b := &querymocks.MockQueryBus{}
	b.On("Dispatch", mock.Anything, mock.Anything).Return(result, err)
	return b
}
