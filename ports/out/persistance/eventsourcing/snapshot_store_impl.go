package eventsourcing

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SaveSnapshot save eventsourcing.Aggregate snapshot
func (m *rdbEventStore) SaveSnapshot(ctx context.Context, aggregate Aggregate) error {
	snapshot, err := NewSnapshotFromAggregate(aggregate)
	if err != nil {
		return errors.Wrap(err, "NewSnapshotFromAggregate")
	}

	return m.snapshotRepository.Save(ctx, snapshot)
}

// GetSnapshot load eventsourcing.Aggregate snapshot
func (m *rdbEventStore) GetSnapshot(ctx context.Context, id string) (*Snapshot, error) {
	snapshot, err := m.snapshotRepository.FindOneByAggregateId(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "db.QueryRow")
	}

	log.Info(fmt.Sprintf("(GetSnapshot) snapshot: %s", snapshot.String()))
	return snapshot, nil
}
