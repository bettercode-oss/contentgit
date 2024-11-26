package eventsourcing

import (
	"contentgit/ports/out/persistance/eventsourcing/serializer"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Snapshot Event Sourcing Snapshotting is an optimisation that reduces time spent on reading event from an event store.
type Snapshot struct {
	gorm.Model
	AggregateId string        `gorm:"type:varchar(100);not null;uniqueIndex:idx_snapshot_unique;index:idx_snapshot_aggregate_id_version"`
	TenantId    string        `gorm:"type:varchar(100);not null"`
	Type        AggregateType `gorm:"column:aggregate_type;type:varchar(250);not null"`
	State       string        `gorm:"column:data;type:jsonb"`
	Version     uint64        `gorm:"not null;index:idx_snapshot_aggregate_id_version"`
}

func (*Snapshot) TableName() string {
	return "snapshots"
}
func (s *Snapshot) String() string {
	return fmt.Sprintf("AggregateID: %s, TenantId: %s, Type: %s, StateSize: %d, Version: %d",
		s.AggregateId,
		s.TenantId,
		string(s.Type),
		len(s.State),
		s.Version,
	)
}

// NewSnapshotFromAggregate create new Snapshot from the Aggregate state.
func NewSnapshotFromAggregate(aggregate Aggregate) (*Snapshot, error) {
	aggregateJson, err := serializer.Marshal(aggregate)
	if err != nil {
		return nil, errors.Wrapf(err, "serializer.Marshal aggregateID: %s", aggregate.GetID())
	}

	return &Snapshot{
		AggregateId: aggregate.GetID(),
		TenantId:    aggregate.GetTenantId(),
		Type:        aggregate.GetType(),
		State:       aggregateJson,
		Version:     aggregate.GetVersion(),
	}, nil
}
