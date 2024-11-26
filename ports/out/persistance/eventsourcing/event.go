package eventsourcing

import (
	"contentgit/ports/out/persistance/eventsourcing/serializer"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// EventType is the type of any event, used as its unique identifier.
type EventType string

// Event is an internal representation of an event, returned when the Aggregate
// uses NewEvent to create a new event. The events loaded from the db is
// represented by each DBs internal event type, implementing Event.
type Event struct {
	gorm.Model
	AggregateID   string        `gorm:"type:varchar(100);not null;uniqueIndex:idx_unique;index:idx_aggregate_id_version"`
	TenantId      string        `gorm:"type:varchar(100);not null"`
	AggregateType AggregateType `gorm:"type:varchar(250);not null"`
	EventType     EventType     `gorm:"type:varchar(250);not null"`
	Data          string        `gorm:"type:jsonb"`
	Metadata      *string       `gorm:"type:jsonb"`
	Version       uint64        `gorm:"not null;uniqueIndex:idx_unique;index:idx_aggregate_id_version"`
}

func (*Event) TableName() string {
	return "events"
}

// NewBaseEvent new base Event constructor with configured EventID, Aggregate properties and Timestamp.
func NewBaseEvent(aggregate Aggregate, eventType EventType) Event {
	return Event{
		AggregateType: aggregate.GetType(),
		TenantId:      aggregate.GetTenantId(),
		AggregateID:   aggregate.GetID(),
		Version:       aggregate.GetVersion(),
		EventType:     eventType,
	}
}

func NewEvent(aggregate Aggregate, eventType EventType, data string, metadata *string) Event {
	return Event{
		AggregateID:   aggregate.GetID(),
		TenantId:      aggregate.GetTenantId(),
		EventType:     eventType,
		AggregateType: aggregate.GetType(),
		Version:       aggregate.GetVersion(),
		Data:          data,
		Metadata:      metadata,
	}
}

// GetEventID get EventID of the Event.
func (e *Event) GetEventID() uint {
	return e.ID
}

// GetTimeStamp get timestamp of the Event.
func (e *Event) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// GetData The data attached to the Event serialized to bytes.
func (e *Event) GetData() string {
	return e.Data
}

// SetData add the data attached to the Event serialized to bytes.
func (e *Event) SetData(data string) *Event {
	e.Data = data
	return e
}

// GetJsonData json unmarshal data attached to the Event.
func (e *Event) GetJsonData(data interface{}) error {
	return serializer.Unmarshal(e.GetData(), data)
}

// SetJsonData serialize to json and set data attached to the Event.
func (e *Event) SetJsonData(data interface{}) error {
	dataJson, err := serializer.Marshal(data)
	if err != nil {
		return err
	}

	e.Data = dataJson
	return nil
}

// GetEventType returns the EventType of the event.
func (e *Event) GetEventType() EventType {
	return e.EventType
}

// GetAggregateType is the AggregateType that the Event can be applied to.
func (e *Event) GetAggregateType() AggregateType {
	return e.AggregateType
}

// SetAggregateType set the AggregateType that the Event can be applied to.
func (e *Event) SetAggregateType(aggregateType AggregateType) {
	e.AggregateType = aggregateType
}

// GetAggregateID is the AggregateID of the Aggregate that the Event belongs to
func (e *Event) GetAggregateID() string {
	return e.AggregateID
}

// GetVersion is the version of the Aggregate after the Event has been applied.
func (e *Event) GetVersion() uint64 {
	return e.Version
}

// SetVersion set the version of the Aggregate.
func (e *Event) SetVersion(aggregateVersion uint64) {
	e.Version = aggregateVersion
}

// GetMetadata is app-specific metadata such as request AggregateID, originating user etc.
func (e *Event) GetMetadata() *string {
	return e.Metadata
}

// SetMetadata add app-specific metadata serialized as json for the Event.
func (e *Event) SetMetadata(metadata interface{}) error {
	metadataJson, err := serializer.Marshal(metadata)
	if err != nil {
		return err
	}

	e.Metadata = &metadataJson
	return nil
}

// GetJsonMetadata unmarshal app-specific metadata serialized as json for the Event.
func (e *Event) GetJsonMetadata(metaData interface{}) error {
	return serializer.Unmarshal(*e.GetMetadata(), metaData)
}

// GetString A string representation of the Event.
func (e *Event) GetString() string {
	return fmt.Sprintf("event: %+v", e)
}

func (e *Event) GetTenantId() string {
	return e.TenantId
}

func (e *Event) String() string {
	return fmt.Sprintf("(Event) AggregateID: %s, TenantId: %s, Version: %d, EventType: %s, AggregateType: %s, Metadata: %s, TimeStamp: %s, EventID: %d",
		e.AggregateID,
		e.TenantId,
		e.Version,
		e.EventType,
		e.AggregateType,
		*e.Metadata,
		e.CreatedAt,
		e.ID,
	)
}
