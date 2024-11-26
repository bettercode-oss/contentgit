package eventsourcing

import (
	"contentgit/ports/out/persistance/eventsourcing/serializer"
	"context"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Load eventsourcing.Aggregate events using snapshots with given frequency
func (m *rdbEventStore) Load(ctx context.Context, aggregate Aggregate) error {
	snapshot, err := m.GetSnapshot(ctx, aggregate.GetID())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if snapshot != nil {
		if err := serializer.Unmarshal(snapshot.State, aggregate); err != nil {
			log.Info("(Load) serializer.Unmarshal err", err)
			return errors.Wrap(err, "json.Unmarshal")
		}

		err := m.loadAggregateEventsByVersion(ctx, aggregate)
		if err != nil {
			return err
		}

		log.Info(fmt.Sprintf("(Load Aggregate By Version) aggregate: %s", aggregate.String()))
		//span.LogFields(log.String("aggregate with events", aggregate.String()))
		return nil
	}

	err = m.loadEvents(ctx, aggregate)
	if err != nil {
		return err
	}

	log.Printf("(Load Aggregate): aggregate: %s", aggregate.String())
	//span.LogFields(log.String("aggregate with events", aggregate.String()))
	return nil
}

// Save eventsourcing.Aggregate events using snapshots with given frequency
func (m *rdbEventStore) Save(ctx context.Context, aggregate Aggregate) error {
	if len(aggregate.GetChanges()) == 0 {
		log.Info("(Save) aggregate.GetChanges()) == 0")
		return nil
	}

	changes := aggregate.GetChanges()
	events := make([]Event, 0, len(changes))

	for i := range changes {
		event, err := m.serializer.SerializeEvent(aggregate, changes[i])
		if err != nil {
			return errors.Wrap(err, "(Save) serializer.SerializeEvent err")
		}
		events = append(events, event)
	}

	if err := m.saveEventsTx(ctx, events); err != nil {
		return errors.Wrap(err, "saveEventsTx")
	}

	if aggregate.GetVersion()%snapshotFrequency == 0 {
		aggregate.ToSnapshot()
		if err := m.saveSnapshotTx(ctx, aggregate); err != nil {
			return errors.Wrap(err, "saveSnapshotTx")
		}
	}

	if err := m.processEvents(ctx, events); err != nil {
		return errors.Wrap(err, "processEvents")
	}

	log.Info(fmt.Sprintf("(Save Aggregate): aggregate: %s", aggregate.String()))

	return nil
}
