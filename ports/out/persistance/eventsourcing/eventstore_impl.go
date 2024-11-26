package eventsourcing

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	eventsCapacity    = 10
	snapshotFrequency = 5
)

type rdbEventStore struct {
	eventBus           EventsBus
	serializer         Serializer
	eventRepository    *EventRepository
	snapshotRepository *SnapshotRepository
}

func NewRdbEventStore(eventBus EventsBus, serializer Serializer,
	eventRepository *EventRepository, snapshotRepository *SnapshotRepository) *rdbEventStore {
	return &rdbEventStore{eventBus: eventBus, serializer: serializer, eventRepository: eventRepository, snapshotRepository: snapshotRepository}
}

// SaveEvents save aggregate uncommitted events as one batch and process with event bus using transaction
func (m *rdbEventStore) SaveEvents(ctx context.Context, events []Event) error {
	if err := m.handleConcurrency(ctx, events); err != nil {
		return errors.Wrap(err, "(SaveEvents) Concurrency err")
	}

	if err := m.eventRepository.Save(ctx, events); err != nil {
		return errors.Wrap(err, "(SaveEvents) tx.Exec err")
	}

	if err := m.processEvents(ctx, events); err != nil {
		return errors.Wrap(err, "(SaveEvents) processEvents err")
	}

	return nil
}

// LoadEvents load aggregate events by id
func (m *rdbEventStore) LoadEvents(ctx context.Context, aggregateID string) ([]Event, error) {
	return m.eventRepository.FindByAggregateId(ctx, aggregateID)
}

// LoadEvents load aggregate events by id
func (m *rdbEventStore) loadEvents(ctx context.Context, aggregate Aggregate) error {
	events, err := m.eventRepository.FindByAggregateId(ctx, aggregate.GetID())
	if err != nil {
		return errors.Wrap(err, "(loadEvents) db.Query err")
	}

	for _, event := range events {
		deserializedEvent, err := m.serializer.DeserializeEvent(event)
		if err != nil {
			return errors.Wrap(err, "(loadEvents) serializer.DeserializeEvent err")
		}

		if err := aggregate.RaiseEvent(deserializedEvent); err != nil {
			return errors.Wrap(err, "(loadEvents) aggregate.RaiseEvent err")
		}
	}

	return nil
}

// Exists check for exists aggregate by id
func (m *rdbEventStore) Exists(ctx context.Context, aggregateID string) (bool, error) {
	_, err := m.eventRepository.FindOneByAggregateId(ctx, aggregateID, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, errors.Wrap(err, "(Exists) db.QueryRow err")
	}

	return true, nil
}

func (m *rdbEventStore) loadEventsByVersion(ctx context.Context, aggregateID string, versionFrom uint64) ([]Event, error) {
	return m.eventRepository.FindByAggregateIdAndVersion(ctx, aggregateID, versionFrom)
}

func (m *rdbEventStore) loadAggregateEventsByVersion(ctx context.Context, aggregate Aggregate) error {
	events, err := m.eventRepository.FindByAggregateIdAndVersion(ctx, aggregate.GetID(), aggregate.GetVersion())
	if err != nil {
		return err
	}

	for _, event := range events {
		deserializedEvent, err := m.serializer.DeserializeEvent(event)
		if err != nil {
			return errors.Wrap(err, "(loadAggregateEventsByVersion) serializer.DeserializeEvent err")
		}

		if err := aggregate.RaiseEvent(deserializedEvent); err != nil {
			return errors.Wrap(err, "(loadAggregateEventsByVersion) aggregate.RaiseEvent err")
		}
	}

	return nil
}

func (m *rdbEventStore) loadEventsByVersionTx(ctx context.Context, aggregateID string, versionFrom uint64) ([]Event, error) {
	return m.eventRepository.FindByAggregateIdAndVersion(ctx, aggregateID, versionFrom)
}

func (m *rdbEventStore) saveEventsTx(ctx context.Context, events []Event) error {
	if err := m.handleConcurrency(ctx, events); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("(saveEventsTx) AggregateID: %s, AggregateVersion: %v, AggregateType: %s", events[0].GetAggregateID(), events[0].GetVersion(), events[0].GetAggregateType()))
	return m.eventRepository.Save(ctx, events)
}

func (m *rdbEventStore) saveSnapshotTx(ctx context.Context, aggregate Aggregate) error {
	snapshot, err := NewSnapshotFromAggregate(aggregate)
	if err != nil {
		return errors.Wrap(err, "(saveSnapshotTx) NewSnapshotFromAggregate err")
	}
	log.Info(fmt.Sprintf("(saveSnapshotTx) snapshot: %s", snapshot.String()))
	return m.snapshotRepository.Save(ctx, snapshot)
}

func (m *rdbEventStore) processEvents(ctx context.Context, events []Event) error {
	return m.eventBus.ProcessEvents(ctx, events)
}

func (m *rdbEventStore) handleConcurrency(ctx context.Context, events []Event) error {
	_, err := m.eventRepository.FindOneByAggregateId(ctx, events[0].GetAggregateID(), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return errors.Wrap(err, "(handleConcurrency) tx.Exec")
	}

	return nil
}
