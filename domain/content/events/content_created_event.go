package events

import "contentgit/ports/out/persistance/eventsourcing"

const (
	ContentCreatedEventType eventsourcing.EventType = "CONTENT_CREATED_V1"
)

type ContentCreatedEventV1 struct {
	Content     map[string]any `json:"content"`
	ContentType string         `json:"contentType"`
	Metadata    *string        `json:"-"`
}
