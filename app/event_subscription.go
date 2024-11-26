package app

import (
	"contentgit/foundation"
	"contentgit/ports/out/persistance/eventsourcing"
	"contentgit/ports/out/queue/pgmq"
	"context"
)

func (a *App) subscribeToEvents() {
	contentEventConsumer := eventsourcing.NewEventConsumer(pgmq.NewPostgresMessagingQueue(), a.componentRegistry.Get("ContentEventHandler").(eventsourcing.EventHandler))
	go func() {
		consumerCtx := foundation.ContextProvider().SetDB(context.TODO(), a.gormDB)
		contentEventConsumer.Consume(consumerCtx)
	}()
}
