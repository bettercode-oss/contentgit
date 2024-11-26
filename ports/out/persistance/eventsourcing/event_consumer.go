package eventsourcing

import (
	persistence "contentgit/ports/out/persistance"
	"contentgit/ports/out/persistance/eventsourcing/serializer"
	"contentgit/ports/out/queue"
	"context"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type EventConsumer struct {
	messagingQueue queue.MessagingQueue
	eventHandler   EventHandler
}

func NewEventConsumer(messagingQueue queue.MessagingQueue, eventHandler EventHandler) *EventConsumer {
	return &EventConsumer{messagingQueue: messagingQueue, eventHandler: eventHandler}
}

func (c *EventConsumer) Consume(ctx context.Context) {
	messageChan := make(chan *queue.MessageEnvelope)
	errChan := make(chan error)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			messageEnvelope, err := c.messagingQueue.ReadMessage(ctx, string(c.eventHandler.GetAggregateType()), 30)
			if err != nil {
				if errors.Is(err, persistence.ErrRecordNotFound) {
					continue
				}
				errChan <- err
				continue
			}

			if messageEnvelope == nil {
				continue
			}

			messageChan <- messageEnvelope
		}
	}()

	for {
		select {
		case messageEnvelope := <-messageChan:
			event := Event{}
			err := serializer.Unmarshal(messageEnvelope.Message, &event)
			if err != nil {
				log.Error(err)
				continue
			}

			if err := c.eventHandler.Handle(ctx, event); err != nil {
				log.Error(err)
				continue
			}

			_, err = c.messagingQueue.DeleteMessage(ctx, string(c.eventHandler.GetAggregateType()), messageEnvelope.MsgId)
			if err != nil {
				log.Error(err)
				continue
			}
		case err := <-errChan:
			log.Error(err)
		case <-ctx.Done():
			return
		}
	}
}
