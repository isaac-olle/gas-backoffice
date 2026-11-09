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

func NewCommandExecutor(err error) *commandmocks.MockCommandExecutor {
	e := &commandmocks.MockCommandExecutor{}
	e.On("SingleExecution", mock.Anything, mock.Anything).Return(err)
	e.On("MultipleExecution", mock.Anything, mock.Anything, mock.Anything).Return(err)
	e.On("MultipleExecution", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(err)
	return e
}

func NewQueryBus(result any, err error) *querymocks.MockQueryBus {
	b := &querymocks.MockQueryBus{}
	b.On("Dispatch", mock.Anything, mock.Anything).Return(result, err)
	return b
}
