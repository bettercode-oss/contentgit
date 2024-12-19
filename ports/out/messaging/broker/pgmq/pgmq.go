package pgmq

import (
	"contentgit/foundation"
	"contentgit/ports/out/messaging/broker"
	persistence "contentgit/ports/out/persistance"
	"contentgit/ports/out/persistance/eventsourcing"
	"contentgit/ports/out/persistance/eventsourcing/serializer"
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const vtDefault = 30

type PostgresMessagingQueue struct {
}

func NewPostgresMessagingQueue() *PostgresMessagingQueue {
	return &PostgresMessagingQueue{}
}

func (e *PostgresMessagingQueue) ProcessEvents(ctx context.Context, events []eventsourcing.Event) error {
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		eventBytes, err := serializer.Marshal(event)
		if err != nil {
			return errors.Wrap(err, "failed to marshal events")
		}

		err = e.PublishMessage(ctx, string(event.GetAggregateType()), eventBytes)
		if err != nil {
			return errors.Wrap(err, "failed to publish message")
		}
	}

	return nil
}

func (e *PostgresMessagingQueue) PublishMessage(ctx context.Context, queueName, message string) error {
	db := foundation.ContextProvider().GetDB(ctx)

	if err := db.Exec("SELECT * from pgmq.send(queue_name  => ?, msg => ?)", queueName, message).Error; err != nil {
		return errors.Wrap(err, "failed to send message to pgmq")
	}

	if queueName == "members" {
		if err := db.Exec("SELECT * from pgmq.send(queue_name  => ?, msg => ?)", "members_for_console", message).Error; err != nil {
			return errors.Wrap(err, "failed to send message to pgmq")
		}
	}

	return nil
}

func (e *PostgresMessagingQueue) ReadMessage(ctx context.Context, queueName string, vt uint) (*broker.MessageEnvelope, error) {
	if vt == 0 {
		vt = vtDefault
	}

	var messageEnvelope broker.MessageEnvelope
	db := foundation.ContextProvider().GetDB(ctx)

	if err := db.Raw("SELECT * FROM pgmq.read(?, ?, ?)", queueName, vt, 1).Scan(&messageEnvelope).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, persistence.ErrRecordNotFound
		}
		return nil, errors.Wrap(err, "failed to read message from pgmq")
	}

	if len(messageEnvelope.Message) == 0 {
		return nil, persistence.ErrRecordNotFound
	}

	return &messageEnvelope, nil
}

func (e *PostgresMessagingQueue) DeleteMessage(ctx context.Context, queueName string, msgId int64) (bool, error) {
	var deleted bool

	db := foundation.ContextProvider().GetDB(ctx)
	if err := db.Raw("SELECT pgmq.archive(?, ?::bigint)", queueName, msgId).Scan(&deleted).Error; err != nil {
		return false, errors.Wrap(err, "failed to delete message from pgmq")
	}

	return deleted, nil
}
