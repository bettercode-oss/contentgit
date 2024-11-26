package events

import (
	"contentgit/ports/out/persistance/eventsourcing"
	"time"
)

const (
	FieldUpdatedEventType eventsourcing.EventType = "CONTENT_FIELD_UPDATED_V1"
)

type FieldUpdatedEventV1 struct {
	FieldName     string    `json:"fieldName"`
	BeforeValue   any       `json:"beforeValue"`
	AfterValue    any       `json:"afterValue"`
	CreatedById   string    `json:"createdById"`
	CreatedByName string    `json:"createdByName"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Metadata      *string   `json:"-"`
}
