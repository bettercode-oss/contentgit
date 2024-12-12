package app

import (
	"contentgit/foundation"
	"contentgit/ports/out/messaging/broker/pgmq"
	"contentgit/ports/out/messaging/consumer"
	"context"
)

func (a *App) subscribeToEvents() {
	contentEventConsumer := consumer.NewEventConsumer(pgmq.NewPostgresMessagingQueue(), a.componentRegistry.Get("ContentEventHandler").(consumer.EventHandler))
	go func() {
		consumerCtx := foundation.ContextProvider().SetDB(context.TODO(), a.gormDB)
		contentEventConsumer.Consume(consumerCtx)
	}()
}
