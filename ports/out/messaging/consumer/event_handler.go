package consumer

import (
	"contentgit/ports/out/persistance/eventsourcing"
	"context"
)

type EventHandler interface {
	Handle(ctx context.Context, event eventsourcing.Event) error
	GetAggregateType() eventsourcing.AggregateType
}
