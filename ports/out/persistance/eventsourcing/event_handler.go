package eventsourcing

import "context"

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	GetAggregateType() AggregateType
}
