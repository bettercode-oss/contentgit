package content

import (
	"contentgit/domain/content/events"
	"contentgit/ports/out/persistance/eventsourcing"
	"contentgit/ports/out/persistance/eventsourcing/serializer"
	"github.com/pkg/errors"
)

var (
	ErrInvalidEvent = errors.New("invalid event")
)

type eventSerializer struct {
}

func NewEventSerializer() *eventSerializer {
	return &eventSerializer{}
}

func (s *eventSerializer) SerializeEvent(aggregate eventsourcing.Aggregate, event any) (eventsourcing.Event, error) {
	eventJson, err := serializer.Marshal(event)
	if err != nil {
		return eventsourcing.Event{}, errors.Wrapf(err, "serializer.Marshal aggregateID: %s", aggregate.GetID())
	}

	switch evt := event.(type) {
	case *events.ContentCreatedEventV1:
		return eventsourcing.NewEvent(aggregate, events.ContentCreatedEventType, eventJson, evt.Metadata), nil
	case *events.FieldUpdatedEventV1:
		return eventsourcing.NewEvent(aggregate, events.FieldUpdatedEventType, eventJson, evt.Metadata), nil
	case *events.FieldCommentAddedEventV1:
		return eventsourcing.NewEvent(aggregate, events.FieldCommentAddedEventType, eventJson, evt.Metadata), nil
	default:
		return eventsourcing.Event{}, errors.Wrapf(ErrInvalidEvent, "aggregateID: %s, type: %T", aggregate.GetID(), event)
	}
}

func (s *eventSerializer) DeserializeEvent(event eventsourcing.Event) (any, error) {
	switch event.GetEventType() {
	case events.ContentCreatedEventType:
		return deserializeEvent(event, new(events.ContentCreatedEventV1))
	case events.FieldUpdatedEventType:
		return deserializeEvent(event, new(events.FieldUpdatedEventV1))
	case events.FieldCommentAddedEventType:
		return deserializeEvent(event, new(events.FieldCommentAddedEventV1))
	default:
		return nil, errors.Wrapf(ErrInvalidEvent, "type: %s", event.GetEventType())
	}
}

func deserializeEvent(event eventsourcing.Event, targetEvent any) (any, error) {
	if err := event.GetJsonData(&targetEvent); err != nil {
		return nil, errors.Wrapf(err, "event.GetJsonData type: %s", event.GetEventType())
	}
	return targetEvent, nil
}
