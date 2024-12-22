package content

import (
	"contentgit/domain/content/events"
	"contentgit/ports/out/persistance/eventsourcing"
	"context"

	"github.com/pkg/errors"
)

const (
	ContentAggregateType eventsourcing.AggregateType = "content"
)

type ContentAggregate struct {
	*eventsourcing.AggregateBase
	Content       map[string]any `json:"content"`
	ContentType   string         `json:"contentType"`
	FieldComments []FieldComment `json:"fieldComments"`
}

func NewContentAggregate(id string, tenantId string) (*ContentAggregate, error) {
	if id == "" || tenantId == "" {
		return nil, errors.New("id and tenantId are required.")
	}

	contentAggregate := &ContentAggregate{
		Content:       make(map[string]any),
		FieldComments: []FieldComment{},
	}

	aggregateBase := eventsourcing.NewAggregateBase(contentAggregate.When)
	aggregateBase.SetType(ContentAggregateType)
	aggregateBase.SetID(id)
	aggregateBase.SetTenantId(tenantId)
	contentAggregate.AggregateBase = aggregateBase

	return contentAggregate, nil
}

func NewContentAggregateWithType(id string, tenantId string, contentType string) (*ContentAggregate, error) {
	aggregate, err := NewContentAggregate(id, tenantId)
	if err != nil {
		return nil, err
	}

	aggregate.ContentType = contentType
	return aggregate, nil
}

func (a *ContentAggregate) CreateContent(ctx context.Context, content map[string]any) error {
	if content == nil {
		return errors.New("content is required.")
	}

	event := &events.ContentCreatedEventV1{
		Content:     content,
		ContentType: a.ContentType,
	}

	return a.Apply(event)
}

func (a *ContentAggregate) UpdateField(ctx context.Context, fieldName string, beforeValue any, afterValue any, createdById string, createdByName string) error {
	event := &events.FieldUpdatedEventV1{
		FieldName:     fieldName,
		BeforeValue:   beforeValue,
		AfterValue:    afterValue,
		CreatedById:   createdById,
		CreatedByName: createdByName,
	}

	return a.Apply(event)
}

func (a *ContentAggregate) AddFieldComment(ctx context.Context, fieldName string, comment string, createdById string, createdByName string) error {
	event := &events.FieldCommentAddedEventV1{
		FieldName:     fieldName,
		Comment:       comment,
		CreatedById:   createdById,
		CreatedByName: createdByName,
	}

	return a.Apply(event)
}

func (a *ContentAggregate) When(event any) error {
	switch evt := event.(type) {
	case *events.ContentCreatedEventV1:
		return a.handleContentCreatedEvent(evt)
	case *events.FieldUpdatedEventV1:
		return a.handleFieldUpdatedEvent(evt)
	case *events.FieldCommentAddedEventV1:
		return a.handleFieldCommentAddedEvent(evt)
	default:
		return errors.Wrapf(ErrUnknownEventType, "event: %#v", event)
	}
}

func (a *ContentAggregate) handleContentCreatedEvent(evt *events.ContentCreatedEventV1) error {
	a.Content = evt.Content
	a.ContentType = evt.ContentType
	return nil
}

func (a *ContentAggregate) handleFieldUpdatedEvent(evt *events.FieldUpdatedEventV1) error {
	contentFieldValue, ok := a.Content[evt.FieldName]
	if !ok {
		return ErrFieldNotFound
	}

	if contentFieldValue != evt.BeforeValue {
		return ErrFieldUpdateConflict
	}

	a.Content[evt.FieldName] = evt.AfterValue
	return nil
}

func (a *ContentAggregate) handleFieldCommentAddedEvent(evt *events.FieldCommentAddedEventV1) error {
	_, ok := a.Content[evt.FieldName]
	if !ok {
		return ErrFieldNotFound
	}

	fieldComment := Comment{
		Comment:       evt.Comment,
		CreatedById:   evt.CreatedById,
		CreatedByName: evt.CreatedByName,
	}

	for i := 0; i < len(a.FieldComments); i++ {
		if a.FieldComments[i].FieldName == evt.FieldName {
			a.FieldComments[i].Comments = append(a.FieldComments[i].Comments, fieldComment)
			return nil
		}
	}

	a.FieldComments = append(a.FieldComments, FieldComment{
		FieldName: evt.FieldName,
		Comments:  []Comment{fieldComment},
	})

	return nil
}

type FieldComment struct {
	FieldName string    `json:"fieldName"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	Comment       string `json:"comment"`
	CreatedById   string `json:"createdById"`
	CreatedByName string `json:"createdByName"`
}
