package content

import (
	"contentgit/domain/content/events"
	"contentgit/domain/content/projections"
	"contentgit/dtos"
	"contentgit/ports/out/persistance/eventsourcing"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type ContentEventHandler struct {
	serializer               eventsourcing.Serializer
	contentProjectRepository ContentProjectionRepository
}

func NewContentEventHandler(serializer eventsourcing.Serializer, contentProjectRepository ContentProjectionRepository) *ContentEventHandler {
	return &ContentEventHandler{serializer: serializer, contentProjectRepository: contentProjectRepository}
}

func (c *ContentEventHandler) Handle(ctx context.Context, esEvent eventsourcing.Event) error {
	deserializedEvent, err := c.serializer.DeserializeEvent(esEvent)
	if err != nil {
		return errors.Wrapf(err, "serializer.DeserializeEvent aggregateID: %s, type: %s", esEvent.GetAggregateID(), esEvent.GetEventType())
	}

	switch event := deserializedEvent.(type) {
	case *events.ContentCreatedEventV1:
		return c.onContentCreated(ctx, esEvent, event)

	case *events.FieldUpdatedEventV1:
		return c.onFieldUpdated(ctx, esEvent, event)

	case *events.FieldCommentAddedEventV1:
		return c.onFieldCommentAdded(ctx, esEvent, event)
	default:
		return errors.New(fmt.Sprintf("unknown event type: %s", esEvent.GetEventType()))
	}
}

func (c *ContentEventHandler) GetAggregateType() eventsourcing.AggregateType {
	return ContentAggregateType
}

func (c *ContentEventHandler) onContentCreated(ctx context.Context, esEvent eventsourcing.Event, event *events.ContentCreatedEventV1) error {
	if esEvent.GetVersion() != 1 {
		return errors.Wrapf(eventsourcing.ErrInvalidEventVersion, "type: %s, version: %d", esEvent.GetEventType(), esEvent.GetVersion())
	}

	contentProjection := projections.NewContentProjection(esEvent.AggregateID, esEvent.TenantId, event.Content, uint(esEvent.Version))
	if err := c.contentProjectRepository.Create(ctx, contentProjection); err != nil {
		return errors.Wrap(err, "failed to create content projection")
	}
	return nil
}

func (c *ContentEventHandler) onFieldUpdated(ctx context.Context, esEvent eventsourcing.Event, event *events.FieldUpdatedEventV1) error {
	contentProjection, err := c.contentProjectRepository.FindByID(ctx, esEvent.TenantId, esEvent.AggregateID)
	if err != nil {
		return errors.Wrap(err, "failed to find content projection")
	}

	updateField := dtos.ContentUpdateField{
		BeforeValue:   event.BeforeValue,
		AfterValue:    event.AfterValue,
		CreatedById:   event.CreatedById,
		CreatedByName: event.CreatedByName,
	}

	contentProjection.UpdateField(event.FieldName, updateField)
	contentProjection.Version = uint(esEvent.Version)

	return c.contentProjectRepository.Save(ctx, contentProjection)
}

func (c *ContentEventHandler) onFieldCommentAdded(ctx context.Context, esEvent eventsourcing.Event, event *events.FieldCommentAddedEventV1) error {
	contentProjection, err := c.contentProjectRepository.FindByID(ctx, esEvent.TenantId, esEvent.AggregateID)
	if err != nil {
		return errors.Wrap(err, "failed to find content projection")
	}
	contentProjection.AddFieldComment(event.FieldName, event.Comment, event.CreatedById, event.CreatedByName)
	contentProjection.Version = uint(esEvent.Version)

	return c.contentProjectRepository.Save(ctx, contentProjection)
}
