package eventsourcing

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type EventsBus interface {
	ProcessEvents(ctx context.Context, events []Event) error
}

type EventsBusMock struct {
	mock.Mock
}

func (e *EventsBusMock) ProcessEvents(ctx context.Context, events []Event) error {
	args := e.Called(ctx, events)
	return args.Error(0)
}
