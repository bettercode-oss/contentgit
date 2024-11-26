package eventsourcing

import (
	"contentgit/foundation"
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type EventRepository struct {
}

func (r EventRepository) Save(ctx context.Context, events []Event) error {
	db := foundation.ContextProvider().GetDB(ctx)
	if err := db.Create(events).Error; err != nil {
		return errors.Wrap(err, "(SaveEvents) tx.Exec err")
	}
	return nil
}

func (r EventRepository) FindByAggregateIdAndVersion(ctx context.Context, aggregateID string, versionFrom uint64) ([]Event, error) {
	db := foundation.ContextProvider().GetDB(ctx)
	events := make([]Event, 0)

	if err := db.Where("aggregate_id = ? AND version > ?", aggregateID, versionFrom).Order("version ASC").Find(&events).Error; err != nil {
		return nil, errors.Wrap(err, "(FindByAggregateIdAndVersion) db.Query err")
	}

	return events, nil
}

func (r EventRepository) FindOneByAggregateId(ctx context.Context, aggregateID string, forUpdate bool) (*Event, error) {
	db := foundation.ContextProvider().GetDB(ctx)

	db = db.Where("aggregate_id = ?", aggregateID)
	if forUpdate {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	event := Event{}
	if err := db.First(&event).Error; err != nil {
		return nil, err
	}

	return nil, nil
}

func (r EventRepository) FindByAggregateId(ctx context.Context, aggregateID string) ([]Event, error) {
	db := foundation.ContextProvider().GetDB(ctx)
	events := make([]Event, 0)

	if err := db.Where("aggregate_id = ?", aggregateID).Order("version ASC").Find(&events).Error; err != nil {
		return nil, errors.Wrap(err, "(FindByAggregateId) db.Query err")
	}

	return events, nil
}

type SnapshotRepository struct {
}

func (r SnapshotRepository) Save(ctx context.Context, snapshot *Snapshot) error {
	db := foundation.ContextProvider().GetDB(ctx)

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "aggregate_id"}},
		DoUpdates: clause.Assignments(map[string]any{"data": snapshot.State, "version": snapshot.Version, "updated_at": time.Now()}),
	}).Create(&snapshot).Error; err != nil {
		return errors.Wrap(err, "(Save Snapshot) tx.Exec err")
	}

	return nil
}

func (r SnapshotRepository) FindOneByAggregateId(ctx context.Context, aggregateId string) (*Snapshot, error) {
	db := foundation.ContextProvider().GetDB(ctx)

	snapshot := Snapshot{}
	if err := db.Where("aggregate_id = ?", aggregateId).First(&snapshot).Error; err != nil {
		return nil, err
	}

	return &snapshot, nil
}
