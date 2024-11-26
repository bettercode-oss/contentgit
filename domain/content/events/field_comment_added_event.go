package events

import (
	"contentgit/ports/out/persistance/eventsourcing"
)

const (
	FieldCommentAddedEventType eventsourcing.EventType = "CONTENT_FIELD_COMMENT_ADDED_V1"
)

type FieldCommentAddedEventV1 struct {
	TenantId      string  `json:"tenantId"`
	FieldName     string  `json:"fieldName"`
	Comment       string  `json:"comment"`
	CreatedById   string  `json:"createdById"`
	CreatedByName string  `json:"createdByName"`
	Metadata      *string `json:"-"`
}
